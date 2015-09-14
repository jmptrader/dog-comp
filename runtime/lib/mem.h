#ifndef MEM_H
#define MEM_H

#include "error.h"


#define Mem_new(p)                  \
    do {                            \
        (p) = malloc(sizeof(*(p))); \
    }while(0)

#define Mem_newSize(p, n)                   \
    do {                                    \
        if (n<0)                            \
            ERROR("invalid buffer size");   \
        (p)=malloc((n)*sizeof(*(p)));       \
    }while(0)

#endif
