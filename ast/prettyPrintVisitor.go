package ast

import (
	"../util"
	"fmt"
	"runtime"
	"strconv"
)

type PrettyPrintVisitor struct {
	indentLevel int
}

func NewPP() *PrettyPrintVisitor {
	pp := new(PrettyPrintVisitor)
	pp.indentLevel = 2

	return pp
}

func (this *PrettyPrintVisitor) indent() {
	this.indentLevel += 2
}

func (this *PrettyPrintVisitor) unIndent() {
	this.indentLevel -= 2
	if this.indentLevel <= 0 {
		_, filename, line, _ := runtime.Caller(0)
		util.Bug("indent error", filename, line)
	}
}

func (this *PrettyPrintVisitor) sayln() {
	this.say("\n")
	this.printSpeaces()
}
func (this *PrettyPrintVisitor) say(s string) {
	fmt.Print(s)
}

func (this *PrettyPrintVisitor) printSpeaces() {
	i := this.indentLevel
	for ; i != 0; i-- {
		fmt.Printf(" ")
	}
}

func (this *PrettyPrintVisitor) visitDec(e Dec) {
	switch v := e.(type) {
	case *VarDec:
		this.say("var ")
		this.say(v.Name + " ")
		v.Tp.accept(this)
	case *FieldDec:
		this.say(v.Name + " ")
		v.Tp.accept(this)
	case *StructDec:
		this.say("type " + v.Name + " struct {")
		this.indent()
		for _, fdec := range v.Fields {
			fdec.accept(this)
			this.sayln()
		}
		this.unIndent()
		this.sayln()
		this.say("}")
	}
}
func (this *PrettyPrintVisitor) visitFunc(e Func) {
	switch v := e.(type) {
	case *MethodSingle:
		this.say("func (" + v.Firstarg + " *")
		v.BindingType.accept(this)
		this.say(") " + v.Methodname + "(")
		commer := ""
		for _, vv := range v.Formals {
			this.say(commer)
			commer = ","
			vv.accept(this)
		}
		this.say(") ")
		v.RetType.accept(this)
		this.say(" {\n")
		this.indent()
		for _, vardec := range v.VarDecs {
			this.printSpeaces()
			vardec.accept(this)
			this.sayln()
		}
		for _, stm := range v.Stms {
			this.printSpeaces()
			stm.accept(this)
			this.sayln()
		}

		this.say("return ")
		v.RetExp.accept(this)
		this.unIndent()
		this.sayln()
		this.say("}\n")

	default:
		_, filename, line, _ := runtime.Caller(0)
		util.Bug("Func need MethodSingle", filename, line)
	}
}
func (this *PrettyPrintVisitor) visitMain(e MainFunc) {
	switch v := e.(type) {
	case *MainFuncSingle:
		this.say("func main() {")
		this.indent()
		this.sayln()
		v.PrintStm.accept(this)
		this.unIndent()
		this.sayln()
		this.say("}\n")
	default:
	}
}

func (this *PrettyPrintVisitor) visitProg(e Prog) {
	switch v := e.(type) {
	case *ProgramSingle:
		this.sayln()
		for _, vv := range v.VarDecs {
			vv.accept(this)
			this.sayln()
		}
		for _, vv := range v.StrDecs {
			vv.accept(this)
			this.sayln()
		}
		v.Mfunc.accept(this)
        this.sayln()
		for _, vv := range v.Methods {
			vv.accept(this)
			this.sayln()
		}
	default:
		_, filename, line, _ := runtime.Caller(0)
		util.Bug("need ProgramSingle", filename, line)
	}
}
func (this *PrettyPrintVisitor) visitExp(e Exp) {
	switch v := e.(type) {
	case *Add:
		v.Left.accept(this)
		this.say("+")
		v.Right.accept(this)
	case *And:
		v.Left.accept(this)
		this.say("&&")
		v.Right.accept(this)
	case *Times:
		v.Left.accept(this)
		this.say("*")
		v.Right.accept(this)
	case *ArraySelect:
		v.ArrayName.accept(this)
		this.say("[")
		v.Index.accept(this)
		this.say("]")
	case *Call:
		v.Callee.accept(this)
		this.say("." + v.MethodName + "(")
		for s, arg := range v.ArgsList {
			arg.accept(this)
			if s != len(v.ArgsList)-1 {
				this.say(",")
			}
		}
		this.say(")")
	case *False:
		this.say("false")
	case *True:
		this.say("true")
	case *Id:
		this.say(v.Name)
	case *Len:
		this.say("len(")
		v.Arrayref.accept(this)
		this.say(")")
	case *Lt:
		v.Left.accept(this)
		this.say("<")
		v.Right.accept(this)
	case *NewIntArray:
		this.say("make([]int ")
		v.Size.accept(this)
		this.say(")")
	case *NewObject:
		this.say("new(" + v.Name + ")")
	case *Not:
		this.say("!")
		v.E.accept(this)
	case *Num:
		this.say(strconv.Itoa(v.Value))
	case *Sub:
		v.Left.accept(this)
		this.say("-")
		v.Right.accept(this)
	default:
		_, filename, line, _ := runtime.Caller(0)
		util.Bug("Exp type error", filename, line)

	}
}
func (this *PrettyPrintVisitor) visitStm(e Stm) {
	switch v := e.(type) {
	case *Assign:
		this.say(v.Name)
		this.say("=")
		v.E.accept(this)
	case *Derive:
		this.say(v.Name)
		this.say(":=")
		v.E.accept(this)
	case *AssignArray:
		this.say(v.Name + "[")
		v.Index.accept(this)
		this.say("]=")
		v.E.accept(this)
	case *Block:
		this.say("{")
		this.indent()
		this.sayln()
		for s, stm := range v.Stms {
			stm.accept(this)
			if s != len(v.Stms)-1 {
				this.sayln()
			}
		}
		this.unIndent()
		this.sayln()
		this.say("}")
	case *If:
		this.say("if ")
		v.Condition.accept(this)
		v.Thenn.accept(this)
		this.say("else")
		v.Elsee.accept(this)
	case *Print:
		this.say("fmt.Println(")
		v.E.accept(this)
		this.say(")")
	case *For:
		this.say("for ")
		v.Condition.accept(this)
		v.Body.accept(this)
	}
}

func (this *PrettyPrintVisitor) visitType(e Type) {
	switch v := e.(type) {
	case *Int:
		this.say("int")
	case *Boolean:
		this.say("bool")
	case *IntArray:
		this.say("[]int")
	case *ClassType:
		this.say(v.Name)
	}
}

/**
* Implements the Visitor
 */
func (this *PrettyPrintVisitor) visit(e Acceptable) {
	switch v := e.(type) {
	case Dec:
		this.visitDec(v)
	case Func:
		this.visitFunc(v)
	case MainFunc:
		this.visitMain(v)
	case Prog:
		this.visitProg(v)
	case Exp:
		this.visitExp(v)
	case Stm:
		this.visitStm(v)
	case Type:
		this.visitType(v)
	}
}

func (this *PrettyPrintVisitor) DumpProg(e Acceptable) {
	this.visit(e)
}
