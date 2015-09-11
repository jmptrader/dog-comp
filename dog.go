package main

import (
	"./ast"
	"./control"
	"./elaborator"
	"./parser"
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
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if control.Control_Lexer_test == true {
		lex := parser.NewLexer(filename, buf)
		tk := lex.NextToken()
		for tk.Kind != parser.TOKEN_EOF {
			fmt.Println(tk.ToString())
			tk = lex.NextToken()
		}
		fmt.Println(tk.ToString())
		os.Exit(0)
	}
	Ast := dog_Parser(filename, buf)
	if control.Control_Ast_dumpAst == true {
		ast.NewPP().DumpProg(Ast)
	}
	elaborator.Elaborate(Ast)
}
