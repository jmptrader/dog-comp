package parser

import "fmt"
import "os"
import "strconv"

type Token struct {
	Kind    int
	Lexeme  string
	LineNum int
}

func newToken(kind int, lexeme string, lineNum int) *Token {
	return &Token{kind, lexeme, lineNum}
}

var tokenMap map[string]int
var tMap map[int]string

func initTokenMap() {
	tokenMap = make(map[string]int)
	tokenMap["+"] = TOKEN_ADD
	tokenMap["&&"] = TOKEN_AND
	tokenMap["="] = TOKEN_ASSIGN
	tokenMap["bool"] = TOKEN_BOOL
	tokenMap["struct"] = TOKEN_STRUCT
	tokenMap[","] = TOKEN_COMMER
	tokenMap["."] = TOKEN_DOT
	tokenMap["else"] = TOKEN_ELSE
	tokenMap["EOF"] = TOKEN_EOF
	tokenMap["false"] = TOKEN_FALSE
	//id
	tokenMap["if"] = TOKEN_IF
	tokenMap["int"] = TOKEN_INT
	tokenMap["{"] = TOKEN_LBRACE
	tokenMap["["] = TOKEN_LBRACK
	tokenMap["len"] = TOKEN_LEN
	tokenMap["("] = TOKEN_LPAREN
	tokenMap["<"] = TOKEN_LT
	tokenMap["main"] = TOKEN_MAIN
	tokenMap["new"] = TOKEN_NEW
	tokenMap["make"] = TOKEN_MAKE
	tokenMap["!"] = TOKEN_NOT
	tokenMap[";"] = TOKEN_NEWLINE
	//num
	tokenMap["fmt"] = TOKEN_FMT
	tokenMap["Println"] = TOKEN_PRINTLN
	tokenMap["func"] = TOKEN_FUNC
	tokenMap["}"] = TOKEN_RBRACE
	tokenMap["]"] = TOKEN_RBRACK
	tokenMap["return"] = TOKEN_RETURN
	tokenMap[")"] = TOKEN_RPAREN
	tokenMap[":="] = TOKEN_DERIVE
	tokenMap[";"] = TOKEN_SEMI
	tokenMap["*"] = TOKEN_STAR
	tokenMap["-"] = TOKEN_SUB
	tokenMap["true"] = TOKEN_TRUE
	tokenMap["void"] = TOKEN_VOID
	tokenMap["var"] = TOKEN_VAR
	tokenMap["for"] = TOKEN_FOR
	tokenMap["type"] = TOKEN_TYPE

	tMap = make(map[int]string)
	tMap[TOKEN_ADD] = "TOKEN_ADD"
	tMap[TOKEN_AND] = "TOKEN_AND"
	tMap[TOKEN_ASSIGN] = "TOKEN_ASSIGN"
	tMap[TOKEN_BOOL] = "TOKEN_BOOL"
	tMap[TOKEN_STRUCT] = "TOKEN_STRUCT"
	tMap[TOKEN_COMMER] = "TOKEN_COMMER"
	tMap[TOKEN_DOT] = "TOKEN_DOT"
	tMap[TOKEN_ELSE] = "TOKEN_ELSE"
	tMap[TOKEN_EOF] = "TOKEN_EOF"
	tMap[TOKEN_FALSE] = "TOKEN_FALSE"
	tMap[TOKEN_IF] = "TOKEN_IF"
	tMap[TOKEN_INT] = "TOKEN_INT"
	tMap[TOKEN_ID] = "TOKEN_ID"
	tMap[TOKEN_LBRACE] = "TOKEN_LBRACE"
	tMap[TOKEN_LBRACK] = "TOKEN_LBRACK"
	tMap[TOKEN_LEN] = "TOKEN_LEN"
	tMap[TOKEN_LPAREN] = "TOKEN_LPAREN"
	tMap[TOKEN_LT] = "TOKEN_LT"
	tMap[TOKEN_MAIN] = "TOKEN_MAIN"
	tMap[TOKEN_NEW] = "TOKEN_NEW"
	tMap[TOKEN_NEWLINE] = "TOKEN_NEWLINE"
	tMap[TOKEN_NUM] = "TOKEN_NUM"
	tMap[TOKEN_MAKE] = "TOKEN_MAKE"
	tMap[TOKEN_NOT] = "TOKEN_NOT"
	tMap[TOKEN_FMT] = "TOKEN_FMT"
	tMap[TOKEN_PRINTLN] = "TOKEN_PRINTLN"
	tMap[TOKEN_FUNC] = "TOKEN_FUNC"
	tMap[TOKEN_RBRACE] = "TOKEN_RBRACE"
	tMap[TOKEN_RBRACK] = "TOKEN_RBRACK"
	tMap[TOKEN_RETURN] = "TOKEN_RETURN"
	tMap[TOKEN_RPAREN] = "TOKEN_RPAREN"
	tMap[TOKEN_DERIVE] = "TOKEN_DERIVE"
	tMap[TOKEN_SEMI] = "TOKEN_SEMI"
	tMap[TOKEN_STAR] = "TOKEN_STAR"
	tMap[TOKEN_SUB] = "TOKEN_SUB"
	tMap[TOKEN_TRUE] = "TOKEN_TRUE"
	tMap[TOKEN_VOID] = "TOKEN_VOID"
	tMap[TOKEN_VAR] = "TOKEN_VAR"
	tMap[TOKEN_FOR] = "TOKEN_FOR"
	tMap[TOKEN_TYPE] = "TOKEN_TYPE"

}

type Kind int

const (
	TOKEN_ADD = iota
	TOKEN_AND
	TOKEN_ASSIGN
	TOKEN_BOOL
	TOKEN_STRUCT
	TOKEN_COMMER
	TOKEN_DOT
	TOKEN_DERIVE
	TOKEN_ELSE
	TOKEN_EOF
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUNC
	TOKEN_FMT
	TOKEN_ID
	TOKEN_IF
	TOKEN_INT
	TOKEN_LBRACE
	TOKEN_LBRACK
	TOKEN_LEN
	TOKEN_LPAREN
	TOKEN_LT
	TOKEN_MAIN
	TOKEN_MAKE
	TOKEN_NEW
	TOKEN_NEWLINE
	TOKEN_NOT
	TOKEN_NUM
	TOKEN_PRINTLN
	TOKEN_RBRACE
	TOKEN_RBRACK
	TOKEN_RETURN
	TOKEN_RPAREN
	TOKEN_STAR
	TOKEN_SEMI
	TOKEN_SUB
	TOKEN_TRUE
	TOKEN_TYPE
	TOKEN_VAR
	TOKEN_VOID
)

func (this *Token) ToString() string {
	var s string
	if this.LineNum == 0 {
		fmt.Println("error")
		os.Exit(0)
	}

	s = ": " + this.Lexeme + " at LINE:" + strconv.Itoa(this.LineNum)
	return tMap[this.Kind] + s
}
