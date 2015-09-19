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
	Name string
}

func (this *ClassType) _type()  {}
func (this *ClassType) accept() {}
func (this *ClassType) GetType() int {
	return TYPE_CLASSTYPE
}
func (this *ClassType) String() string {
	return "@" + this.Name
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
	Left  Exp
	Right Exp
}

func (this *Add) _exp()   {}
func (this *Add) accept() {}

type And struct {
	Left  Exp
	Right Exp
}

func (this *And) _exp()   {}
func (this *And) accept() {}

type ArraySelect struct {
	Arrayref Exp
	Index    Exp
}

func (this *ArraySelect) _exp()   {}
func (this *ArraySelect) accept() {}

type Call struct {
	New_id  string //callee id new Sub().Name -> x_0 = new Sub(), x_0.Name(Args)
	E       Exp
	Name    string //method name
	Args    []Exp
	RetType Type
}

func (this *Call) _exp()   {}
func (this *Call) accept() {}

type Id struct {
	Name    string
	IsField bool
}

func (this *Id) _exp()   {}
func (this *Id) accept() {}

type Length struct {
	Arrayref Exp
}

func (this *Length) _exp()   {}
func (this *Length) accept() {}

type Lt struct {
	Left  Exp
	Right Exp
}

func (this *Lt) _exp()   {}
func (this *Lt) accept() {}

type NewIntArray struct {
	E    Exp
	Name string
}

func (this *NewIntArray) _exp()  {}
func (this NewIntArray) accept() {}

type NewObject struct {
	Class_name string //this field is used to name the allocation
	Name       string
}

func (this *NewObject) _exp()   {}
func (this *NewObject) accept() {}

type Not struct {
	E Exp
}

func (this *Not) _exp()   {}
func (this *Not) accept() {}

type Num struct {
	Value int
}

func (this *Num) _exp()   {}
func (this *Num) accept() {}

type Sub struct {
	Left  Exp
	Right Exp
}

func (this *Sub) _exp()   {}
func (this *Sub) accept() {}

type This struct {
}

func (this *This) _exp()   {}
func (this *This) accept() {}

type Times struct {
	Left  Exp
	Right Exp
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
	Name    string
	E       Exp
	IsField bool
}

func (this *Assign) _stm()   {}
func (this *Assign) accept() {}

type AssignArray struct {
	Name    string
	Index   Exp
	E       Exp
	IsField bool
}

func (this *AssignArray) _stm()   {}
func (this *AssignArray) accept() {}

type Block struct {
	Stms []Stm
}

func (this *Block) _stm()   {}
func (this *Block) accept() {}

type If struct {
	Cond  Exp
	Thenn Stm
	Elsee Stm
}

func (this *If) _stm()   {}
func (this *If) accept() {}

type Print struct {
	E Exp
}

func (this *Print) _stm()   {}
func (this *Print) accept() {}

type While struct {
	Cond Exp
	Body Stm
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
	Tp   Type
	Name string
}

func (this *DecSingle) _dec()   {}
func (this *DecSingle) accept() {}
func (this *DecSingle) GetType() int {
	return this.Tp.GetType()
}
func (this *DecSingle) String() string {
	return this.Tp.String() + " " + this.Name
}

type MainMethod interface {
	_maincmethod()
	accept()
}
type MainMethodSingle struct {
	Locals []Dec
	Stms   Stm
}

func (this *MainMethodSingle) _maincmethod() {}
func (this *MainMethodSingle) accept()       {}

type Method interface {
	_method()
	accept()
}

type MethodSingle struct {
	RetType Type
	ClassId string // class name
	Name    string //method name
	Formals []Dec
	Locals  []Dec
	Stms    []Stm
	RetExp  Exp
}

func (this *MethodSingle) _method() {}
func (this *MethodSingle) accept()  {}

type Class interface {
	_class()
	accept()
}
type ClassSingle struct {
	Name string
	Decs []*Tuple
}

func (this *ClassSingle) _class() {}
func (this *ClassSingle) accept() {}

type Vtable interface {
	_vtable()
	accept()
}
type VtableSingle struct {
	Name    string
	Methods []*Ftuple
}

func (this *VtableSingle) _vtable() {}
func (this *VtableSingle) accept()  {}

type Program interface {
	_prog()
	accept()
}
type ProgramC struct {
	Classes    []Class
	Vtables    []Vtable
	Methods    []Method
	Mainmethod MainMethod
}

func (this *ProgramC) _prog()  {}
func (this *ProgramC) accept() {}
