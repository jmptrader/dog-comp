package parser

import (
	"fmt"
	"os"
)

type Lexer struct {
	fname   string
	s       string
	lineNum int
	buf     []byte
	/**
	 * the buf index, can increse or decrese to implements reset
	 */
	fp int
}

func NewLexer(fname string, buf []byte) *Lexer {
	initTokenMap()
	lex := Lexer{}
	lex.fname = fname
	lex.s = ""
	lex.lineNum = 1
	lex.buf = buf
	lex.fp = 0

	return &lex
}

func (this *Lexer) NextToken() *Token {
	var t *Token
	t = nil
	for t == nil {
		t = this.nextTokenInternal()
	}
	return t
}

func (this *Lexer) expectKeyword(expect string) bool {
	reset := this.fp
	for _, e := range expect {
		if e == int32(this.buf[this.fp]) {
			this.fp++
			continue
		} else {
			this.fp = reset
			return false
		}
	}
	return true
}

func (this *Lexer) expectIdOrKey(c byte) *Token {
	kind, exist := tokenMap[this.s]
	if exist {
		tk := newToken(kind, this.s, this.lineNum)
		this.s = ""
		this.fp--
		return tk
	} else if this.s == "" {
		if c != ' ' {
			kk := tokenMap[string(c)]
			tk := newToken(kk, string(c), this.lineNum)
			return tk
		} else {
			return nil
		}
	} else {
		tk := newToken(TOKEN_ID, this.s, this.lineNum)
		this.s = ""
		this.fp--
		return tk
	}
}

func (this *Lexer) lex_Comments(c byte) {
	ex := this.buf[this.fp]
	this.fp++
	if ex == '/' {
		for ex != '\n' && this.fp < len(this.buf) {
			ex = this.buf[this.fp]
			this.fp++
		}
		if this.fp == len(this.buf) {
			this.fp--
			return
		} else {
			this.lineNum++
		}
	} else if ex == '*' {
		ex = this.buf[this.fp]
		for (c != '*' || ex != '/') && this.fp < len(this.buf) {
			c = ex
			ex = this.buf[this.fp]
			this.fp++
		}
		if this.fp == len(this.buf) {
			fmt.Println("error")
			os.Exit(0)
		}
	} else {
		fmt.Println("error")
		os.Exit(0)
	}
}

func (this *Lexer) lex_Num(c byte) string {
	var s string
	s += string(c)

	for {
		next := this.buf[this.fp]
		this.fp++
		if next >= '0' && next <= '9' {
			s += string(next)
			continue
		}

		//999abc is not number
		if (next == '_') || (next >= 'a' && next <= 'z') ||
			(next >= 'A' && next <= 'Z') {
			fmt.Println("ilegal number")
			os.Exit(0)
		}
		break
	}

	this.fp--
	return s

}

func (this *Lexer) nextTokenInternal() *Token {
	if this.fp == len(this.buf) {
		return newToken(TOKEN_EOF, "EOF", this.lineNum)
	}

	c := this.buf[this.fp]
	this.fp++

	for c == '\t' || '\n' == c || '\r' == c {
		if c == '\n' {
			this.lineNum++
		}
		if this.fp >= len(this.buf) {
			return newToken(TOKEN_EOF, "EOF", this.lineNum)
		}
		c = this.buf[this.fp]
		this.fp++
	}

	if this.fp >= len(this.buf) {
		return newToken(TOKEN_EOF, "EOF", this.lineNum)
	}

	switch c {
	case '&':
		if this.s == "" {
			if this.expectKeyword("&") {
				return newToken(TOKEN_AND, "&&", this.lineNum)
			} else {
				panic("expect &&")
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case ' ':
		fallthrough
	case '+':
		fallthrough
	case '=':
		fallthrough
	case ',':
		fallthrough
	case '.':
		fallthrough
	case '{':
		fallthrough
	case '[':
		fallthrough
	case '(':
		fallthrough
	case '<':
		fallthrough
	case '!':
		fallthrough
	case '}':
		fallthrough
	case ']':
		fallthrough
	case ')':
		fallthrough
	case ';':
		fallthrough
	case '-':
		fallthrough
	case '*':
		return this.expectIdOrKey(c)
	case '0':
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		fallthrough
	case '4':
		fallthrough
	case '5':
		fallthrough
	case '6':
		fallthrough
	case '7':
		fallthrough
	case '8':
		fallthrough
	case '9':
		if this.s == "" {
			return newToken(TOKEN_NUM, this.lex_Num(c), this.lineNum)
		}
		this.s += string(c)
	case '/':
		this.lex_Comments(c)
	default:
		this.s += string(c)
	}

	return nil

}
