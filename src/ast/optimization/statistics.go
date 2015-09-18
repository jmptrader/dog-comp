package ast_opt

import (
	"../../ast"
)

func Statistics_Ast(prog ast.Program) int {

	var do func(e ast.Acceptable)
	var size int

	do_Exp := func(ee ast.Exp) {
		switch e := ee.(type) {
		case *ast.Add:
			do(e.Left)
			do(e.Right)
			size++
		case *ast.And:
			do(e.Left)
			do(e.Right)
			size++
		case *ast.ArraySelect:
			do(e.Arrayref)
			do(e.Index)
			size++
		case *ast.Call:
			do(e.Callee)
			size++
		case *ast.False:
			size++
		case *ast.Id:
			size++
		case *ast.Length:
			do(e.Arrayref)
			size++
		case *ast.Lt:
			do(e.Left)
			do(e.Right)
			size++
		case *ast.NewIntArray:
			do(e.Size)
			size++
		case *ast.NewObject:
			size++
		case *ast.Not:
			do(e.E)
			size++
		case *ast.Num:
			size++
		case *ast.Sub:
			do(e.Left)
			do(e.Right)
			size++
		case *ast.This:
			size++
		case *ast.Times:
			do(e.Left)
			do(e.Right)
			size++
		case *ast.True:
			size++
		default:
			panic("impossible")
		}
	}

	do_Stm := func(ss ast.Stm) {
		switch s := ss.(type) {
		case *ast.Assign:
			do(s.E)
			size++
		case *ast.AssignArray:
			do(s.E)
			do(s.Index)
			size++
		case *ast.Block:
			for _, stm := range s.Stms {
				do(stm)
			}
			size++
		case *ast.If:
			do(s.Condition)
			do(s.Thenn)
			do(s.Elsee)
			size++
		case *ast.Print:
			do(s.E)
			size++
		case *ast.While:
			do(s.E)
			do(s.Body)
			size++
		default:
			panic("impossible")
		}
	}

	do_Dec := func(dd ast.Dec) {
		switch d := dd.(type) {
		case *ast.DecSingle:
			do(d.Tp)
			size++
		default:
			panic("impossible")
		}
	}

	do_Type := func(tt ast.Type) {
		size++
	}

	do_Method := func(mm ast.Method) {
		switch m := mm.(type) {
		case *ast.MethodSingle:
			do(m.RetType)
			for _, dec := range m.Formals {
				do(dec)
			}
			for _, dec := range m.Locals {
				do(dec)
			}
			for _, stm := range m.Stms {
				do(stm)
			}
			do(m.RetExp)
			size++
		default:
			panic("impossible")
		}
	}

	do_Class := func(cc ast.Class) {
		switch c := cc.(type) {
		case *ast.ClassSingle:
			for _, dec := range c.Decs {
				do(dec)
			}
			for _, m := range c.Methods {
				do(m)
			}
			size++
		default:
			panic("impossible")
		}
	}

	do_MainClass := func(cc ast.MainClass) {
		switch c := cc.(type) {
		case *ast.MainClassSingle:
			do(c.Stms)
			size++
		default:
			panic("impossible")
		}
	}

	do_Program := func(pr ast.Program) {
		switch p := pr.(type) {
		case *ast.ProgramSingle:
			do(p.Mainclass)
			for _, c := range p.Classes {
				do(c)
			}
			size++
		default:
			panic("impossible")
		}
	}

	do = func(e ast.Acceptable) {
		switch v := e.(type) {
		case ast.Class:
			do_Class(v)
		case ast.Dec:
			do_Dec(v)
		case ast.Exp:
			do_Exp(v)
		case ast.MainClass:
			do_MainClass(v)
		case ast.Method:
			do_Method(v)
		case ast.Program:
			do_Program(v)
		case ast.Stm:
			do_Stm(v)
		case ast.Type:
			do_Type(v)
		default:
			panic("impossible")
		}
	}

	do(prog)
	return size
}
