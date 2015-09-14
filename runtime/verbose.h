#ifndef VERBOSE_H
#define VERBOSE_H

#include <time.h>
#include <stdio.h>
#include "lib/trace.h"


#define Verbose_TRACE(s, f, x, r, level)        \
    do{                                         \
        int exist = Control_Verb_order(level, Control_verbose); \
        if (exist)                              \
        {                                       \
            Trace_spaces();                     \
            printf("%s starting\n", s);         \
            Trace_indent();                     \
        }                                       \
         r = f x;                               \
        if (exist)                              \
        {                                       \
            Trace_unindent();                   \
            Trace_spaces();                     \
            printf("%s finished\n", s);         \
        }                                       \
    }while(0)

#endif
