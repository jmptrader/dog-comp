package parser
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
    currentNext *Token
    isSpecial bool
    isField bool
    linenum int
}

func NewParse(fname string, buf []byte) *Parser {
    lexer := NewLexer(fname, buf)
    p := new(Parser)
    p.lexer = lexer
    p.current = p.lexer.NextToken()

    return p
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
    case TOKEN_BOOLEAN:
        this.eatToken(TOKEN_BOOLEAN)
        return &ast.Boolean{}
    case TOKEN_LBRACK:
        this.eatToken(TOKEN_LBRACK)
        this.eatToken(TOKEN_RBRACK)
        this.eatToken(TOKEN_INT)
        return &ast.IntArray{}
    default:
        name := this.current.Lexeme
        this.eatToken(TOKEN_ID)
        return &ast.ClassType{name}
    }
}

func (this *Parser) parseFormalList() []ast.Dec {
    flist := []ast.Dec{}
    var tp ast.Type
    var id string

    if this.current.Kind == TOKEN_ID ||
        this.current.Kind == TOKEN_INT ||
        this.current.Kind == TOKEN_BOOLEAN{
            tp = this.parseType()
            id = this.current.Lexeme
            this.eatToken(TOKEN_ID)
            flist = append(flist, &ast.DecSingle{tp, id, this.isField})

        for this.current.Kind == TOKEN_COMMER {
            this.eatToken(TOKEN_COMMER)
            tp = this.parseType()
            id = this.current.Lexeme
            this.eatToken(TOKEN_ID)
            flist = append(flist, &ast.DecSingle{tp, id, this.isField})
        }
    }
    return flist
}

//AtomExp   -> (exp)
//          -> INTEGER_LITERAL
//          -> true
//          -> false
//          -> this
//          -> id
//          -> new int[exp]
//          -> new id()
func (this *Parser) parseAtomExp() ast.Exp {
    switch this.current.Kind {
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
    case TOKEN_THIS:
        this.advance()
        return &ast.This{}
    case TOKEN_ID:
        id := this.current.Lexeme
        this.advance()
        return &ast.Id{id, nil, false}
    case TOKEN_NEW:
        this.advance()
        switch this.current.Kind {
        case TOKEN_INT:
            this.advance()
            this.eatToken(TOKEN_LBRACK)
            exp := this.parseExp()
            this.eatToken(TOKEN_RBRACK)
            return &ast.NewIntArray{exp}
        case TOKEN_ID:
            s := this.current.Lexeme
            this.advance()
            this.eatToken(TOKEN_LPAREN)
            this.eatToken(TOKEN_RPAREN)
            return &ast.NewObject{s}
        default:
            _, filename, line, _ := runtime.Caller(0)
            util.Bug("parse error", filename, line)
        }
    default:
        fmt.Println(tMap[this.current.Kind])
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

//NotExp    -> AtomExp
//          -> AtomExp.id(explist)
//          -> AtomExp[exp]
//          -> AtomExp.length
func (this *Parser) parseNotExp() ast.Exp {
    exp := this.parseAtomExp()
    for this.current.Kind == TOKEN_DOT||
    this.current.Kind == TOKEN_LBRACK{
        switch this.current.Kind {
        case TOKEN_DOT:
            this.advance()
            if this.current.Kind == TOKEN_LENGTH {
                this.advance()
                return &ast.Length{exp}
            }
            //else ast.Call
            methodname := this.current.Lexeme
            this.eatToken(TOKEN_ID)
            this.eatToken(TOKEN_LPAREN)
            args := this.parseExpList()
            this.eatToken(TOKEN_RPAREN)
            return &ast.Call{exp, methodname, args, "", nil, nil}
        case TOKEN_LBRACK://[exp]
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

//TimesExp  -> !TimesExp
//          -> NotExp
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

//AddSubExp -> TimesExp * TimesExp
//          -> TimesExp
func (this *Parser) parseAddSubExp() ast.Exp {
    left := this.parseTimeExp()
    for this.current.Kind == TOKEN_TIMES {
        this.advance()
        right := this.parseTimeExp()
        return &ast.Times{left, right}
    }
    return left
}

//LtExp -> AddSubExp + AddSubExp
//      -> AddSubExp - AddSubExp
//      -> AddSubExp
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

//AndExp    -> LtExp < LtExp
//          -> LtExp
func (this *Parser) parseAndExp() ast.Exp {
    left := this.parseLtExp()
    for this.current.Kind == TOKEN_LT {
        this.advance()
        right := this.parseLtExp()
        return &ast.Lt{left, right}
    }
    return left
}

//Exp -> AndExp && AndExp
//    -> AndExp
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
        if this.isSpecial == true {
            this.current = this.currentNext
            switch this.current.Kind {
            case TOKEN_ASSIGN:
                this.eatToken(TOKEN_ASSIGN)
                e := this.parseExp()
                this.eatToken(TOKEN_SEMI)
                this.isSpecial = false
                assign := new(ast.Assign)
                assign.Name = id
                assign.E = e
                return assign
            case TOKEN_LBRACK:
                this.eatToken(TOKEN_LBRACK)
                index := this.parseExp()
                this.eatToken(TOKEN_RBRACK)
                this.eatToken(TOKEN_ASSIGN)
                e := this.parseExp()
                this.eatToken(TOKEN_SEMI)
                this.isSpecial = false
                return &ast.AssignArray{id, index, e, nil, false}
            default:
                _, filename, line, _ := runtime.Caller(0)
                util.Bug("test bug", filename, line)
            }
        }else {
            this.eatToken(TOKEN_ID)
            switch this.current.Kind {
            case TOKEN_ASSIGN:
                this.eatToken(TOKEN_ASSIGN)
                exp := this.parseExp()
                this.eatToken(TOKEN_SEMI)
                assign := new(ast.Assign)
                assign.Name = id
                assign.E = exp
                return assign
            case TOKEN_LBRACK:
                this.eatToken(TOKEN_LBRACE)
                index := this.parseExp()
                this.eatToken(TOKEN_RBRACK)
                this.eatToken(TOKEN_ASSIGN)
                exp := this.parseExp()
                this.eatToken(TOKEN_SEMI)
                return &ast.AssignArray{id, index, exp, nil, false}
            default:
                _, filename, line, _ := runtime.Caller(0)
                util.Bug("test bug", filename, line)
            }
        }
    case TOKEN_IF:
        this.eatToken(TOKEN_IF)
        this.eatToken(TOKEN_LPAREN)
        condition := this.parseExp()
        this.eatToken(TOKEN_RPAREN)
        thenn := this.parseStatement()
        this.eatToken(TOKEN_ELSE)
        elsee := this.parseStatement()
        return &ast.If{condition, thenn, elsee}
    case TOKEN_WHILE:
        this.eatToken(TOKEN_WHILE)
        this.eatToken(TOKEN_LPAREN)
        exp := this.parseExp()
        this.eatToken(TOKEN_RPAREN)
        body := this.parseStatement()
        return &ast.While{exp, body}
    case TOKEN_SYSTEM:
        this.eatToken(TOKEN_SYSTEM)
        this.eatToken(TOKEN_DOT)
        this.eatToken(TOKEN_OUT)
        this.eatToken(TOKEN_DOT)
        this.eatToken(TOKEN_PRINTLN)
        this.eatToken(TOKEN_LPAREN)
        e := this.parseExp()
        this.eatToken(TOKEN_RPAREN)
        this.eatToken(TOKEN_SEMI)
        return &ast.Print{e}
    default:
        _, filename, line, _ := runtime.Caller(0)
        util.Bug("token error", filename, line)
    }
    return nil
}

func (this *Parser) parseStatements() []ast.Stm {
    stms := []ast.Stm{}
    for this.current.Kind == TOKEN_LBRACE ||
    this.current.Kind == TOKEN_ID ||
    this.current.Kind == TOKEN_IF ||
    this.current.Kind == TOKEN_WHILE ||
    this.current.Kind == TOKEN_SYSTEM{
        stms = append(stms, this.parseStatement())
    }
    return stms
}

func (this *Parser) parseVarDecl() ast.Dec {
    var dec *ast.DecSingle
    var id string

    if !this.isSpecial {
        tp := this.parseType()
        id := this.current.Lexeme
        dec = &ast.DecSingle{tp, id, this.isField}
        this.eatToken(TOKEN_ID)
        this.eatToken(TOKEN_SEMI)
        return dec
    }else {
        tp := &ast.ClassType{this.current.Lexeme}
        this.current = this.currentNext
        id = this.current.Lexeme
        dec = &ast.DecSingle{tp, id, this.isField}
        this.eatToken(TOKEN_ID)
        this.eatToken(TOKEN_SEMI)
        this.isSpecial = false
        return dec
    }
}

func (this *Parser) parseVarDecls() []ast.Dec {
    decs := []ast.Dec{}
    for this.current.Kind == TOKEN_INT ||
            this.current.Kind == TOKEN_BOOLEAN||
            this.current.Kind == TOKEN_ID{
        if this.current.Kind != TOKEN_ID {
            decs = append(decs, this.parseVarDecl())
        }else {
            id := this.current.Lexeme
            linenum := this.current.LineNum
            this.eatToken(TOKEN_ID)
            if this.current.Kind == TOKEN_ASSIGN {
                this.currentNext = this.current
                this.current = newToken(TOKEN_ID, id, linenum)
                this.isSpecial = true
                return decs
            }else if this.current.Kind == TOKEN_LBRACK {
               this.currentNext = this.current
               this.current = newToken(TOKEN_ID, id, linenum)
               this.isSpecial = true
               return decs
            }else {
                this.currentNext = this.current
                this.current = newToken(TOKEN_ID, id, linenum)
                this.isSpecial = true
                decs = append(decs, this.parseVarDecl())
            }
        }
    }
    return decs
}

func (this *Parser) parseMethod() ast.Method {
    this.eatToken(TOKEN_PUBLIC)
    reType := this.parseType()
    method_name := this.current.Lexeme
    this.eatToken(TOKEN_ID)
    this.eatToken(TOKEN_LPAREN)
    formals := this.parseFormalList()
    this.eatToken(TOKEN_RPAREN)
    this.eatToken(TOKEN_LBRACE)
    locals := this.parseVarDecls()
    stms := this.parseStatements()
    this.eatToken(TOKEN_RETURN)
    retExp := this.parseExp()
    this.eatToken(TOKEN_SEMI)
    this.eatToken(TOKEN_RBRACE)

    return &ast.MethodSingle{reType,
                            method_name,
                            formals,
                            locals,
                            stms,
                            retExp}
}

func (this *Parser) parseMethodDecls() []ast.Method {
    methods := []ast.Method{}
    for this.current.Kind == TOKEN_PUBLIC {
        this.isField = false
        methods = append(methods, this.parseMethod())
    }
    this.isField = true
    return methods
}

func (this *Parser) parseClassDecl() ast.Class{
    var id, extends string

    this.eatToken(TOKEN_CLASS)
    id = this.current.Lexeme
    this.eatToken(TOKEN_ID)
    if this.current.Kind == TOKEN_EXTENDS {
        this.eatToken(TOKEN_EXTENDS)
        extends = this.current.Lexeme
        this.eatToken(TOKEN_ID)
    }
    this.eatToken(TOKEN_LBRACE)
    decs := this.parseVarDecls()
    methods := this.parseMethodDecls()
    this.eatToken(TOKEN_RBRACE)
    return &ast.ClassSingle{id, extends, decs, methods}
}


func (this *Parser) parseClassDecls() []ast.Class {
    classes := []ast.Class{}
    for this.current.Kind == TOKEN_CLASS{
        classes = append(classes, this.parseClassDecl())
        //fmt.Printf("current.Kind=%s\n", tMap[this.current.Kind])
    }
    return classes
}

func (this *Parser) parseMainClass() ast.MainClass {
    this.eatToken(TOKEN_CLASS)
    id := this.current.Lexeme
    this.eatToken(TOKEN_ID)
    this.eatToken(TOKEN_LBRACE)
    this.eatToken(TOKEN_PUBLIC)
    this.eatToken(TOKEN_STATIC)
    this.eatToken(TOKEN_VOID)
    this.eatToken(TOKEN_MAIN)
    this.eatToken(TOKEN_LPAREN)
    this.eatToken(TOKEN_STRING)
    this.eatToken(TOKEN_LBRACK)
    this.eatToken(TOKEN_RBRACK)
    arg := this.current.Lexeme
    this.eatToken(TOKEN_ID)
    this.eatToken(TOKEN_RPAREN)
    this.eatToken(TOKEN_LBRACE)
    stm := this.parseStatement()
    this.eatToken(TOKEN_RBRACE)
    this.eatToken(TOKEN_RBRACE)
    return &ast.MainClassSingle{id, arg, stm}
}

func (this *Parser) parseProgram() ast.Program {
    main_class := this.parseMainClass()
    classes := this.parseClassDecls()
    return &ast.ProgramSingle{main_class, classes}
}

func (this *Parser) Parser() ast.Program {
    p := this.parseProgram()
    return p
}
