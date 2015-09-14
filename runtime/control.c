#include <stdio.h>
#include "control.h"

int Control_heapSize = 512;

/*{{{ Verbose*/
Verbose_t Control_verbose = VERBOSE_SILENT;
Verbose_t Control_verboseDefault = VERBOSE_SILENT;


int Control_Verb_order(Verbose_t v1, Verbose_t v2)
{
    return v1<=v2;
}
/*}}}*/

