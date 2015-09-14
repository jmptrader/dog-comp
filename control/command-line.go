package control

import (
	"../util"
	"fmt"
	"os"
	"strings"
)

type Kind int
type Arg struct {
	name       string
	option     string
	desription string
	kind       Kind
	action     func(interface{})
}

const (
	EMPTY = iota
	BOOL
	INT
	STRING
	STRINGLIST
)

const (
	VERSION = "dog v0.0.1 linux/386"
	WEBSITE = "https://github.com/qc1iu/dog-comp"
)

func printSpeaces(i int) int {
	r := i
	for ; i > 0; i-- {
		fmt.Print(" ")
	}
	return r
}

func printAllArg() {
	const INDENT_1 = 12
	const INDENT_2 = 36
	for _, arg := range all_Arg {
		i, _ := fmt.Print("  " + arg.name)
		i += printSpeaces(INDENT_1 - i)
		_i, _ := fmt.Print(arg.option)
		i += _i
		i += printSpeaces(INDENT_2 - i)
		fmt.Println(arg.desription)
	}
}

func Usage() {
	fmt.Println("The Dog compiler. Copyright (C) 2015, CSS of USTC.\n")
	fmt.Println("Usage:\n")
	fmt.Println("\tcommand [arguments]\n")
	fmt.Println("The commands are:\n")
	printAllArg()
	fmt.Println("")
	fmt.Println(VERSION)
	fmt.Printf("See %s for more details.\n", WEBSITE)
}

func argException(s ...interface{}) {
	fmt.Print("Args error: ")
	for _, v := range s {
		fmt.Print(v)
		fmt.Print(" ")
	}
	fmt.Println("")
	os.Exit(1)
}

var all_Arg = []Arg{
	{"codegen",
		"{bytecode|C|dalvik|x86}",
		"which code generator to use",
		STRING,
		func(c interface{}) {
			switch c.(type) {
			case string:
				if c == "bytecode" {
					Control_CodeGen_codegen = Bytecode
				} else if c == "C" {
					Control_CodeGen_codegen = C
				} else if c == "dalvik" {
					Control_CodeGen_codegen = Dalvik
				} else if c == "x86" {
					Control_CodeGen_codegen = X86
				} else {
					argException("-codegen {bytecode|C|dalvik|x86}")
				}
			default:
				argException("bad argument")
			}
		}},
	{"dump",
		"{ast}",
		"dump information about the given ir",
		STRING,
		func(c interface{}) {
			switch c.(type) {
			case string:
				if c == "ast" {
					Control_Ast_dumpAst = true
				} else {
					argException("-dump {ast}")
				}
			default:
				argException("bad argument")
			}
		}},
	{"elab",
		"{classtable|methodtable}",
		"dump information about elaboration",
		STRING,
		func(c interface{}) {
			if s, ok := c.(string); ok {
				if s == "classtable" {
					Control_Elab_classTable = true
				} else if s == "methodtable" {
					Control_Elab_methodTable = true
				} else {
					argException("-elab {classtable|methodtable}")
				}
			} else {
				argException("bad argument")
			}
		}},
	{"lex",
		"",
		"dump the result of lexical analysis",
		EMPTY,
		func(c interface{}) {
			Control_Lexer_dump = true
		}},
	{"testlexer",
		"",
		"whether or not to test the lexer",
		EMPTY,
		func(c interface{}) {
			Control_Lexer_test = true
		}},
	{"output",
		"<outfile>",
		"set the name of the output file",
		STRING,
		func(c interface{}) {
			if s, ok := c.(string); ok {
				Control_CodeGen_outputName = s
			} else {
				panic("impossible")
			}
		}},
	{"help",
		"",
		"show this help information",
		EMPTY,
		nil},
}

func Do_arg(args []string) string {
	filename := ""
	found := false
	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			if filename == "" {
				filename = args[i]
				continue
			} else {
				argException("can only compile one Java file")
			}
		} else {
		}
		for _, arg := range all_Arg {
			if !strings.EqualFold(arg.name, strings.TrimPrefix(args[i], "-")) {
				continue
			}
			found = true
			switch arg.kind {
			case EMPTY:
				arg.action(nil)
			case BOOL:
				util.Todo()
			case INT:
				util.Todo()
				if i >= len(args) {
					argException("need <INT>")
				}
			case STRING:
				i++
				if i >= len(args) {
					argException("need <STRING>")
				}
				theArg := args[i]
				(arg.action)(theArg)
			case STRINGLIST:
				util.Todo()
			default:
			}
			break
		}
		if !found {
			argException("invalid option " + args[i])
		}
	}
	return filename
}
