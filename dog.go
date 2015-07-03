package main

import (
	//"./lexer"
	"./parser"
	"./util"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

func main() {

	fmt.Println("dog-cmp Start...\n")

	filename := "test/sum.go"
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	/*
		lex := parser.NewLexer(filename, buf)
		tk := lex.NextToken()
		for tk.Kind != parser.TOKEN_EOF {
			fmt.Println(tk.ToString())
			tk = lex.NextToken()
		}
		fmt.Println(tk.ToString())

		_, filename, line, _ := runtime.Caller(0)
		util.Bug("test bug", filename, line)
	*/
	pser := parser.NewParse(filename, buf)
	Ast := pser.Parser()
	fmt.Printf("%T\n", Ast)
	_, filename, line, _ := runtime.Caller(0)
	util.Bug("test bug", filename, line)
}
