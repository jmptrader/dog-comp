package ast_opt

import (
	"../../ast"
	"../../util"
	"container/list"
	"fmt"
)

var worklist *list.List = list.New()
var classes *util.HashSet = util.HashSet_new()

func opt_Exp(e ast.Exp) {
	switch v := e.(type) {
	case *ast.Add:
		opt(v.Left)
		opt(v.Right)
	case *ast.And:
		opt(v.Left)
		opt(v.Right)
	case *ast.ArraySelect:
		opt(v.Arrayref)
		opt(v.Index)
	case *ast.Call:
		opt(v.Callee)
		for _, arg := range v.ArgsList {
			opt(arg)
		}
		//ret Type??
	case *ast.False:
	case *ast.Id:
	case *ast.Length:
		opt(v.Arrayref)
	case *ast.Lt:
		opt(v.Left)
		opt(v.Right)
	case *ast.NewIntArray:
		opt(v.Size)
	case *ast.NewObject:
		if !classes.Contains(v.Name) {
			worklist.PushBack(v.Name)
			classes.Add(v.Name)
		}
	case *ast.Not:
		opt(v.E)
	case *ast.Num:
	case *ast.Sub:
		opt(v.Left)
		opt(v.Right)
	case *ast.This:
	case *ast.Times:
		opt(v.Left)
		opt(v.Right)
	case *ast.True:
	default:
		panic("impossible")
	}
}

func opt_Stm(stm ast.Stm) {
	switch s := stm.(type) {
	case *ast.Assign:
		opt(s.E)
	case *ast.AssignArray:
		opt(s.Index)
		opt(s.E)
	case *ast.Block:
		for _, ss := range s.Stms {
			opt(ss)
		}
	case *ast.If:
		opt(s.Condition)
		opt(s.Thenn)
		opt(s.Elsee)
	case *ast.Print:
		opt(s.E)
	case *ast.While:
		opt(s.E)
		opt(s.Body)
	default:
		panic("impossible")

	}
}

func opt_Method(method ast.Method) {
	switch m := method.(type) {
	case *ast.MethodSingle:
		//omit the formals and locals
		//Statements
		for _, s := range m.Stms {
			opt(s)
		}
		//return exp
		opt(m.RetExp)
	default:
		panic("impossible")
	}
}

func opt_Class(class ast.Class) {
	switch c := class.(type) {
	case *ast.ClassSingle:
		//Add super class into worklist and set
		if c.Extends != "" {
			exist := false
			for e := worklist.Front(); e != nil; e = e.Next() {
				classid := e.Value
				if classid == c.Extends {
					exist = true
					break
				}
			}
			if exist == false {
				worklist.PushBack(c.Extends)
				classes.Add(c.Extends)
			}
		}
		//methods
		for _, m := range c.Methods {
			opt(m)
		}
	default:
		panic("impossible")
	}
}

func opt_MainClass(c ast.MainClass) {
	switch cc := c.(type) {
	case *ast.MainClassSingle:
		classes.Add(cc.Name)
		opt(cc.Stms)
	default:
		panic("impossible")
	}
}

func opt(e ast.Acceptable) {
	switch v := e.(type) {
	case ast.Exp:
		opt_Exp(v)
	case ast.Stm:
		opt_Stm(v)
	case ast.Type:
		//no need
	case ast.Dec:
		//no need
	case ast.Method:
		opt_Method(v)
	case ast.MainClass:
		opt_MainClass(v)
	case ast.Class:
		opt_Class(v)
	default:
		panic("impossible")
	}
}

func DeadClass_Opt(prog ast.Program) ast.Program {
	var p *ast.ProgramSingle
	if v, ok := prog.(*ast.ProgramSingle); ok {
		p = v
	} else {
		panic("impossible")
	}
	//add mainclass to the classes set
	opt(p.Mainclass)

	for worklist.Len() != 0 {
		e := worklist.Front()
		worklist.Remove(e)
		classid := e.Value
		for _, c := range p.Classes {
			if cc, ok := c.(*ast.ClassSingle); ok {
				if cc.Name == classid {
					opt(c)
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
			if classes.Contains(cc.Name) {
				newclasses = append(newclasses, c)
			}
		} else {
			panic("impossible")
		}
	}

	Ast := &ast.ProgramSingle{p.Mainclass, newclasses}

	//trace
	if util.Trace_contains("deadclass") == true {
		fmt.Println("before deadclass opt:")
		ast.NewPP().DumpProg(prog)
		fmt.Println("\nafter deadclass opt:")
		ast.NewPP().DumpProg(Ast)
	}
	return &ast.ProgramSingle{p.Mainclass, newclasses}
}
