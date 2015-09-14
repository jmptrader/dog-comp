#ifndef STRING_H
#define STRING_H


#define T String_t

typedef char* String_t;



int String_equals(T, T);

long String_hashCode(T);

char** String_split(T x, T delim);

T String_concat(T s, ...);

T String_new(T x);

#undef T

#endif
