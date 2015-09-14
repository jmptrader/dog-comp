#include <stdio.h>
#include <stdlib.h>
#include <string.h>

extern void Tiger_main();
extern void Tiger_heap_init (int);
extern void CommandLine_doarg (int argc, char **argv);
int main (int argc, char **argv)
{
  // Lab 4, exercise 13:
  // You should add some command arguments to the generated executable
  // to control the behaviour of your Gimple garbage collector.
  // For instance, you should run:
  //   $ a.out @tiger -heapSize 1 @
  // to set the Java heap size to 1K. Or you can run
  //   $ a.out @tiger -gcLog @
  // to generate the log (which is discussed in this exercise).
  // You can use the offered function in file "control.c"
  // and "command-line.c"
  // Your code here:

  // initialize the Java heap

  CommandLine_doarg(argc,argv);

  Tiger_heap_init (Control_heapSize);

  // enter Java code...
  Tiger_main ();
}
