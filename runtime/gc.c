/*^_^*--------------------------------------------------------------*//*{{{*/
/* Copyright (C) SSE-USTC, 2014-2015                                */
/*                                                                  */
/*  FILE NAME             :  gc.c                                   */
/*  PRINCIPAL AUTHOR      :  qc1iu                                  */
/*  LANGUAGE              :  C                                      */
/*  TARGET ENVIRONMENT    :  ANY                                    */
/*  DATE OF FIRST RELEASE :  2014/10/05                             */
/*  DESCRIPTION           :  the tiger compiler 'gc                 */
/*------------------------------------------------------------------*/

/*
 * Revision log:
 * ---------------------
 * 2014/12/06
 * 1>add Exchange()
 * 2>add RewriteObj()
 * --------------------
 *
 * 2014/12/08
 * 1>add the copyCount
 *
 * 2015/08/26
 *  refactor by qc1iu
 *
 *
 */
/*}}}*/
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include "control.h"
#include "verbose.h"
#include "lib/assert.h"
#include "lib/error.h"

#define O Object_t
#define F Frame_t
#define V Vtable_t

/**
 * If use clang, should add the option -D__CLANG__
 */
#ifdef __CLANG__
#define GET_STACK_ARG_ADDRESS(base, index)  \
    (((char*)base)-(index)*sizeof(void*))
#else
#define GET_STACK_ARG_ADDRESS(base, index)  \
    (((char*)base)+(index)*sizeof(void*))
#endif

#define CAST_ADDR(p)        ((int*)(p))
#define GET_LOCAL_ADDRESS(frame, index)     \
    (((char*)(frame+1))+index*(sizeof(int)))

#define GET_OBJECT_FIELD_ADDR(obj, index)   \
    (((char*)(obj+1))+index*(sizeof(int)))

#define GET_OBJECT(addr)    (*(O*)addr)

#define NO_ENOUGH_SPACE(remains, need)      \
    ((remains)<(need))

#define ADDRESS_COMPARE(add1, op, add2)     \
    ((int*)(add1)op(int*)(add2))

#define IN_TO_SPACE(addr)                             \
    (ADDRESS_COMPARE(addr, <, heap.toStart+heap.size)&&    \
     ADDRESS_COMPARE(addr, >=, heap.toStart))

#define IN_FROM_SPACE(addr)                           \
    (ADDRESS_COMPARE(addr, <, heap.from+heap.size)&&  \
     ADDRESS_COMPARE(addr, >=, heap.from))

typedef enum
{
    TYPE_OBJECT,
    TYPE_ARRAY
}OBJECT_TYPE;

typedef struct O *O;
typedef struct F *F;
typedef struct V *V;


struct V
{
    char* class_map;
};

/*    
      ----------------
      | vptr         | (this field should be empty for an array)
      |--------------|
      | isObjOrArray | (1: for array)
      |--------------|
      | length       |
      |--------------|
      | forwarding   |
      |--------------|\
      p---->| e_0          | \
      |--------------|  s
      | ...          |  i
      |--------------|  z
      | e_{length-1} | /e
      ----------------/
      */
struct O
{
    V vptr;
    OBJECT_TYPE isArray;
    int length;
    void* forwarding;
};


struct F
{
    F prev;
    char* arguments_gc_map;
    int* arguments_base_address;
    int locals_gc_map;
};



extern int Log;

extern char* logname;

/**
 * The "prev" pointer, pointing to the top frame on the GC stack.
 */
void *previous = 0;

static int copyCount;

static int gcNum;

static const int WORLD = sizeof(int);

static const int OBJECT_HEADER_SIZE = sizeof(struct O);

static const int FRAME_HEADER_SIZE = sizeof(struct F);


// The Gimple Garbage Collector.


//===============================================================//
// The Java Heap data structure.

/*
   ----------------------------------------------------
   |From Space              |To Space                 |
   ----------------------------------------------------
   ^\                      /^
   | \<~~~~~~~ size ~~~~~>/ |
   from                       to
   */
struct JavaHeap
{
    int size;         // in bytes, note that this is for semi-heap size
    char *from;       // the "from" space pointer
    char *fromFree;   // the next "free" space in the from space
    char *to;         // the "to" space pointer
    char *toStart;    // "start" address in the "to" space
    char *toNext;     // "next" free space pointer in the to space
};

/**
 * The Java heap, which is initialized by the following
 * "heap_init" function.
 */
struct JavaHeap heap;

/**
 * Given the heap size (in bytes), allocate a Java heap
 * in the C heap, initialize the relevant fields.
 */
void Tiger_heap_init (int heapSize)
{
    // #1: allocate a chunk of memory of size "heapSize"
    char* jheap=(char*)malloc(heapSize);

    // #2: initialize the "size" field, note that "size" field
    // is for semi-heap, but "heapSize" is for the whole heap.
    heap.size=heapSize/2;
    // #3: initialize the "from" field.
    heap.from=(char*)jheap;
    // #4: initialize the "fromFree" field.
    heap.fromFree=heap.from;
    // #5: initialize the "to" field.
    heap.to=heap.fromFree+heap.size;
    // #6: initizlize the "toStart" field. 
    heap.toNext=(char*)heap.to+1;
    // #7: initialize the "toNext" field.
    heap.toStart=(char*)heap.to+1;

    return;
}

static void dumpObject(O obj)
{
    fprintf(stdout, "-----------------\n");
    fprintf(stdout, "obj->vptr %x\n", (unsigned int)obj->vptr);
    fprintf(stdout, "obj->isArray %x\n", obj->isArray);
    fprintf(stdout, "obj->length %x\n", obj->length);
    fprintf(stdout, "obj->forwarding %x\n", (unsigned int)obj->forwarding);
    fprintf(stdout, "-----------------\n\n");

}

//===============================================================//
// The Gimple Garbage Collector
// A copying collector based-on Cheney's algorithm.


static int swapAndCleanUp()
{
    char* swap;

    swap = heap.from;
    heap.from = heap.toStart;
    heap.to = (char*)heap.from+heap.size;
    heap.fromFree = heap.toNext;
    heap.toStart = swap;
    heap.toNext = swap;
    memset(heap.toStart, 0, heap.size);

    return 0;
}


static int objectSize(O obj)
{
    int size;

    size = 0;
    switch (obj->isArray)
    {
        case TYPE_OBJECT:
            size = obj->length;
            break;
        case TYPE_ARRAY:
            size = obj->length*sizeof(int)+OBJECT_HEADER_SIZE;
            break;
        default:
            ERROR("wrong type");
    }

    Assert_ASSERT(size != 0);
    return size;
}


static int copyCollection(int** addr_addr)
{
    O old_obj;
    O new_obj;

    old_obj = (O)*addr_addr;

    if (!IN_FROM_SPACE(old_obj))
      return 0;

    void* forwarding = old_obj->forwarding;
    if (IN_TO_SPACE(forwarding))
    {
        *addr_addr = (int*)forwarding;
        return 0;
    }
    else if (IN_FROM_SPACE(forwarding)||ADDRESS_COMPARE(forwarding, ==, 0))
    {
        copyCount++;
        new_obj = (O)heap.toNext;

        int size = objectSize(old_obj);
        memcpy(new_obj, old_obj, size);
        old_obj->forwarding = new_obj;

        //XXX copy finished, also need to change original address
        *addr_addr = (int*)forwarding;
        heap.toNext+=size;

        return 0;
    }
    else
    {
        ERROR("impossible!");
    }

    return 0;
}

static void collectedObjectField(O obj)
{
    Assert_ASSERT(obj->isArray == TYPE_OBJECT);

    V vtable;
    char* class_map;
    int field_count;

    vtable = obj->vptr;
    class_map = vtable->class_map;
    Assert_ASSERT(class_map);
    field_count = strlen(class_map);

    if (field_count <= 0)
      return;

    int i=0;
    for (i=0; i<field_count; i++)
    {
        if (class_map[i] == '0')
          continue;

        int r;
        Verbose_TRACE("copyCollection", copyCollection, 
                    ((int**)GET_OBJECT_FIELD_ADDR(obj, i)), r, VERBOSE_SUBPASS);
    }
}

static int collectedField()
{
    char* to_ptr;
    O obj;

    for(to_ptr = heap.toStart; copyCount > 0; copyCount--)
    {
        obj = (O)to_ptr;
        switch (obj->isArray)
        {
            case TYPE_OBJECT:
                Assert_ASSERT(obj->vptr);
                collectedObjectField(obj);
                break;
            case TYPE_ARRAY:
                break;
            default:
                ERROR("impossible");
        }
        to_ptr = (char*)to_ptr + objectSize(obj);
    }

    return 0;
}

static int doArg(F frame)
{
    int len;
    char* arg_map;

    arg_map = frame->arguments_gc_map;
    if (arg_map == NULL)
      return 0;

    len = strlen(arg_map);
    int i = 0;
    for (i=0; i<len; i++)
    {
        if (arg_map[i] == '0')
          continue;

        int* arg_addr = CAST_ADDR(
                    GET_STACK_ARG_ADDRESS(
                        frame->arguments_base_address, i));
        int r;
        Verbose_TRACE("copyCollection", copyCollection, 
                    ((int**)arg_addr), r, VERBOSE_SUBPASS);
    }

    return 0;
}

static int doLocals(F frame)
{
    int locals_map;

    locals_map = frame->locals_gc_map;
    Assert_ASSERT(locals_map>=0);
    if (locals_map == 0)
      return 0;

    int i=0;
    for (i=0; i<locals_map; i++)
    {
        int r;
        Verbose_TRACE("copyCollection", copyCollection,
                    ((int**)GET_LOCAL_ADDRESS(frame, i)), r, VERBOSE_SUBPASS);
    }

    return 0;
}

static int frameSingle(F frame)
{
    if (frame == NULL)
      return 0;

    int r;
    Verbose_TRACE("doArg", doArg, (frame), r, VERBOSE_SUBPASS);
    Verbose_TRACE("doLocals", doLocals, (frame), r, VERBOSE_SUBPASS);
    Verbose_TRACE("frameSingle", frameSingle, (frame->prev), r, VERBOSE_SUBPASS);

    return 0;
}

static int Verbose_Tiger_gc()
{
    Assert_ASSERT(copyCount == 0);

    int r;
    Verbose_TRACE("frameSingle", frameSingle, (previous), r, VERBOSE_SUBPASS);
    Verbose_TRACE("collectedfield", collectedField, (), r, VERBOSE_SUBPASS);
    Verbose_TRACE("swapAndCleanUp", swapAndCleanUp, (), r, VERBOSE_SUBPASS);

    return 0;
}

static void Tiger_gc ()
{
    int before_gc;
    clock_t start;
    clock_t end;
    int gcByte;
    double sec;

    before_gc = heap.to-heap.fromFree;
    start = clock();
    int r;
    Verbose_TRACE("Tiger_gc", Verbose_Tiger_gc, (), r, VERBOSE_PASS);
    end = clock();
    gcByte = heap.to-heap.fromFree-before_gc;
    sec =  (double)(end-start)/CLOCKS_PER_SEC;
    gcNum++;
    if(Log)
    {
        FILE *fp;
        fp = fopen(logname, "at");
        if (!fp)
        {
            printf("File cannot be opened");
            exit(1);
        }
        fprintf(fp,"%d round of gc :%fs.%d byte reclaim\n",
                    gcNum, sec, gcByte);
        fclose(fp);
    }
    return;
}


//===============================================================//
// Object Model And allocation


// "new" a new object, do necessary initializations, and
// return the pointer (reference).
/*    ----------------
      | vptr      ---|----> (points to the virtual method table)
      |--------------|
      | isObjOrArray | (0: for normal objects)
      |--------------|
      | length       | (this field should be empty for normal objects)
      |--------------|
      | forwarding   |
      |--------------|\
p---->| v_0          | \
      |--------------|  s
      | ...          |  i
      |--------------|  z
      | v_{size-1}   | /e
      ----------------/
      */
/**
 * Try to allocate an object in the "from" space of the Java
 * heap. Read Tiger book chapter 13.3 for details on the
 * allocation.
 * There are two cases to consider:
 *   1. If the "from" space has enough space to hold this object, then
 *      allocation succeeds, return the apropriate address (look at
 *      the above figure, be careful);
 *   2. if there is no enough space left in the "from" space, then
 *      you should call the function "Tiger_gc()" to collect garbages.
 *      and after the collection, there are still two sub-cases:
 *        a: if there is enough space, you can do allocations just as case 1;
 *        b: if there is still no enough space, you can just issue
 *           an error message ("OutOfMemory") and exit.
 *           (However, a production compiler will try to expand
 *           the Java heap.)
 */
static void* newObject(void* vtable, int size)
{
    O obj;

    obj  = (O)heap.fromFree;
    memset(obj, 0,  size);
    obj->vptr = vtable;
    obj->isArray = TYPE_OBJECT;
    obj->length = size;
    obj->forwarding = 0;

    heap.fromFree+=size;

    return obj;
}

static void* newArray(int length)
{
    O obj;

    obj = (O)heap.fromFree;
    memset(obj ,0,length*sizeof(int)+OBJECT_HEADER_SIZE);
    obj->vptr = NULL;
    obj->isArray = TYPE_ARRAY;
    obj->length = length;
    obj->forwarding = 0;
    heap.fromFree += (length*sizeof(int))+OBJECT_HEADER_SIZE;


    return obj;
}

void *Tiger_new (void *vtable, int size)
{
    if(NO_ENOUGH_SPACE(heap.to-heap.fromFree, size))
    {
        Tiger_gc();

        if(NO_ENOUGH_SPACE(heap.to-heap.fromFree, size))
        {
            printf("Tiger_gc can not collecte enough space...\n");
            printf("There is %d byte remained,but you need:%d\n", 
                        (int)(heap.to-heap.fromFree),size);
            ERROR("OutOfMemory. Use '-help' to see the command-line option");
        }
    }

    O obj = newObject(vtable, size);

    return obj;
}

// "new" an array of size "length", do necessary
// initializations. And each array comes with an
// extra "header" storing the array length and other information.
/*    ----------------
      | vptr         | (this field should be empty for an array)
      |--------------|
      | isObjOrArray | (1: for array)
      |--------------|
      | length       |
      |--------------|
      | forwarding   |
      |--------------|\
p---->| e_0          | \
      |--------------|  s
      | ...          |  i
      |--------------|  z
      | e_{length-1} | /e
      ----------------/
      */

/**
 * Try to allocate an array object in the "from" space of the Java
 * heap. Read Tiger book chapter 13.3 for details on the
 * allocation.
 * There are two cases to consider:
 *  1. If the "from" space has enough space to hold this array object, then
 *     allocation succeeds, return the apropriate address (look at
 *     the above figure, be careful);
 *
 * //XXX Return the base address. And in codegenC,
 *       array[0] translate to array[0+4].
 *       array.length translate to array[2].
 *
 *  2. if there is no enough space left in the "from" space, then
 *     call the function "Tiger_gc()" to collect garbages.
 *     and after the collection, there are still two sub-cases:
 *       a: if there is enough space, do allocations just as case 1;
 *       b: if there is still no enough space, just issue
 *          an error message ("OutOfMemory") and exit.
 *          (However, a production compiler will try to expand
 *          the Java heap.)
 */
void *Tiger_new_array (int length)
{
    if(NO_ENOUGH_SPACE(heap.to-heap.fromFree,(
                        length*sizeof(int))+OBJECT_HEADER_SIZE))
    {
        Tiger_gc();
        if(NO_ENOUGH_SPACE(heap.to-heap.fromFree, 
                        (length*sizeof(int))+OBJECT_HEADER_SIZE))
        {
            ERROR("OutOfMemory. Use '-help' to see the command-line option");
        }
    }
    O array = newArray(length);

    return (array);
}


#undef O
#undef F
#undef V
