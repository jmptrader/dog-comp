#ifndef CONTROL_H
#define CONTROL_H



typedef enum {
    VERBOSE_SILENT,
    VERBOSE_PASS,
    VERBOSE_SUBPASS,
    VERBOSE_DETAIL
}Verbose_t;

// size of the Java heap (in bytes)
extern int Control_heapSize;

extern int Control_Verb_order(Verbose_t v1, Verbose_t v2);

extern Verbose_t Control_verbose;

extern void Control_setLogFile(FILE* fd);


#endif
