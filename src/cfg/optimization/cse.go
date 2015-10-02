package cfg_opt

import (
	. "../../cfg"
)

func Cse(prog Program) Program {

	var cse func(Acceptable)
	var f_stm Stm

	emit := func(dst string, src string) {
		src_var := &Var{src, false}
		f_stm = &Move{dst, &IntType{}, src_var, false}
	}

	operand_equal := func(op1, op2 Operand) bool {
		switch op := op1.(type) {
		case *Int:
			if o, ok := op2.(*Int); ok {
				if op.Value == o.Value {
					return true
				}
			}
			return false
		case *Var:
			if o, ok := op2.(*Var); ok {
				if op.Name == o.Name {
					return true
				}
			}
			return false
		default:
			panic("impossible")
		}
	}

	cse_Stm := func(ss Stm) {
		avail_exp := stmExpIn[ss]
		f_stm = ss
		switch s := ss.(type) {
		case *Add:
			for x, _ := range avail_exp {
				if xx, ok := x.(*Add); ok {
					if operand_equal(s.Left, xx.Left) && operand_equal(s.Right, xx.Right) {
						emit(s.Dst, xx.Dst)
					}
				}
			}
		case *And:
			for x, _ := range avail_exp {
				if xx, ok := x.(*And); ok {
					if operand_equal(s.Left, xx.Left) && operand_equal(s.Right, xx.Right) {
						emit(s.Dst, xx.Dst)
					}
				}
			}
		case *ArraySelect:
		case *AssignArray:
		case *InvokeVirtual:
		case *Length:
		case *Lt:
			for x, _ := range avail_exp {
				if xx, ok := x.(*Lt); ok {
					if operand_equal(s.Left, xx.Left) && operand_equal(s.Right, xx.Right) {
						emit(s.Dst, xx.Dst)
					}
				}
			}
		case *Move:
		case *NewIntArray:
		case *NewObject:
		case *Not:
		case *Print:
		case *Sub:
			for x, _ := range avail_exp {
				if xx, ok := x.(*Sub); ok {
					if operand_equal(s.Left, xx.Left) && operand_equal(s.Right, xx.Right) {
						emit(s.Dst, xx.Dst)
					}
				}
			}
		case *Times:
			for x, _ := range avail_exp {
				if xx, ok := x.(*Times); ok {
					if operand_equal(s.Left, xx.Left) && operand_equal(s.Right, xx.Right) {
						emit(s.Dst, xx.Dst)
					}
				}
			}
		default:
			panic("impossible")
		}
	}

	cse_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			stms := make([]Stm, 0)
			for _, s := range b.Stms {
				cse(s)
				stms = append(stms, f_stm)
			}
			b.Stms = stms
		default:
			panic("impossible")
		}
	}

	cse_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			for _, b := range m.Blocks {
				cse(b)
			}
		default:
			panic("impossible")
		}
	}

	cse_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			for _, m := range p.Methods {
				cse(m)
			}
		default:
			panic("impossible")
		}
	}

	cse = func(e Acceptable) {
		switch v := e.(type) {
		case Block:
			cse_Block(v)
		case Class:
		case Dec:
		case MainMethod:
		case Method:
			cse_Method(v)
		case Operand:
		case Program:
			cse_Program(v)
		case Stm:
			cse_Stm(v)
		case Transfer:
		case Type:
		case Vtable:
		default:
			panic("impossible")
		}
	}

	cse(prog)
	return prog
}
