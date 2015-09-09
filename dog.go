package main

import (
	//"./ast"
	"./control"
	"./parser"
	"./util"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

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

	pser := parser.NewParse(filename, buf)
	Ast := pser.Parser()
	fmt.Printf("%T\n", Ast)

	//pp := ast.NewPP()
	//pp.DumpProg(Ast)
	_, filename, line, _ := runtime.Caller(0)
	util.Bug("test bug", filename, line)
}
