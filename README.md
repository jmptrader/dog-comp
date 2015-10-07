DOG
===
[![Build Status](https://drone.io/github.com/qc1iu/dog-comp/status.png)](https://drone.io/github.com/qc1iu/dog-comp/latest)

The Dog compiler. Copyright (C) 2015, CSS of USTC.

`Dog` is a **MiniJava lanuage**(see [Tiger book](http://www.cs.princeton.edu/~appel/modern/java/) Appendix MiniJava Language
Reference Manual) compile implement in **`Golang`**. Just for fun, it has a garbage collector using copy collection algorithm and CFG-based Optimizations.

##Getting Start

### compile to C programming language
use 

	dog-comp/src$ make
	dog-comp/src$ ./dog-comp ../test/LinkedList.java

can compile LinkedList.java to a.c. And you can compile this C source file with runtime.

### more command option
	Usage:

		command [arguments]

	The commands are:

		-codegen   {bytecode|C|dalvik|x86}  which code generator to use
		-dump      {ast|c}                  dump information about the given ir
		-elab      {classtable|methodtable} dump information about elaboration
		-lex                                dump the result of lexical analysis
		-testlexer                          whether or not to test the lexer
		-o         <outfile>                set the name of the output file
		-trace     <passname>               trace compile pass
		-skip      <passname>               skip opt pass
		-verbose   {0|1|2|3}                verbose pass
		-visualize {pdf|ps|svg|jpg}         graph format
		-opt       {0~6}                    optimization level
		-help                               show this help information


###Play In Runtime

Since the evaluation order of the argument list, when using `clang`, need add the compile options **`__CLANG__`**
	
	dog-comp/src$ clang a.c ../runtime/runtime.c -D__CLANG__

## CFG-based Optimizations
Dog-comp has various optimizations on  [control-flow graph](https://en.wikipedia.org/wiki/Control_flow_graph) (CFG). CFG is a graph-based program representation with basic blocks as nodes and control transfers between blocks as edges. A control-flow graph is a good candidate intermediate representation for it makes optimizations easier to perform.

For example: 

use
	
	dog-comp/src$ ./dog-comp ../test/LinkedList.java -visualize svg 

to generate CFG for LinkedList.java.

![](https://raw.githubusercontent.com/qc1iu/dog-comp/master/screenshots/Element_Equal.jpg)

CFG-based program representations make program optimizations much easier, due to the fact that it's much easier to calculate the control flow and data flow information required to perform optimizations. Generally speaking, an optimization can be divided into two phases: [program analysis](https://en.wikipedia.org/wiki/Program_analysis) and rewriting. During the program analysis phase (usually data-flow analysis), the dog-comp analyzes the target program to get useful static information as an approximation of the program's dynamic behaviour. And such information will be used in the second phase: program rewriting, the dog-comp rewrites the program according to the static information obtained in the first phase. As a result fo these two steps, the program is altered in some way---it is optimized. For instance, in the [dead-code elimination optimization](https://github.com/qc1iu/dog-comp/blob/master/src/cfg/optimization/dead-code.go), the program analysis phase will determine whether or not the variable x defined by x = e will be used anywhere in this program; if x will not be used, the program rewriting phase will remove this assignment from this program (and the program shrinks). 

Dog-comp solves three [data-flow equations](https://en.wikipedia.org/wiki/Data-flow_analysis):

1. [Liveness Analysis](https://en.wikipedia.org/wiki/Live_variable_analysis)
2. [Reaching Definitions](https://en.wikipedia.org/wiki/Reaching_definition)
3. [Available Expressions](https://en.wikipedia.org/wiki/Available_expression)

and then do some optimization

1. Dead-code Elimination
2. Constant Propagation
3. Copy Propagation
4. Common Sub-Expression Elimination (CSE)



##GC Support

In tiger, we use a gc named **Gimple garbage collector**, whitch means `gc is simple`. And we use the algorithm called [Cheney's algorithm](https://en.wikipedia.org/wiki/Cheney's_algorithm) which uses a breadth-first strategy.

###Object Model

An object model is the strategy that how an object can be laid out in memory (the heap). A good object model should support all kinds of possible operations efficiently on an object: virtual method dispatching, locking, garbage collection, etc., whereas at the same time be as compact as possible to save storage. 

There are two different forms of objects in MiniJava: normal objects and (integer) array objects. So, to support virtual method dispatching on normal objects, each object is represented as a sequence of its fields plus a unique virtual method table pointer at the zero offset as a header. And to support array length operations on array objects, arrays contain a length field as a header.


	// "new" a new object, do necessary initializations, and
	// return the pointer (reference).
	
		  ----------------
	      | vptr      ---|----> (points to the virtual method table)
	      |--------------|
	      | isObjOrArray | (0: for normal objects)
	      |--------------|
	      | length       | (this field should be empty for normal objects)
	      |--------------|
	      | forwarding   |
	      |--------------|\
	p---->| v_0          | \
	      |--------------|  s
	      | ...          |  i
	      |--------------|  z
	      | v_{size-1}   | /e
	      ----------------/


	// "new" an array of size "length", do necessary
	// initializations. And each array comes with an
	// extra "header" storing the array length and other information.
		  ----------------
		  | vptr         | (this field should be empty for an array)
		  |--------------|
		  | isObjOrArray | (1: for array)
		  |--------------|
		  | length       |
		  |--------------|
	      | forwarding   |
		  |--------------|\
	p---->| e_0          | \
      	  |--------------|  s
      	  | ...          |  i
      	  |--------------|  z
      	  | e_{length-1} | /e
      	  ----------------/



###command-line

	-heapSize <n>            set the Java heap size (in bytes)
	-verbose {0|1|2|3}       trace method execuated
   	-gcLog                   generate GClog
   	-help                    help

