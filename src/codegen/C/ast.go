package codegen_c

/* Type {{{*/

const (
	TYPE_INT = iota
	TYPE_INTARRAY
	TYPE_CLASSTYPE
)

type Type interface {
	accept()
	_type()
	GetType() int
	String() string
}

type ClassType struct {
	id string
}

func (this *ClassType) _type()  {}
func (this *ClassType) accept() {}
func (this *ClassType) GetType() int {
	return TYPE_CLASSTYPE
}
func (this *ClassType) String() string {
	return "@" + this.id
}

type Int struct {
}

func (this *Int) _type()  {}
func (this *Int) accept() {}
func (this *Int) GetType() int {
	return TYPE_INT
}
func (this *Int) String() string {
	return "@int"
}

type IntArray struct {
}

func (this *IntArray) _type()         {}
func (this *IntArray) accept()        {}
func (this *IntArray) GetType() int   { return TYPE_INTARRAY }
func (this *IntArray) String() string { return "@int[]" }

/*}}}*/

/* Exp {{{*/
type Exp interface {
	_exp()
	accept()
}

type Add struct {
	left  Exp
	right Exp
}

func (this *Add) _exp()   {}
func (this *Add) accept() {}

type And struct {
	left  Exp
	right Exp
}

func (this *And) _exp()   {}
func (this *And) accept() {}

type ArraySelect struct {
	array Exp
	index Exp
}

func (this *ArraySelect) _exp()   {}
func (this *ArraySelect) accept() {}

type Call struct {
	assign string
	e      Exp
	name   string
	args   []Exp
}

func (this *Call) _exp()   {}
func (this *Call) accept() {}

type Id struct {
	name    string
	isField bool
}

func (this *Id) _exp()   {}
func (this *Id) accept() {}

type Length struct {
	array Exp
}

func (this *Length) _exp()   {}
func (this *Length) accept() {}

type Lt struct {
	left  Exp
	right Exp
}

func (this *Lt) _exp()   {}
func (this *Lt) accept() {}

type NewIntArray struct {
	e    Exp
	name string
}

func (this *NewIntArray) _exp()  {}
func (this NewIntArray) accept() {}

type NewObject struct {
	id   string
	name string
}

func (this *NewObject) _exp()   {}
func (this *NewObject) accept() {}

type Not struct {
	e Exp
}

func (this *Not) _exp()   {}
func (this *Not) accept() {}

type Num struct {
	num int
}

func (this *Num) _exp()   {}
func (this *Num) accept() {}

type Sub struct {
	left  Exp
	right Exp
}

func (this *Sub) _exp()   {}
func (this *Sub) accept() {}

type This struct {
}

func (this *This) _exp()   {}
func (this *This) accept() {}

type Times struct {
	left  Exp
	right Exp
}

func (this *Times) _exp()   {}
func (this *Times) accept() {}

/*}}}*/

/* Stm {{{*/
type Stm interface {
	_stm()
	accept()
}

type Assign struct {
	id      string
	e       Exp
	isField bool
}

func (this *Assign) _stm()   {}
func (this *Assign) accept() {}

type AssignArray struct {
	id      string
	index   Exp
	e       Exp
	isField bool
}

func (this *AssignArray) _stm()   {}
func (this *AssignArray) accept() {}

type Block struct {
	stms []Stm
}

func (this *Block) _stm()   {}
func (this *Block) accept() {}

type If struct {
	cond  Exp
	thenn Stm
	elsee Stm
}

func (this *If) _stm()   {}
func (this *If) accept() {}

type Print struct {
	e Exp
}

func (this *Print) _stm()   {}
func (this *Print) accept() {}

type While struct {
	cond Exp
	body Stm
}

func (this *While) _stm()   {}
func (this *While) accept() {}

/*}}}*/

type Dec interface {
	_dec()
	accept()
	GetType() int
	String() string
}
type DecSingle struct {
	tp Type
	id string
}

func (this *DecSingle) _dec()   {}
func (this *DecSingle) accept() {}
func (this *DecSingle) GetType() int {
	return this.tp.GetType()
}
func (this *DecSingle) String() string {
	return this.tp.String() + " " + this.id
}

type MainMethod interface {
	_maincmethod()
	accept()
}
type MainMethodSingle struct {
	locals []Dec
	stm    Stm
}

func (this *MainMethodSingle) _maincmethod() {}
func (this *MainMethodSingle) accept()       {}

type Method interface {
	_method()
	accept()
}

type MethodSingle struct {
	retType Type
	classId string // class name
	id      string //method name
	formals []Dec
	locals  []Dec
	stms    []Stm
	retExp  Exp
}

func (this *MethodSingle) _method() {}
func (this *MethodSingle) accept()  {}

type Class interface {
	_class()
	accept()
}
type ClassSingle struct {
	id   string
	decs []*Tuple
}

func (this *ClassSingle) _class() {}
func (this *ClassSingle) accept() {}

type Vtable interface {
	_vtable()
	accept()
}
type VtableSingle struct {
	id      string
	methods []*Ftuple
}

func (this *VtableSingle) _vtable() {}
func (this *VtableSingle) accept()  {}

type Program interface {
	_prog()
	accept()
}
type ProgramC struct {
	classes    []Class
	vtables    []Vtable
	methods    []Method
	mainMethod MainMethod
}

func (this *ProgramC) _prog()  {}
func (this *ProgramC) accept() {}
