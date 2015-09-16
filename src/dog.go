package main

import (
	"./ast"
	"./ast/optimization"
	"./codegen/C"
	"./control"
	"./elaborator"
	"./parser"
	"./util"
	"fmt"
	"io/ioutil"
	"os"
)

func dog_Parser(filename string, buf []byte) ast.Program {
	return parser.NewParse(filename, buf).Parser()
}

func main() {
	args := os.Args[1:len(os.Args)]
	filename := control.Do_arg(args)
	if filename == "" {
		control.Usage()
		os.Exit(0)
	}
	control.CodeGen_fileName = filename
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if control.Lexer_test == true {
		lex := parser.NewLexer(filename, buf)
		tk := lex.NextToken()
		for tk.Kind != parser.TOKEN_EOF {
			fmt.Println(tk.ToString())
			tk = lex.NextToken()
		}
		fmt.Println(tk.ToString())
		os.Exit(0)
	}
	//setp1: lexer&&parser
	Ast := dog_Parser(filename, buf)
	if control.Ast_dumpAst == true {
		ast.NewPP().DumpProg(Ast)
	}
	//step2: elaborate
	elaborator.Elaborate(Ast)

	Ast = ast_opt.DeadClass_Opt(Ast)

	//set3: codegen
	var Ast_c codegen_c.Program
	switch control.CodeGen_codegen {
	case control.C:
		Ast_c = codegen_c.TransC(Ast)
	case control.Bytecode:
		util.Todo()
	case control.Dalvik:
		util.Todo()
	case control.X86:
		util.Todo()
	default:
		panic("impossible")
	}

	codegen_c.CodegenC(Ast_c)
}
