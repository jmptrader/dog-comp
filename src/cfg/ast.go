package cfg

import (
	"../util"
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
	label    *util.Label
	stms     []Stm
	transfer Transfer
}

func (this *BlockSingle) accept() {}
func (this *BlockSingle) _block() {}

//Class
type Class interface {
	accept()
	_class()
}

type ClassSingle struct {
	id   string
	decs []*Tuple
}

func (this *ClassSingle) accept() {}
func (this *ClassSingle) _class() {}

//Dec
type Dec interface {
	accept()
	_dec()
}

type DecSingle struct {
	tp Type
	id string
}

func (this *DecSingle) accept() {}
func (this *DecSingle) _dec()   {}

//MainMethod
type MainMethod interface {
	accept()
	_mainclass()
}

type MainMethodSingle struct {
	locals []Dec
	blocks []Block
}

func (this *MainMethodSingle) accept()     {}
func (this *MainMethodSingle) _mainclass() {}

//Method
type Method interface {
	accept()
	_method()
}

type MethodSingle struct {
	ret_type Type
	name     string
	classId  string
	formals  []Dec
	locals   []Dec
	blocks   []Block
	entry    *util.Label
}

func (this *MethodSingle) accept()  {}
func (this *MethodSingle) _method() {}

//Operand
type Operand interface {
	accept()
	_operand()
}

type Int struct {
	value int
}

func (this *Int) accept()   {}
func (this *Int) _operand() {}

type Var struct {
	id      string
	isField bool
}

func (this *Var) accept()   {}
func (this *Var) _operand() {}

//prog
type Program interface {
	accept()
	_prog()
}

type ProgramSingle struct {
	classes     []Class
	vtables     []Vtable
	methods     []Method
	main_method MainMethod
}

func (this *ProgramSingle) accept() {}
func (this *ProgramSingle) _prog()  {}

//Stm
type Stm interface {
	accept()
	_stm()
}

type Add struct {
	dst   string
	tp    Type
	left  Operand
	right Operand
}

func (this *Add) accept() {}
func (this *Add) _stm()   {}

type And struct {
	dst   string
	left  Operand
	right Operand
}

func (this *And) accept() {}
func (this *And) _stm()   {}

type ArraySelect struct {
	id    string
	array Operand
	index Operand
}

func (this *ArraySelect) accept() {}
func (this *ArraySelect) _stm()   {}

type AssignArray struct {
	dst     string
	index   Operand
	exp     Operand
	isField bool
}

func (this *AssignArray) accept() {}
func (this *AssignArray) _stm()   {}

type InvokeVirtual struct {
	dst  string
	obj  string
	f    string
	args []Operand //type of the destination variable
}

func (this *InvokeVirtual) accept() {}
func (this *InvokeVirtual) _stm()   {}

type Length struct {
	dst   string
	array Operand
}

func (this *Length) accept() {}
func (this *Length) _stm()   {}

type Lt struct {
	dst   string
	tp    Type
	left  Operand
	right Operand
}

func (this *Lt) accept() {}
func (this *Lt) _stm()   {}

type Move struct {
	dst     string
	tp      Type
	src     Operand
	IsField bool
}

func (this *Move) accept() {}
func (this *Move) _stm()   {}

type NewIntArray struct {
	dst string
	exp Operand
}

func (this *NewIntArray) accept() {}
func (this *NewIntArray) _stm()   {}

type NewObject struct {
	dst string
	c   string
}

func (this *NewObject) accept() {}
func (this *NewObject) _stm()   {}

type Not struct {
	dst string
	exp Operand
}

func (this *Not) accept() {}
func (this *Not) _stm()   {}

type Print struct {
	arg Operand
}

func (this *Print) accept() {}
func (this *Print) _stm()   {}

type Sub struct {
	dst   string
	tp    Type
	left  Operand
	right Operand
}

func (this *Sub) accept() {}
func (this *Sub) _stm()   {}

type Times struct {
	dst   string
	tp    Type
	left  Operand
	right Operand
}

func (this *Times) accept() {}
func (this *Times) _stm()   {}

//Transfer
type Transfer interface {
	accept()
	_transfer()
}

type Goto struct {
	label *util.Label
}

func (this *Goto) accept()    {}
func (this *Goto) _transfer() {}

type If struct {
	cond   Operand
	truee  *util.Label
	falsee *util.Label
}

func (this *If) accept()    {}
func (this *If) _transfer() {}

type Return struct {
	op Operand
}

func (this *Return) accept()    {}
func (this *Return) _transfer() {}

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
	id string
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
	id      string
	methods []*Ftuple
}

func (this *VtableSingle) accept()  {}
func (this *VtableSingle) _vtable() {}
