package cfg

import (
	"../util"
	"strconv"
)

type Acceptable interface {
	accept()
}

//Block
type Block interface {
	accept()
	_block()
}
type BlockSingle struct {
	Label_id *util.Label
	Stms     []Stm
	Trans    Transfer
}

func (this *BlockSingle) accept() {}
func (this *BlockSingle) _block() {}

//Class
type Class interface {
	accept()
	_class()
}

type ClassSingle struct {
	Name string
	Decs []*Tuple
}

func (this *ClassSingle) accept() {}
func (this *ClassSingle) _class() {}

//Dec
type Dec interface {
	accept()
	_dec()
}

type DecSingle struct {
	Tp   Type
	Name string
}

func (this *DecSingle) accept() {}
func (this *DecSingle) _dec()   {}

//MainMethod
type MainMethod interface {
	accept()
	_mainclass()
}

type MainMethodSingle struct {
	Locals []Dec
	Blocks []Block
}

func (this *MainMethodSingle) accept()     {}
func (this *MainMethodSingle) _mainclass() {}

//Method
type Method interface {
	accept()
	_method()
}

type MethodSingle struct {
	Ret_type Type
	Name     string
	ClassId  string
	Formals  []Dec
	Locals   []Dec
	Blocks   []Block
	Entry    *util.Label
}

func (this *MethodSingle) accept()  {}
func (this *MethodSingle) _method() {}

//Operand
type Operand interface {
	accept()
	_operand()
	String() string
}

type Int struct {
	Value int
}

func (this *Int) accept()   {}
func (this *Int) _operand() {}
func (this *Int) String() string {
	return strconv.Itoa(this.Value)
}

type Var struct {
	Name    string
	IsField bool
}

func (this *Var) accept()   {}
func (this *Var) _operand() {}
func (this *Var) String() string {
	return this.Name
}

//prog
type Program interface {
	accept()
	_prog()
}

type ProgramSingle struct {
	Classes     []Class
	Vtables     []Vtable
	Methods     []Method
	Main_method MainMethod
}

func (this *ProgramSingle) accept() {}
func (this *ProgramSingle) _prog()  {}

//Stm
type Stm interface {
	accept()
	_stm()
	String() string
}

type Add struct {
	Dst   string
	Tp    Type
	Left  Operand
	Right Operand
}

func (this *Add) accept() {}
func (this *Add) _stm()   {}
func (this *Add) String() string {
	return this.Dst + " = " + this.Left.String() + " + " +
		this.Right.String() + ";"
}

type And struct {
	Dst   string
	Left  Operand
	Right Operand
}

func (this *And) accept() {}
func (this *And) _stm()   {}
func (this *And) String() string {
	return this.Dst + " = " + this.Left.String() + " && " +
		this.Right.String() + ";"
}

type ArraySelect struct {
	Name     string
	Arrayref Operand
	Index    Operand
}

func (this *ArraySelect) accept() {}
func (this *ArraySelect) _stm()   {}
func (this *ArraySelect) String() string {
	return this.Name + " = " + this.Arrayref.String() + "[" + this.Index.String() +
		"+4];"
}

type AssignArray struct {
	Dst     string
	Index   Operand
	E       Operand
	IsField bool
}

func (this *AssignArray) accept() {}
func (this *AssignArray) _stm()   {}
func (this *AssignArray) String() string {
	return this.Dst + "[" + this.Index.String() + "+4]=" + this.E.String() + ";"
}

type InvokeVirtual struct {
	Dst  string
	Obj  string
	F    string
	Args []Operand //type of the destination variable
}

func (this *InvokeVirtual) accept() {}
func (this *InvokeVirtual) _stm()   {}
func (this *InvokeVirtual) String() string {
	return this.Dst + " = " + this.Obj + "->vptr->" +
		this.F + "(" + this.Obj + " ... );"
}

type Length struct {
	Dst      string
	Arrayref Operand
}

func (this *Length) accept() {}
func (this *Length) _stm()   {}
func (this *Length) String() string {
	return this.Dst + " = " + this.Arrayref.String() + "[2];"
}

type Lt struct {
	Dst   string
	Tp    Type
	Left  Operand
	Right Operand
}

func (this *Lt) accept() {}
func (this *Lt) _stm()   {}
func (this *Lt) String() string {
	return this.Dst + " = " + this.Left.String() + " < " +
		this.Right.String() + ";"
}

type Move struct {
	Dst     string
	Tp      Type
	Src     Operand
	IsField bool
}

func (this *Move) accept() {}
func (this *Move) _stm()   {}
func (this *Move) String() string {
	return this.Dst + " = " + this.Src.String()
}

type NewIntArray struct {
	Dst string
	E   Operand
}

func (this *NewIntArray) accept() {}
func (this *NewIntArray) _stm()   {}
func (this *NewIntArray) String() string {
	return this.Dst + "= (int*)Tiger_new_array(" +
		this.E.String() + ")"
}

type NewObject struct {
	Dst        string
	Class_name string
}

func (this *NewObject) accept() {}
func (this *NewObject) _stm()   {}
func (this *NewObject) String() string {
	return this.Dst + " = " + this.Class_name + ";"
}

type Not struct {
	Dst string
	E   Operand
}

func (this *Not) accept() {}
func (this *Not) _stm()   {}
func (this *Not) String() string {
	return this.Dst + " = !" + this.E.String()
}

type Print struct {
	Args Operand
}

func (this *Print) accept() {}
func (this *Print) _stm()   {}
func (this *Print) String() string {
	return "System_out_printtln(" + this.Args.String() + ");"
}

type Sub struct {
	Dst   string
	Tp    Type
	Left  Operand
	Right Operand
}

func (this *Sub) accept() {}
func (this *Sub) _stm()   {}
func (this *Sub) String() string {
	return this.Dst + " = " + this.Left.String() + " - " +
		this.Right.String() + ";"
}

type Times struct {
	Dst   string
	Tp    Type
	Left  Operand
	Right Operand
}

func (this *Times) accept() {}
func (this *Times) _stm()   {}
func (this *Times) String() string {
	return this.Dst + " = " + this.Left.String() + " * " +
		this.Right.String() + ";"
}

//Transfer
type Transfer interface {
	accept()
	_transfer()
	String() string
}

type Goto struct {
	Label_id *util.Label
}

func (this *Goto) accept()    {}
func (this *Goto) _transfer() {}
func (this *Goto) String() string {
	return "goto " + this.Label_id.String()
}

type If struct {
	Cond   Operand
	Truee  *util.Label
	Falsee *util.Label
}

func (this *If) accept()    {}
func (this *If) _transfer() {}
func (this *If) String() string {
	return "if " + this.Cond.String() + " " + "goto " + this.Truee.String() +
		" else " + "goto " + this.Falsee.String()
}

type Return struct {
	Op Operand
}

func (this *Return) accept()    {}
func (this *Return) _transfer() {}
func (this *Return) String() string {
	return "return " + this.Op.String() + ";"
}

//Type
const (
	TYPE_INT = iota
	TYPE_INTARRAY
	TYPE_CLASSTYPE
)

type Type interface {
	accept()
	_type()
	GetType() int
}

type ClassType struct {
	Name string
}

func (this *ClassType) accept()      {}
func (this *ClassType) _type()       {}
func (this *ClassType) GetType() int { return TYPE_CLASSTYPE }

type IntType struct {
}

func (this *IntType) accept()      {}
func (this *IntType) _type()       {}
func (this *IntType) GetType() int { return TYPE_INT }

type IntArrayType struct {
}

func (this *IntArrayType) accept()      {}
func (this *IntArrayType) _type()       {}
func (this *IntArrayType) GetType() int { return TYPE_INTARRAY }

//Vtable
type Vtable interface {
	accept()
	_vtable()
}

type VtableSingle struct {
	Name    string
	Methods []*Ftuple
}

func (this *VtableSingle) accept()  {}
func (this *VtableSingle) _vtable() {}
