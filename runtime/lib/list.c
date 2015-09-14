#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "poly.h"
#include "list.h"
#include "mem.h"
#include "error.h"
#include "assert.h"


#define T List_t
#define P Poly_t


T List_new() 
{
    T l;
    Mem_new(l);
    l->data = 0;
    l->next = 0;

    return l;
}


/**
 * @parm 
 *      x data
 * @parm
 *      l next pointer
 *
 * @return a new List_t struct 
 */
static T List_newNode(P x, T l) 
{
    T p;
    Mem_new(p);
    p->data = x;
    p->next = l;

    return p;
}

 /**
  * Inserts the specified element at the beginning of this list.
  *
  * @param e the element to add
  * NOTE: we use head->data to record tail
  */
void List_addFirst(T l, P x) 
{
    T t;
    Assert_ASSERT(l);

    t = List_new();
    t = List_newNode(x, l->next);
    l->next = t;

    if (l->data == NULL) {
        l->data = t;
    }
    return;
}

/**
 * Appends the specified element to the end of this list.
 *
 * @param e the element to add
 */
void List_addLast(T l, P x) 
{
    T tail, p;
    Assert_ASSERT(l);

    if (l->next == NULL) {
        List_addFirst(l, x);
        return;
    } 

    tail = (T)l->data;
    p = List_newNode(x, NULL);
    tail->next = p;
    l->data = p;

    return;
}


/**
 * need traivals the list
 */
int List_size(T l) 
{
    Assert_ASSERT(l);
    T p;
    int i = 0;
    p = l->next;
    while (p) {
        i++;
        p = p->next;
    }

    return i;
}

int List_isEmpty(T l)
{
    Assert_ASSERT(l);
    return (0==l->next);
}

/**
 * Returns the element at the specified position in this list.
 *
 * @param index index of the element to return
 * @return the element at the specified position in this list
 */
P List_getIndexOf(T l, int index) 
{
    Assert_ASSERT(l);

    T p = l->next;
    if (index < 0) {
        ERROR("invalid argument");
        return 0;
    }
    while (p) {
        if (0 == index)
            return p->data;//found!
        index--;
        p=p->next;
    }
    return 0;
}

/**
 * getIndexOf(l, 0)
 */
P List_getFirst(T l)
{
    Assert_ASSERT(l);
    return l->next;
}

/**
 * return the first fit
 *
 */
static P List_containsInternal(T l, P x, Poly_tyEquals f)
{
    T p;
    Assert_ASSERT(l);
    Assert_ASSERT(f);

    p = l->next;
    while (p) {
        if (f(x, p->data))
          return p->data;

        p = p->next;
    }

    return 0;
}

/**
 * Returns {@code true} if this list contains the specified element.
 *
 * @param o element whose presence in this list is to be tested
 * @return {@code true} if this list contains the specified element
 *                                         */
int List_contains(T l, P x, Poly_tyEquals f)
{
    P p = List_containsInternal(l, x, f);
    if (p)
      return 1;
    else
      return 0;
}

/**
 *
 *
 */
P List_removeFirst(T l)
{
    T p;
    Assert_ASSERT(l);

    if (0 == l->next)
      ERROR("List size is 0\n"); 

    p = l->next;
    l->next = l->next->next;
    //XXX omit the l->data. 

    return p->data;

}
/**
 * Removes the first occurrence of the specified element from this list,
 * if it is present (optional operation).  If this list does not contain
 * the element, it is unchanged.  More formally, removes the element with
 * the lowest index
 *
 * @param o element to be removed from this list, if present.
 * @return <tt>data</tt> if this list contained the specified element.
 */
P List_remove(T l, P x, Poly_tyEquals equals)
{
    T prev;
    T current;
    P r;

    Assert_ASSERT(l);
    Assert_ASSERT(x);
    Assert_ASSERT(equals);

    prev = l;
    current = l->next;
    while (current)
    {
        if (equals(x, current->data))
        {
            r = current->data;
            current = current->next;
            prev->next = current;

            //XXX free(current)
            return r;
        }

        prev = current;
        current = current->next;

    }

    return NULL;
}



#undef T
#undef P
