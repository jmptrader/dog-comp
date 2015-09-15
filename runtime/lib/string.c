#include <string.h>
#include "string.h"
#include <stdarg.h>
#include "assert.h"
#include "mem.h"

#define T String_t


#define MULTIPLIER 31
long String_hashCode(T s)
{
    long h=0;

    while(*s)
    {
        h = h*MULTIPLIER+(unsigned)*s++;
    }
    return h;
}

int String_equals(T x, T y)
{
    if (0 == strcmp((char*)x, (char*)y))
      return 1;

    return 0;
}

char** String_split(T x, T delim)
{
    Assert_ASSERT(x);
    Assert_ASSERT(delim);

    int len = strlen(x);
    int len2 = strlen(delim);
    Assert_ASSERT(len2==1);

    char* buf;
    Mem_newSize(buf, len+1);
    strcpy(buf, x);
    buf[len] = '\0';

    char c = delim[0];
    int size = 0;
    char* p = x;
    while(*p)
    {
        if (*p == c)
          size++;
        p++;
    }
    size+=2;

    char** t;
    Mem_newSize(t, size);
    char* value = buf;
    char* result = NULL;

    int i = 0;
    result = strsep(&value, delim);

    for (i = 0; result != NULL; i++)
    {
        char* s;
        Mem_newSize(s, strlen(result)+1);
        strcpy(s, result);
        s[strlen(result)] = '\0';
        t[i] = s;

        result = strsep(&value, delim);
    }
    t[i] = NULL;

    return t;
}

T String_concat(T s, ...)
{
    int totalSize = 0;
    char* current = s;
    char* temp;
    char* head;

    va_list ap;
    va_start(ap, s);
    while (current)
    {
        totalSize += strlen(current);
        current = va_arg(ap, char *);
    }
    va_end(ap);


    Mem_newSize(temp, (totalSize+1));
    head = temp;
    current = s;
    va_start(ap, s);
    while (current)
    {
        strcpy(temp, current);
        temp += strlen(current);
        current = va_arg(ap, char *);
    }

    return head;
}

T String_new(T x)
{
    int len;
    char* s;

    len = strlen(x);
    Mem_newSize(s, len+1);
    strcpy(s, x);

    if (String_equals(s, x))
      return s;

    ERROR("string_new err!");
}



#undef MULTIPLIER
#undef T
