package ast

import (
	"container/list"
)

type Struct interface {
	accept(v Visitor)
	_struct()
}

type Dec interface {
	accept(v Visitor)
	_dec()
}

type Exp interface {
	accept(v Visitor)
	_exp()
}

type MainFunc interface {
	accept(v Visitor)
	_mainfunc()
}

type Func interface {
	accept(v Visitor)
	_func()
}

type Prog interface {
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

/*Class-------------------------------------*//*{{{*/
//Class.StructSingle    /*{{{*/
type StructSingle struct {
	id   string
	desc *list.List
	//methods *list.List
}

func (this *StructSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *StructSingle) _struct() {
}
/*}}}*/

/*}}}*/

/*Dec*//*{{{*/

//VarDec  /*{{{*/
type VarDec struct {
	Tp      Type
	Name      string
}
func (this *VarDec) accept(v Visitor) {
	v.visit(this)
}
func (this *VarDec) _dec() {
}
/*}}}*/

//FieldDec /*{{{*/
type FieldDec struct {
	Tp Type
	Id string
}
func (this *FieldDec) accept(v Visitor) {
	v.visit(this)
}
func (this *FieldDec) _dec() {
}
/*}}}*/

//StructDec    /*{{{*/
type StructDec struct {
    Id string
    Fields *list.List
}
func (this *StructDec) accept(v Visitor) {
    v.visit(this)
}
func (this *StructDec) _dec() {
}
/*}}}*/

/*}}}*/

/*MainFunc*/ /*{{{*/
type MainFuncSingle struct {
	id  string
	stm Stm
}

func (this *MainFuncSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *MainFuncSingle) _mainfuc() {
}
/*}}}*/

//Func  /*{{{*/
type MethodSingle struct {
    Firstasg string
    BindingType Type//??? or Type
    Methodname string
    Formals *list.List
    RetType Type
    VarDecs *list.List
    Stms *list.List
    RetExp Exp
}

func (this *MethodSingle) accept(v Visitor) {
    v.visit(this)
}
func (this *MethodSingle) _func() {
}
/*}}}*/

/*Prog*/ /*{{{*/
type ProgramSingle struct {
	Mfunc   MainFunc
    VarDec *list.List
	StrDecs *list.List
	Methods *list.List
}

func (this *ProgramSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *ProgramSingle) _prog() {
}/*}}}*/

/*Exp*/ /*{{{*/

//Exp.Add /*{{{*/
type Add struct {
    Left Exp
    Right Exp
}

func (this *Add) accept(v Visitor) {
    v.visit(this)
}
func (this *Add) _exp() {
}/*}}}*/

//Exp.And /*{{{*/
type And struct {
    Left Exp
    Right Exp
}
func (this *And) accept(v Visitor) {
    v.visit(this)
}
func (this *And) _exp(){
}/*}}}*/

//Exp.Time  /*{{{*/
type Times struct {
    Left Exp
    Right Exp
}
func (this *Times) accept(v Visitor){
    v.visit(this)
}
func (this *Times) _exp(){
}
/*}}}*/

//Exp.ArraySelect /*{{{*/
type ArraySelect struct {
    ArrayName Exp
    Index Exp
}
func (this *ArraySelect) accept(v Visitor){
    v.visit(this)
}
func (this *ArraySelect) _exp(){
}
/*}}}*/

//Exp.Call /*{{{*/
type Call struct {
    Callee Exp
    MethodName string
    ArgsList *list.List
    ArgsType *list.List
    Rt Type
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
    Tp Type
}
func (this *Id) accept(v Visitor) {
    v.visit(this)
}
func (this *Id) _exp() {
}
/*}}}*/

//Exp.len /*{{{*/
//len(arrayref)
type Len struct {
    Arrayref Exp
}
func (this *Len) accept(v Visitor) {
    v.visit(this)
}
func (this *Len) _exp() {
}
/*}}}*/

//Exp.lt    /*{{{*/
// left < right
type Lt struct {
    Left Exp
    Right Exp
}
func (this *Lt) accept(v Visitor) {
    v.visit(this)
}
func (this *Lt) _exp() {
}
/*}}}*/

//Exp.NewIntArray   /*{{{*/
// make([]int, size)
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
// new(id)
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
    Left Exp
    Right Exp
}
func (this *Sub) accept(v Visitor) {
    v.visit(this)
}
func (this *Sub) _exp() {
}
/*}}}*/

//Exp end/*}}}*/

//Stm   /*{{{*/

//Stm.Assign    /*{{{*/
type Assign struct {
    Name string
    E Exp
    Tp Type
}
func (this *Assign) accept(v Visitor) {
    v.visit(this)
}
func (this *Assign) _stm() {
}
/*}}}*/

//Stm.Derive    /*{{{*/
type Derive struct {
    Name string
    E Exp
    Tp Type
}
func (this *Derive) accept(v Visitor) {
    v.visit(this)
}
func (this *Derive) _stm() {
}
/*}}}*/

//Stm.AssignArray   /*{{{*/
type AssignArray struct {
    // id[index] = e
    Name string
    Index Exp
    E Exp
    Tp Type
}
func (this *AssignArray) accept(v Visitor) {
    v.visit(this)
}
func (this *AssignArray) _stm() {
}
/*}}}*/

//Stm.Block /*{{{*/
type Block struct {
    Stms *list.List
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
    Thenn Stm
    Elsee Stm
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
type For struct {
    Condition Exp
    Body Stm
}
func (this *For) accept(v Visitor) {
    v.visit(this)
}
func (this *For) _stm() {
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
}
func (this *ClassType) accept(v Visitor) {
    v.visit(this)
}
func (this *ClassType) _type() {
}

/*}}}*/

//Type end/*}}}*/
