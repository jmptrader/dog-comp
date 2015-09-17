package ast_opt

import (
	"../../ast"
)

type DeadCode struct {
	new_class  ast.Class
	main_class ast.MainClass
	stm        ast.Stm
	classes    []ast.Class
	methods    []ast.Method
	stms       []ast.Stm
	method     ast.Method
	is_bool    bool
	is_true    bool
}

func DeadCode_new() *DeadCode {
	o := new(DeadCode)
	o.classes = make([]ast.Class, 0)
	o.methods = make([]ast.Method, 0)
	o.stms = make([]ast.Stm, 0)

	return o
}

func (this *DeadCode) opt_Exp(exp ast.Exp) {
	switch e := exp.(type) {
	case *ast.Add:
		this.is_bool = false
	case *ast.And:
		this.opt(e.Left)
		left := this.is_true
		this.opt(e.Right)
		right := this.is_true
		this.is_bool = true
		if left && right {
			this.is_true = true
		} else {
			this.is_true = false
		}
	case *ast.ArraySelect:
		this.is_bool = false
	case *ast.Call:
		this.is_bool = false
	case *ast.False:
		this.is_bool = true
		this.is_true = false
	case *ast.Id:
		this.is_bool = false
	case *ast.Length:
		this.is_bool = false
	case *ast.Lt:
		/*
		 * Although we can do some magic in here to opt
		 * Exp like 1<2 -> true, but the real work is in
		 * const-fold.go
		 */
		this.is_bool = false
	case *ast.NewIntArray:
		this.is_bool = false
	case *ast.NewObject:
		this.is_bool = false
	case *ast.Not:
		this.is_bool = true
		this.opt(e.E)
		if this.is_bool {
			this.is_true = !this.is_true
		}
	case *ast.Num:
		this.is_bool = false
	case *ast.Sub:
		this.is_bool = false
	case *ast.This:
		this.is_bool = false
	case *ast.Times:
		this.is_bool = false
	case *ast.True:
		this.is_bool = true
		this.is_true = true
	default:
		panic("impossible")
	}
}

func (this *DeadCode) opt_Stm(stm ast.Stm) {
	switch s := stm.(type) {
	case *ast.Assign:
		this.stm = s
	case *ast.AssignArray:
		this.stm = s
	case *ast.Block:
		for _, ss := range s.Stms {
			this.opt(ss)
		}
	case *ast.If:
		this.opt(s.Condition)
		if this.is_bool {
			if this.is_true {
				this.stm = s.Thenn
			} else {
				this.stm = s.Elsee
			}
		} else {
			this.stm = s
		}
	case *ast.Print:
		this.stm = s
	case *ast.While:
		this.opt(s.E)
		if this.is_bool && !this.is_true {
			this.stm = nil
		} else {
			this.stm = s
		}
	default:
		panic("impossible")
	}
}

func (this *DeadCode) opt_Method(method ast.Method) {
	switch m := method.(type) {
	case *ast.MethodSingle:
		this.stms = make([]ast.Stm, 0)
		for _, s := range m.Stms {
			this.opt(s)
			if this.stm != nil {
				this.stms = append(this.stms, this.stm)
			}
		}
		this.method = &ast.MethodSingle{m.RetType,
			m.Name,
			m.Formals,
			m.Locals,
			this.stms,
			m.RetExp}
	default:
		panic("impossible")
	}
}

func (this *DeadCode) opt_MainClass(mm ast.MainClass) {
	switch c := mm.(type) {
	case *ast.MainClassSingle:
		this.opt(c.Stms)
		this.main_class = &ast.MainClassSingle{c.Name, c.Args, this.stm}
	default:
		panic("impossilbe")
	}
}

func (this *DeadCode) opt_Class(cc ast.Class) {
	switch c := cc.(type) {
	case *ast.ClassSingle:
		this.methods = make([]ast.Method, 0)
		for _, m := range c.Methods {
			this.opt(m)
			this.methods = append(this.methods, this.method)
		}
		this.new_class = &ast.ClassSingle{c.Name, c.Extends, c.Decs, this.methods}
	default:
		panic("impossible")
	}
}

func (this *DeadCode) opt(e ast.Acceptable) {
	switch v := e.(type) {
	case ast.Class:
		this.opt_Class(v)
	case ast.MainClass:
		this.opt_MainClass(v)
	case ast.Method:
		this.opt_Method(v)
	case ast.Stm:
		this.opt_Stm(v)
	case ast.Exp:
		this.opt_Exp(v)
	case ast.Type:
	case ast.Dec:
	default:
		panic("impossible")
	}
}

func (this *DeadCode) DeadCode_Opt(prog ast.Program) ast.Program {
	var Ast *ast.ProgramSingle
	switch p := prog.(type) {
	case *ast.ProgramSingle:
		this.opt(p.Mainclass)
		this.classes = make([]ast.Class, 0)
		for _, c := range p.Classes {
			this.opt(c)
			this.classes = append(this.classes, this.new_class)
		}
		Ast = &ast.ProgramSingle{this.main_class, this.classes}

	default:
		panic("impossible")
	}
	return Ast
}
