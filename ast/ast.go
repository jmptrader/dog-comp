package ast

import (
//"container/list"
)

type Class interface {
	accept(v Visitor)
	_class()
}

type Dec interface {
	accept(v Visitor)
	_dec()
}

type Exp interface {
	accept(v Visitor)
	_exp()
}

type MainClass interface {
	accept(v Visitor)
	_mainclass()
}

type Method interface {
	accept(v Visitor)
	_method()
}

type Program interface {
	accept(v Visitor)
	_prog()
}

type Stm interface {
	accept(v Visitor)
	_stm()
}

type Type interface {
	accept(v Visitor)
	_type()
}

/*Dec*/ /*{{{*/
type DecSingle struct {
    Tp Type
    Name string
    IsField bool
}

func (this *DecSingle) accept(v Visitor) {
    v.visit(this)
}

func (this *DecSingle) _dec() {
}
/*}}}*/

type MainClassSingle struct {
    Name string
    Args string
    Stms Stm
}

func (this *MainClassSingle) accept(v Visitor) {
    v.visit(this)
}

func (this *MainClassSingle) _mainclass(){
}

/* ClassSingle {{{*/
type ClassSingle struct {
    Name string
    Extends string
    Decs []Dec
    Methods []Method
}
func (this *ClassSingle) accept(v Visitor) {
    v.visit(this)
}
func (this *ClassSingle) _class(){
}
/*}}}*/

//Method  /*{{{*/
type MethodSingle struct {
    RetType Type
    Name string
	Formals     []Dec
    Locals      []Dec
	Stms        []Stm
	RetExp      Exp
}

func (this *MethodSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *MethodSingle) _method() {
}

/*}}}*/

/*Prog*/ /*{{{*/
type ProgramSingle struct {
	Mainclass   MainClass
	Classes []Class
}

func (this *ProgramSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *ProgramSingle) _prog() {
} /*}}}*/

/*Exp*/ /*{{{*/

//Exp.Add /*{{{*/
type Add struct {
	Left  Exp
	Right Exp
}

func (this *Add) accept(v Visitor) {
	v.visit(this)
}
func (this *Add) _exp() {
} /*}}}*/

//Exp.And /*{{{*/
type And struct {
	Left  Exp
	Right Exp
}

func (this *And) accept(v Visitor) {
	v.visit(this)
}
func (this *And) _exp() {
} /*}}}*/

//Exp.Time  /*{{{*/
type Times struct {
	Left  Exp
	Right Exp
}

func (this *Times) accept(v Visitor) {
	v.visit(this)
}
func (this *Times) _exp() {
}

/*}}}*/

//Exp.ArraySelect /*{{{*/
type ArraySelect struct {
	ArrayName Exp
	Index     Exp
}

func (this *ArraySelect) accept(v Visitor) {
	v.visit(this)
}
func (this *ArraySelect) _exp() {
}

/*}}}*/

//Exp.Call /*{{{*/
type Call struct {
	Callee     Exp
	MethodName string
	ArgsList   []Exp
    Firsttype string
	ArgsType   []Type
	Rt         Type
}

func (this *Call) accept(v Visitor) {
	v.visit(this)
}
func (this *Call) _exp() {
}

/*}}}*/

//Exp.False /*{{{*/
type False struct {
}

func (this *False) accept(v Visitor) {
	v.visit(this)
}
func (this *False) _exp() {
}

/*}}}*/

//Exp.True   /*{{{*/
type True struct {
}

func (this *True) accept(v Visitor) {
	v.visit(this)
}
func (this *True) _exp() {
}

/*}}}*/

//Exp.Id /*{{{*/
type Id struct {
	Name string
	Tp   Type
    IsField bool
}

func (this *Id) accept(v Visitor) {
	v.visit(this)
}
func (this *Id) _exp() {
}

/*}}}*/

//Exp.len /*{{{*/
//len(arrayref)
type Length struct {
	Arrayref Exp
}

func (this *Length) accept(v Visitor) {
	v.visit(this)
}
func (this *Length) _exp() {
}

/*}}}*/

//Exp.lt    /*{{{*/
// left < right
type Lt struct {
	Left  Exp
	Right Exp
}

func (this *Lt) accept(v Visitor) {
	v.visit(this)
}
func (this *Lt) _exp() {
}

/*}}}*/

//Exp.NewIntArray   /*{{{*/
type NewIntArray struct {
	Size Exp
}

func (this *NewIntArray) accept(v Visitor) {
	v.visit(this)
}
func (this *NewIntArray) _exp() {
}

/*}}}*/

//Exp.NewObject /*{{{*/
type NewObject struct {
	Name string
}

func (this *NewObject) accept(v Visitor) {
	v.visit(this)
}
func (this *NewObject) _exp() {
}

/*}}}*/

//Exp.Not   /*{{{*/
// !expp
type Not struct {
	E Exp
}

func (this *Not) accept(v Visitor) {
	v.visit(this)
}
func (this *Not) _exp() {
}

/*}}}*/

//Exp.Num   /*{{{*/
type Num struct {
	Value int
}

func (this *Num) accept(v Visitor) {
	v.visit(this)
}
func (this *Num) _exp() {
}

/*}}}*/

//Exp.Sub   /*{{{*/
type Sub struct {
	Left  Exp
	Right Exp
}

func (this *Sub) accept(v Visitor) {
	v.visit(this)
}
func (this *Sub) _exp() {
}

/*}}}*/

type This struct {
}

func (this *This) accept(v Visitor) {
    v.visit(this)
}

func (this *This) _exp() {
}


//Exp end/*}}}*/

//Stm   /*{{{*/

//Stm.Assign    /*{{{*/
type Assign struct {
	Name string
	E    Exp
	Tp   Type
    IsField bool
}

func (this *Assign) accept(v Visitor) {
	v.visit(this)
}
func (this *Assign) _stm() {
}

/*}}}*/


//Stm.AssignArray   /*{{{*/
type AssignArray struct {
	// id[index] = e
	Name  string
	Index Exp
	E     Exp
	Tp    Type
    IsField bool
}

func (this *AssignArray) accept(v Visitor) {
	v.visit(this)
}
func (this *AssignArray) _stm() {
}

/*}}}*/

//Stm.Block /*{{{*/
type Block struct {
	Stms []Stm
}

func (this *Block) accept(v Visitor) {
	v.visit(this)
}
func (this *Block) _stm() {
}

/*}}}*/

//Stm.If    /*{{{*/
type If struct {
	Condition Exp
	Thenn     Stm
	Elsee     Stm
}

func (this *If) accept(v Visitor) {
	v.visit(this)
}
func (this *If) _stm() {
}

/*}}}*/

//Stm.Print /*{{{*/
type Print struct {
	E Exp
}

func (this *Print) accept(v Visitor) {
	v.visit(this)
}
func (this *Print) _stm() {
}

/*}}}*/

//Stm.For   /*{{{*/
type While struct {
    E Exp
    Body Stm
}
func (this *While) accept(v Visitor) {
	v.visit(this)
}
func (this *While) _stm() {
}

/*}}}*/

//Stm end/*}}}*/

//Type   /*{{{*/

//Type.Int  /*{{{*/
type Int struct {
}

func (this *Int) accept(v Visitor) {
	v.visit(this)
}
func (this *Int) _type() {
}

/*}}}*/

//Type.Bool /*{{{*/
type Boolean struct {
}

func (this *Boolean) accept(v Visitor) {
	v.visit(this)
}
func (this *Boolean) _type() {
}

/*}}}*/

//Type.IntArray /*{{{*/
type IntArray struct {
}

func (this *IntArray) accept(v Visitor) {
	v.visit(this)
}
func (this *IntArray) _type() {
}

/*}}}*/

//Type.ClassType    /*{{{*/
type ClassType struct {
	Name string
}

func (this *ClassType) accept(v Visitor) {
	v.visit(this)
}
func (this *ClassType) _type() {
}

/*}}}*/

//Type end/*}}}*/
