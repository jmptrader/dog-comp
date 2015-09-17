package ast_opt

import (
	"../../ast"
	"../../util"
	"container/list"
)

type DeadClass struct {
	worklist *list.List
	classes  *util.HashSet
}

func DeadClass_new() *DeadClass {
	o := new(DeadClass)
	o.worklist = list.New()
	o.classes = util.HashSet_new()
	return o
}

func (this *DeadClass) opt_Exp(e ast.Exp) {
	switch v := e.(type) {
	case *ast.Add:
		this.opt(v.Left)
		this.opt(v.Right)
	case *ast.And:
		this.opt(v.Left)
		this.opt(v.Right)
	case *ast.ArraySelect:
		this.opt(v.Arrayref)
		this.opt(v.Index)
	case *ast.Call:
		this.opt(v.Callee)
		for _, arg := range v.ArgsList {
			this.opt(arg)
		}
		//ret Type??
	case *ast.False:
	case *ast.Id:
	case *ast.Length:
		this.opt(v.Arrayref)
	case *ast.Lt:
		this.opt(v.Left)
		this.opt(v.Right)
	case *ast.NewIntArray:
		this.opt(v.Size)
	case *ast.NewObject:
		if !this.classes.Contains(v.Name) {
			this.worklist.PushBack(v.Name)
			this.classes.Add(v.Name)
		}
	case *ast.Not:
		this.opt(v.E)
	case *ast.Num:
	case *ast.Sub:
		this.opt(v.Left)
		this.opt(v.Right)
	case *ast.This:
	case *ast.Times:
		this.opt(v.Left)
		this.opt(v.Right)
	case *ast.True:
	default:
		panic("impossible")
	}
}

func (this *DeadClass) opt_Stm(stm ast.Stm) {
	switch s := stm.(type) {
	case *ast.Assign:
		this.opt(s.E)
	case *ast.AssignArray:
		this.opt(s.Index)
		this.opt(s.E)
	case *ast.Block:
		for _, ss := range s.Stms {
			this.opt(ss)
		}
	case *ast.If:
		this.opt(s.Condition)
		this.opt(s.Thenn)
		this.opt(s.Elsee)
	case *ast.Print:
		this.opt(s.E)
	case *ast.While:
		this.opt(s.E)
		this.opt(s.Body)
	default:
		panic("impossible")

	}
}

func (this *DeadClass) opt_Method(method ast.Method) {
	switch m := method.(type) {
	case *ast.MethodSingle:
		//omit the formals and locals
		//Statements
		for _, s := range m.Stms {
			this.opt(s)
		}
		//return exp
		this.opt(m.RetExp)
	default:
		panic("impossible")
	}
}

func (this *DeadClass) opt_Class(class ast.Class) {
	switch c := class.(type) {
	case *ast.ClassSingle:
		//Add super class into worklist and set
		if c.Extends != "" {
			exist := false
			for e := this.worklist.Front(); e != nil; e = e.Next() {
				classid := e.Value
				if classid == c.Extends {
					exist = true
					break
				}
			}
			if exist == false {
				this.worklist.PushBack(c.Extends)
				this.classes.Add(c.Extends)
			}
		}
		//methods
		for _, m := range c.Methods {
			this.opt(m)
		}
	default:
		panic("impossible")
	}
}

func (this *DeadClass) opt_MainClass(c ast.MainClass) {
	switch cc := c.(type) {
	case *ast.MainClassSingle:
		this.classes.Add(cc.Name)
		this.opt(cc.Stms)
	default:
		panic("impossible")
	}
}

func (this *DeadClass) opt(e ast.Acceptable) {
	switch v := e.(type) {
	case ast.Exp:
		this.opt_Exp(v)
	case ast.Stm:
		this.opt_Stm(v)
	case ast.Type:
		//no need
	case ast.Dec:
		//no need
	case ast.Method:
		this.opt_Method(v)
	case ast.MainClass:
		this.opt_MainClass(v)
	case ast.Class:
		this.opt_Class(v)
	default:
		panic("impossible")
	}
}

func (this *DeadClass) DeadClass_Opt(prog ast.Program) ast.Program {
	var p *ast.ProgramSingle
	if v, ok := prog.(*ast.ProgramSingle); ok {
		p = v
	} else {
		panic("impossible")
	}
	//add mainclass to the classes set
	this.opt(p.Mainclass)

	for this.worklist.Len() != 0 {
		e := this.worklist.Front()
		this.worklist.Remove(e)
		classid := e.Value
		for _, c := range p.Classes {
			if cc, ok := c.(*ast.ClassSingle); ok {
				if cc.Name == classid {
					this.opt(c)
					break
				}
			} else {
				panic("impossible")
			}
		}
	}

	newclasses := make([]ast.Class, 0)
	for _, c := range p.Classes {
		if cc, ok := c.(*ast.ClassSingle); ok {
			if this.classes.Contains(cc.Name) {
				newclasses = append(newclasses, c)
			}
		} else {
			panic("impossible")
		}
	}

	Ast := &ast.ProgramSingle{p.Mainclass, newclasses}

	return Ast
}
