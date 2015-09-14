#ifndef ERROR_H
#define ERROR_H

#include <stdio.h>
#include <stdlib.h>

#define ERROR(s)                    \
    do {                            \
        fprintf(stderr,             \
                    "\e[31m\e[1mERROR>\e[0m%s %s:%d\n",       \
                    s,              \
                    __FILE__,       \
                    __LINE__        \
               );                   \
        fflush(stderr);             \
        exit(1);                    \
    }while (0)


#define TODO(s)                    \
    do {                            \
        fprintf(stderr,             \
                    "\e[33m\e[1mTODO>\e[0m%s %s:%d\n",       \
                    s,              \
                    __FILE__,       \
                    __LINE__        \
               );                   \
        fflush(stderr);             \
        exit(1);                    \
    }while (0)

#define WARNING(s)                  \
    do {                            \
        fprintf(stderr,             \
                    "\e[33m\e[1mWARNING>\e[0m%s %s:%d\n",       \
                    s,              \
                    __FILE__,       \
                    __LINE__        \
               );                   \
        fflush(stderr);             \
    }while (0)


#endif
