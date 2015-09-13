package codegen_c

import (
	"../../ast"
	"../../util"
)

var table *ClassTable
var classId string
var type_c Type
var dec_c Dec
var stm_c Stm
var exp_c Exp
var method_c Method
var tmpVars_c []Dec
var classes_c []Class
var vtables_c []Vtable
var methods_c []Method

var main_method_c MainMethod
var prog_c Program

func genId() string {
	return util.Next()
}

func trans_Method(m ast.Method) {
	if mm, ok := m.(*ast.MethodSingle); ok {
		tmpVars_c = make([]Dec, 0)
		trans(mm.RetType)
		new_retType := type_c
		var new_formals []Dec
		new_formals = append(new_formals,
			&DecSingle{&ClassType{classId}, "this"})
		for _, dec := range mm.Formals {
			trans(dec)
			new_formals = append(new_formals, dec_c)
		}
		var new_locals []Dec
		for _, dec := range mm.Locals {
			trans(dec)
			new_locals = append(new_locals, dec_c)
		}
		var new_stms []Stm
		for _, s := range mm.Stms {
			trans(s)
			new_stms = append(new_stms, stm_c)
		}
		trans(mm.RetExp)
		retExp := exp_c
		for _, dec := range tmpVars_c {
			new_locals = append(new_locals, dec)
		}
		method_c = &MethodSingle{new_retType,
			classId, mm.Name, new_formals,
			new_locals, new_stms, retExp}
	} else {
		panic("impossible")
	}
}

func trans_Exp_Add(e *ast.Add) {
	trans(e.Left)
	left := exp_c
	trans(e.Right)
	right := exp_c
	exp_c = &Add{left, right}
}

func trans_Exp_And(e *ast.And) {
	trans(e.Left)
	left := exp_c
	trans(e.Right)
	right := exp_c
	exp_c = &And{left, right}
}

func trans_Exp_ArraySelect(e *ast.ArraySelect) {
	trans(e.Arrayref)
	array := exp_c
	trans(e.Index)
	index := exp_c
	exp_c = &ArraySelect{array, index}
}

func trans_Exp_Call(e *ast.Call) {
	trans(e.Callee)
	new_id := genId()
	tmpVars_c = append(tmpVars_c, &DecSingle{&ClassType{e.Firsttype}, new_id})
	exp := exp_c
	args := make([]Exp, 0)
	for _, x := range e.ArgsList {
		trans(x)
		args = append(args, exp_c)
	}
	exp_c = &Call{new_id, exp, e.MethodName, args}
}

func trans_Exp_False(e *ast.False) {
	exp_c = &Num{0}
}

func trans_Exp_Id(e *ast.Id) {
	exp_c = &Id{e.Name, e.IsField}
}

func trans_Exp_Length(e *ast.Length) {
	type_c = &IntArray{}
	trans(e.Arrayref)
	array := exp_c
	exp_c = &Length{array}
}

func trans_Exp_Lt(e *ast.Lt) {
	trans(e.Left)
	left := exp_c
	trans(e.Right)
	right := exp_c
	exp_c = &Lt{left, right}
}

func trans_Exp_NewIntArray(e *ast.NewIntArray) {
	trans(e.Size)
	size := exp_c
	exp_c = &NewIntArray{size, ""}
}

func trans_Exp_NewObject(e *ast.NewObject) {
	exp_c = &NewObject{e.Name, ""}
}

func trans_Exp_Not(e *ast.Not) {
	trans(e.E)
	t := exp_c
	exp_c = &Not{t}
}

func trans_Exp_Num(e *ast.Num) {
	exp_c = &Num{e.Value}
}

func trans_Exp_Sub(e *ast.Sub) {
	trans(e.Left)
	left := exp_c
	trans(e.Right)
	right := exp_c
	exp_c = &Sub{left, right}
}

func trans_Exp_This(e *ast.This) {
	exp_c = &This{}
}

func trans_Exp_Times(e *ast.Times) {
	trans(e.Left)
	left := exp_c
	trans(e.Right)
	right := exp_c
	exp_c = &Times{left, right}
}

func trans_Exp_True(e *ast.True) {
	exp_c = &Num{1}
}

func trans_Exp(e ast.Exp) {
	switch v := e.(type) {
	case *ast.Add:
		trans_Exp_Add(v)
	case *ast.And:
		trans_Exp_And(v)
	case *ast.ArraySelect:
		trans_Exp_ArraySelect(v)
	case *ast.Call:
	case *ast.False:
		trans_Exp_False(v)
	case *ast.Id:
		trans_Exp_Id(v)
	case *ast.Length:
		trans_Exp_Length(v)
	case *ast.Lt:
		trans_Exp_Lt(v)
	case *ast.NewIntArray:
		trans_Exp_NewIntArray(v)
	case *ast.NewObject:
		trans_Exp_NewObject(v)
	case *ast.Not:
		trans_Exp_Not(v)
	case *ast.Num:
		trans_Exp_Num(v)
	case *ast.Sub:
		trans_Exp_Sub(v)
	case *ast.This:
		trans_Exp_This(v)
	case *ast.Times:
		trans_Exp_Times(v)
	case *ast.True:
		trans_Exp_True(v)
	default:
		panic("impossible")
	}
}

func trans_Stm_Assign(s *ast.Assign) {
	isField := s.IsField
	trans(s.E)
	stm_c = &Assign{s.Name, exp_c, isField}
}

func trans_Stm_AssignArray(s *ast.AssignArray) {
	isField := s.IsField
	trans(s.Index)
	index := exp_c
	trans(s.E)
	stm_c = &AssignArray{s.Name, index, exp_c, isField}
}

func trans_Stm_Block(s *ast.Block) {
	var stms []Stm
	for _, ss := range s.Stms {
		trans(ss)
		stms = append(stms, stm_c)
	}
	stm_c = &Block{stms}
}

func trans_Stm_If(s *ast.If) {
	trans(s.Condition)
	cond := exp_c
	trans(s.Thenn)
	thenn := stm_c
	trans(s.Elsee)
	elsee := stm_c
	stm_c = &If{cond, thenn, elsee}
}

func trans_Stm_Print(s *ast.Print) {
	trans(s.E)
	stm_c = &Print{exp_c}
}

func trans_Stm_While(s *ast.While) {
	trans(s.E)
	cond := exp_c
	trans(s.Body)
	body := stm_c
	stm_c = &While{cond, body}
}

func trans_Stm(s ast.Stm) {
	switch v := s.(type) {
	case *ast.Assign:
		trans_Stm_Assign(v)
	case *ast.AssignArray:
		trans_Stm_AssignArray(v)
	case *ast.Block:
		trans_Stm_Block(v)
	case *ast.If:
		trans_Stm_If(v)
	case *ast.Print:
		trans_Stm_Print(v)
	case *ast.While:
		trans_Stm_While(v)
	default:
		panic("impossible")
	}
}

func trans_Type_Int(t *ast.Int) {
	type_c = &Int{}
}

func trans_Type_IntArray(t *ast.IntArray) {
	type_c = &IntArray{}
}

func trans_Type_ClassType(t *ast.ClassType) {
	type_c = &ClassType{t.Name}
}

func trans_Type_Boolean(t *ast.Boolean) {
	type_c = &Int{}
}

func trans_Type(t ast.Type) {
	switch v := t.(type) {
	case *ast.Int:
		trans_Type_Int(v)
	case *ast.IntArray:
		trans_Type_IntArray(v)
	case *ast.ClassType:
		trans_Type_ClassType(v)
	case *ast.Boolean:
		trans_Type_Boolean(v)
	default:
		panic("impossible")
	}
}

func trans_Dec(d ast.Dec) {
	if dec, ok := d.(*ast.DecSingle); ok {
		trans(dec.Tp)
		dec_c = &DecSingle{type_c, dec.Name}
	} else {
		panic("impossible")
	}
}

func trans_Class(c ast.Class) {
	if cc, ok := c.(*ast.ClassSingle); ok {
		cb := table.get(cc.Name)
		classes_c = append(classes_c,
			&ClassSingle{cc.Name, cb.fields})
		vtables_c = append(vtables_c,
			&VtableSingle{cc.Name, cb.methods})
		classId = cc.Name
		for _, m := range cc.Methods {
			trans(m)
			methods_c = append(methods_c, method_c)
		}
	} else {
		panic("impossible")
	}
}

func trans_MainClass(mc ast.MainClass) {
	if cc, ok := mc.(*ast.MainClassSingle); ok {
		cb := table.get(cc.Name)
		newclass := &ClassSingle{cc.Name, cb.fields}
		classes_c = append(classes_c, newclass)
		vtables_c = append(vtables_c,
			&VtableSingle{cc.Name, cb.methods})
		tmpVars_c = make([]Dec, 0)
		trans(cc.Stms)
		mtd := &MainMethodSingle{tmpVars_c, stm_c}
		main_method_c = mtd
	} else {
		panic("impossible")
	}
}

func trans_Program(p ast.Program) {
	transC_init()
	if pp, ok := p.(*ast.ProgramSingle); ok {
		scanProgram(pp)
		trans(pp.Mainclass)
		for _, c := range pp.Classes {
			trans(c)
		}
		prog_c = &ProgramC{classes_c, vtables_c,
			methods_c, main_method_c}
	} else {
		panic("impossible")
	}
}

func trans(e ast.Acceptable) {
	switch v := e.(type) {
	case ast.Exp:
		trans_Exp(v)
	case ast.Stm:
		trans_Stm(v)
	case ast.Program:
		trans_Program(v)
	case ast.Class:
		trans_Class(v)
	case ast.MainClass:
		trans_MainClass(v)
	case ast.Method:
		trans_Method(v)
	case ast.Dec:
		trans_Dec(v)
	case ast.Type:
		trans_Type(v)
	default:
		panic("impossible")
	}
}

func scanClasses(c []ast.Class) {
	for _, cc := range c {
		if class, ok := cc.(*ast.ClassSingle); ok {
			table.init(class.Name, class.Extends)
		} else {
			panic("need *ast.ClassSingle")
		}
	}
	for _, cc := range c {
		var cs *ast.ClassSingle
		var new_decs []Dec
		if class, ok := cc.(*ast.ClassSingle); ok {
			cs = class
		} else {
			panic("need *ast.ClassSingle")
		}
		// decls
		for _, dec := range cs.Decs {
			trans(dec)
			new_decs = append(new_decs, dec_c)
		}
		table.initDecs(cs.Name, new_decs)
		//methods
		for _, mtd := range cs.Methods {
			if m, ok := mtd.(*ast.MethodSingle); ok {
				var new_args []Dec
				new_args = append(new_args, &DecSingle{&ClassType{cs.Name}, "this"})
				for _, arg := range m.Formals {
					trans(arg)
					new_args = append(new_args, dec_c)
				}
				trans(m.RetType)
				new_retType := type_c
				table.initMethod(cs.Name, new_retType, new_args, m.Name)
			} else {
				panic("impossible")
			}
		}
	}
	//calculate all inheritance information
	for _, cc := range c {
		if cs, ok := cc.(*ast.ClassSingle); ok {
			table.inherit(cs.Name)
		}
	}
}

func scanMain(m ast.MainClass) {
	if mm, ok := m.(*ast.MainClassSingle); ok {
		table.init(mm.Name, "")
	}
}

func scanProgram(p ast.Program) {
	if pp, ok := p.(*ast.ProgramSingle); ok {
		scanMain(pp.Mainclass)
		scanClasses(pp.Classes)
	} else {
		panic("impossible")
	}
}

func transC_init() {
	table = ClassTable_new()
	classes_c = make([]Class, 0)
	vtables_c = make([]Vtable, 0)
	methods_c = make([]Method, 0)
}

func CodegenC(p ast.Program) Program {
	trans(p)
	return prog_c
}
