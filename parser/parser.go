package parser
/*

import (
    "fmt"
    "runtime"
    //"../lexer"
    "../ast"
    "../util"
    //"container/list"
    "strconv"
)

type Parser struct {
    lexer   *Lexer
    current *Token
    linenum int
}

func NewParse(fname string, buf []byte) *Parser {
    lexer := NewLexer(fname, buf)
    p := new(Parser)
    p.lexer = lexer
    p.current = p.lexer.NextToken()

    return p
}

func (this *Parser) skipLine() {
    for this.current.Kind == TOKEN_NEWLINE {
        this.eatToken(TOKEN_NEWLINE)
    }
}

func (this *Parser) advance() {
    this.linenum = this.current.LineNum
    this.current = this.lexer.NextToken()
    fmt.Println(this.current.ToString())
}

func (this *Parser) eatToken(kind int) {
    if kind == this.current.Kind {
        this.advance()
    } else {
        util.ParserError(tMap[kind], tMap[this.current.Kind], this.current.LineNum)
    }
}

func (this *Parser) parseType() ast.Type {
    switch this.current.Kind {
    case TOKEN_INT:
        this.eatToken(TOKEN_INT)
        return &ast.Int{}
    case TOKEN_BOOL:
        this.eatToken(TOKEN_BOOL)
        return &ast.Boolean{}
    case TOKEN_LBRACK:
        this.eatToken(TOKEN_LBRACK)
        this.eatToken(TOKEN_RBRACK)
        this.eatToken(TOKEN_INT)
        return &ast.IntArray{}
    case TOKEN_STAR:
        this.eatToken(TOKEN_STAR)
        name := this.current.Lexeme
        this.eatToken(TOKEN_ID)
        return &ast.ClassType{name}
    default:
        name := this.current.Lexeme
        this.eatToken(TOKEN_ID)
        return &ast.ClassType{name}
    }
}

func (this *Parser) parseFieldDec(id string) *ast.FieldDec {
    this.eatToken(TOKEN_ID)
    tp := this.parseType()
    return &ast.FieldDec{tp, id}
}

func (this *Parser) parseFieldDecs() []*ast.FieldDec {
    fields :=[]*ast.FieldDec{}
    id := this.current.Lexeme
    for this.skipLine(); this.current.Kind == TOKEN_ID; this.skipLine() {
        this.eatToken(TOKEN_ID)
        fields = append(fields, this.parseFieldDec(id))
    }

    return fields
}

func (this *Parser) parseStrDec() *ast.StructDec {

    this.eatToken(TOKEN_TYPE)
    id := this.current.Lexeme
    this.eatToken(TOKEN_ID)
    this.eatToken(TOKEN_STRUCT)
    this.eatToken(TOKEN_LBRACE)
    fields := this.parseFieldDecs()
    this.eatToken(TOKEN_RBRACE)

    return &ast.StructDec{id, fields}
}

func (this *Parser) parseStrDecs() []*ast.StructDec {
    strdecs := []*ast.StructDec{}
    for this.skipLine(); this.current.Kind == TOKEN_TYPE; this.skipLine() {
        strdecs = append(strdecs, this.parseStrDec())
        //strdecs.PushBack(this.parseStrDec())
    }
    return strdecs
}

func (this *Parser) parseFormalList() []*ast.FieldDec {
    flist := []*ast.FieldDec{}

    this.eatToken(TOKEN_LPAREN)
    if this.current.Kind == TOKEN_ID {
        id := this.current.Lexeme
        flist = append(flist, this.parseFieldDec(id))

        for this.current.Kind == TOKEN_COMMER {
            this.eatToken(TOKEN_COMMER)
            id := this.current.Lexeme
            flist = append(flist, this.parseFieldDec(id))
        }
    }
    this.eatToken(TOKEN_RPAREN)

    return flist
}

func (this *Parser) parseAtomExp() ast.Exp {
    switch this.current.Kind {
    case TOKEN_LEN:
        this.advance()
        this.eatToken(TOKEN_LPAREN)
        exp := this.parseExp()
        this.eatToken(TOKEN_RPAREN)
        return &ast.Len{exp}
    case TOKEN_SUB:
        this.advance()
        if this.current.Kind == TOKEN_NUM {
            num := this.current.Lexeme
            this.advance()
            s, _ := strconv.Atoi(num)
            s = -s
            return &ast.Num{s}
        } else {
            _, filename, line, _ := runtime.Caller(0)
            util.Bug("error", filename, line)
        }
    case TOKEN_LPAREN:
        this.advance()
        exp := this.parseExp()
        this.eatToken(TOKEN_RPAREN)
        return exp
    case TOKEN_NUM:
        value, _ := strconv.Atoi(this.current.Lexeme)
        this.advance()
        return &ast.Num{value}
    case TOKEN_TRUE:
        this.advance()
        return &ast.True{}
    case TOKEN_FALSE:
        this.advance()
        return &ast.False{}
    case TOKEN_ID:
        id := this.current.Lexeme
        this.advance()
        return &ast.Id{id, nil}
    case TOKEN_NEW:
        this.advance()
        this.eatToken(TOKEN_LPAREN)
        id := this.current.Lexeme
        this.eatToken(TOKEN_ID)
        this.eatToken(TOKEN_RPAREN)
        return &ast.NewObject{id}
    case TOKEN_MAKE:
        this.advance()
        this.eatToken(TOKEN_LPAREN)
        this.eatToken(TOKEN_LBRACK)
        this.eatToken(TOKEN_RBRACK)
        this.eatToken(TOKEN_INT)
        this.eatToken(TOKEN_COMMER)
        size := this.parseExp()
        this.eatToken(TOKEN_RPAREN)
        return &ast.NewIntArray{size}
    default:
        _, filename, line, _ := runtime.Caller(0)
        util.Bug("parse error", filename, line)
    }
    return nil
}

func (this *Parser) parseExpList() []ast.Exp {
    args := []ast.Exp{}
    if this.current.Kind == TOKEN_RPAREN {
        return args
    }

    args = append(args, this.parseExp())
    for this.current.Kind == TOKEN_COMMER {
        this.advance()
        args = append(args, this.parseExp())
    }
    return args
}

func (this *Parser) parseNotExp() ast.Exp {
    exp := this.parseAtomExp()
    for this.current.Kind == TOKEN_DOT||
    this.current.Kind == TOKEN_LBRACK{
        switch this.current.Kind {
        case TOKEN_DOT:
            this.advance()
            methodname := this.current.Lexeme
            this.eatToken(TOKEN_ID)
            this.eatToken(TOKEN_LPAREN)
            args := this.parseExpList()
            this.eatToken(TOKEN_RPAREN)
            return &ast.Call{exp, methodname, args, nil, nil}
        case TOKEN_LBRACK:
            this.advance()
            index := this.parseExp()
            this.eatToken(TOKEN_RBRACK)
            return &ast.ArraySelect{exp, index}
        default:
            _, filename, line, _ := runtime.Caller(0)
            util.Bug("need TOKEN_NOT or TOKEN_LBRACK", filename, line)
        }
    }
    return exp
}
func (this *Parser) parseTimeExp() ast.Exp {
    var exp2 ast.Exp
    for this.current.Kind == TOKEN_NOT {
        this.advance()
        exp2 = this.parseTimeExp()
    }
    if exp2 != nil {
        return &ast.Not{exp2}
    } else {
        return this.parseNotExp()
    }
}

func (this *Parser) parseAddSubExp() ast.Exp {
    left := this.parseTimeExp()
    for this.current.Kind == TOKEN_STAR {
        this.advance()
        right := this.parseTimeExp()
        return &ast.Times{left, right}
    }
    return left
}

func (this *Parser) parseLtExp() ast.Exp {
    left := this.parseAddSubExp()
    for this.current.Kind == TOKEN_ADD ||
    this.current.Kind == TOKEN_SUB {
        switch this.current.Kind {
        case TOKEN_ADD:
            this.advance()
            right := this.parseAddSubExp()
            return &ast.Add{left, right}
        case TOKEN_SUB:
            this.advance()
            right := this.parseAddSubExp()
            return &ast.Sub{left, right}
        default:
            _, filename, line, _ := runtime.Caller(0)
            util.Bug("need TOKEN_ADD or TOKEN_SUB", filename, line)
        }
    }
    return left
}

func (this *Parser) parseAndExp() ast.Exp {
    left := this.parseLtExp()
    for this.current.Kind == TOKEN_LT {
        this.advance()
        right := this.parseLtExp()
        return &ast.Lt{left, right}
    }
    return left
}

func (this *Parser) parseExp() ast.Exp {
    left := this.parseAndExp()
    for this.current.Kind == TOKEN_AND {
        this.advance()
        right := this.parseAndExp()
        return &ast.And{left, right}
    }
    return left
}

func (this *Parser) parseStatement() ast.Stm {
    switch this.current.Kind {
    case TOKEN_LBRACE:
        this.eatToken(TOKEN_LBRACE)
        stms := this.parseStatements()
        this.eatToken(TOKEN_RBRACE)
        return &ast.Block{stms}
    case TOKEN_ID:
        id := this.current.Lexeme
        this.eatToken(TOKEN_ID)
        switch this.current.Kind {
        case TOKEN_DERIVE:
            this.eatToken(TOKEN_DERIVE)
            e := this.parseExp()
            this.eatToken(TOKEN_NEWLINE)
            return &ast.Derive{id, e, nil}
        case TOKEN_ASSIGN:
            this.eatToken(TOKEN_ASSIGN)
            e := this.parseExp()
            this.eatToken(TOKEN_NEWLINE)
            return &ast.Assign{id, e, nil}
        case TOKEN_LBRACK:
            this.eatToken(TOKEN_LBRACK)
            index := this.parseExp()
            this.eatToken(TOKEN_RBRACK)
            this.eatToken(TOKEN_ASSIGN)
            e := this.parseExp()
            this.eatToken(TOKEN_NEWLINE)
            return &ast.AssignArray{id, index, e, nil}
        default:
            _, filename, line, _ := runtime.Caller(0)
            util.Bug("test bug", filename, line)
        }
    case TOKEN_IF:
        this.eatToken(TOKEN_IF)
        condition := this.parseExp()
        thenn := this.parseStatement()
        this.eatToken(TOKEN_ELSE)
        elsee := this.parseStatement()
        return &ast.If{condition, thenn, elsee}
    case TOKEN_FOR:
        this.eatToken(TOKEN_FOR)
        condition := this.parseExp()
        body := this.parseStatement()
        return &ast.For{condition, body}
    case TOKEN_FMT:
        this.eatToken(TOKEN_FMT)
        this.eatToken(TOKEN_DOT)
        this.eatToken(TOKEN_PRINTLN)
        this.eatToken(TOKEN_LPAREN)
        e := this.parseExp()
        this.eatToken(TOKEN_RPAREN)
        return &ast.Print{e}
    default:
        _, filename, line, _ := runtime.Caller(0)
        util.Bug("token error", filename, line)
    }
    return nil
}

func (this *Parser) parseStatements() []ast.Stm {
    stms := []ast.Stm{}
    for this.skipLine(); this.current.Kind == TOKEN_LBRACE ||
    this.current.Kind == TOKEN_ID ||
    this.current.Kind == TOKEN_IF ||
    this.current.Kind == TOKEN_FOR; this.skipLine() {
        stms = append(stms, this.parseStatement())
    }
    return stms
}

func (this *Parser) parseVarDec() *ast.VarDec {
    this.eatToken(TOKEN_VAR)
    id := this.current.Lexeme
    this.eatToken(TOKEN_ID)
    tp := this.parseType()

    return &ast.VarDec{tp, id}
}

func (this *Parser) parseVarDecs() []*ast.VarDec {
    decs := []*ast.VarDec{}
    for this.skipLine(); this.current.Kind == TOKEN_VAR; this.skipLine() {
        decs = append(decs, this.parseVarDec())
    }
    return decs
}

func (this *Parser) parseMethod() ast.Func {
    this.skipLine()

    this.eatToken(TOKEN_FUNC)
    this.eatToken(TOKEN_LPAREN)
    firstarg := this.current.Lexeme
    this.eatToken(TOKEN_ID)
    this.eatToken(TOKEN_STAR)
    bindingType := this.parseType()
    this.eatToken(TOKEN_RPAREN)
    method_name := this.current.Lexeme
    this.eatToken(TOKEN_ID)
    formals := this.parseFormalList()
    rettype := this.parseType()
    this.eatToken(TOKEN_LBRACE)
    locals := this.parseVarDecs()
    stms := this.parseStatements()
    this.eatToken(TOKEN_RETURN)
    retExp := this.parseExp()
    this.skipLine()
    this.eatToken(TOKEN_RBRACE)

    return &ast.MethodSingle{firstarg, bindingType, method_name, formals,rettype ,locals, stms, retExp}
}

func (this *Parser) parseMethods() []ast.Func {
    methods := []ast.Func{}
    for this.skipLine(); this.current.Kind == TOKEN_FUNC; this.skipLine() {
        methods = append(methods, this.parseMethod())
        fmt.Printf("current.Kind=%s\n", tMap[this.current.Kind])
    }
    return methods
}

func (this *Parser) parseMainFunc() ast.MainFunc {
    this.eatToken(TOKEN_FUNC)
    this.eatToken(TOKEN_MAIN)
    this.eatToken(TOKEN_LPAREN)
    this.eatToken(TOKEN_RPAREN)
    this.eatToken(TOKEN_LBRACE)
    this.skipLine()
    stm := this.parseStatement()
    this.skipLine()
    this.eatToken(TOKEN_RBRACE)
    return &ast.MainFuncSingle{stm}
}

func (this *Parser) parseProgram() ast.Prog {
    this.skipLine()
    strdecs := this.parseStrDecs()
    mainfunc := this.parseMainFunc()
    methods := this.parseMethods()
    return &ast.ProgramSingle{mainfunc, nil, strdecs, methods}
}

func (this *Parser) Parser() ast.Prog {
    p := this.parseProgram()
    return p
}
*/
