#include <stdlib.h>
#include <stdio.h>
#include "list.h"
#include "string.h"
#include "poly.h"
#include "assert.h"

#define L List_t
#define P Poly_t

static const int STEP = 2;
static int indent = 0;
static L traceList;

void Trace_indent()
{
    indent += STEP;
}

void Trace_unindent()
{
    indent -= STEP;
}

void Trace_spaces()
{
    int i = indent;
    while (i--)
    {
        printf(" ");
    }
}

int Trace_contains(char* s)
{
    if (!traceList)
        traceList = List_new();

    int exist =  List_contains(traceList,
                s,
                (Poly_tyEquals)String_equals);
    return exist;
}

void Trace_addFunc(char* s)
{
    if (!traceList)
      traceList = List_new();

    List_addFirst(traceList, s);
}

#undef L
#undef P
