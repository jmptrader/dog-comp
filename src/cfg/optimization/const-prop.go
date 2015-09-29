package cfg_opt

import (
	. "../../cfg"
)

func ConstProp(prog Program) Program {

	var opt func(Acceptable)
	var opt_Operand func(Operand)
	var f_stm Stm
	var f_operand Operand
	var f_reaching_def map[Stm]bool

	opt_Operand = func(oo Operand) {
		switch o := oo.(type) {
		case *Int:
			f_operand = o
		case *Var:
			f_operand = o
			f_reaching_def = stmDefIn[f_stm]
			founded := false
			var move *Move = nil
			for ss, _ := range f_reaching_def {
				switch s := ss.(type) {
				case *Add:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *And:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *ArraySelect:
					if o.Name == s.Name {
						if founded == true {
							return
						}
						founded = true
					}
				case *AssignArray:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *InvokeVirtual:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *Length:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *Lt:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *Move:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
						move = s
					}
				case *NewIntArray:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *NewObject:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *Not:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *Print:
				case *Sub:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				case *Times:
					if o.Name == s.Dst {
						if founded == true {
							return
						}
						founded = true
					}
				default:
					panic("impossible")
				}
			}
			if founded {
				if move != nil {
					if op_int, ok := move.Src.(*Int); ok {
						f_operand = &Int{op_int.Value}
					}
				}
			}
		default:
			panic("impossible")
		}
	}

	opt_Stm := func(ss Stm) {
		f_stm = ss
		switch s := ss.(type) {
		case *Add:
			opt_Operand(s.Left)
			s.Left = f_operand
			opt_Operand(s.Right)
			s.Right = f_operand
		case *And:
			opt_Operand(s.Left)
			s.Left = f_operand
			opt_Operand(s.Right)
			s.Right = f_operand
		case *ArraySelect:
			opt_Operand(s.Index)
			s.Index = f_operand
		case *AssignArray:
			opt_Operand(s.Index)
			s.Index = f_operand
		case *InvokeVirtual:
		case *Length:
		case *Lt:
			opt_Operand(s.Left)
			s.Left = f_operand
			opt_Operand(s.Right)
			s.Right = f_operand
		case *Move:
			opt_Operand(s.Src)
			s.Src = f_operand
		case *NewIntArray:
		case *NewObject:
		case *Not:
			opt_Operand(s.E)
			s.E = f_operand
		case *Print:
		case *Sub:
			opt_Operand(s.Left)
			s.Left = f_operand
			opt_Operand(s.Right)
			s.Right = f_operand
		case *Times:
			opt_Operand(s.Left)
			s.Left = f_operand
			opt_Operand(s.Right)
			s.Right = f_operand
		default:
			panic("impossible")
		}
	}

	opt_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			for _, s := range b.Stms {
				opt(s)
			}
		default:
			panic("impossible")
		}
	}

	opt_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			for _, b := range m.Blocks {
				opt(b)
			}
		default:
			panic("impossible")
		}
	}

	opt_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			for _, m := range p.Methods {
				opt(m)
			}
		default:
			panic("impossible")
		}
	}

	opt = func(e Acceptable) {
		switch v := e.(type) {
		case Block:
			opt_Block(v)
		case Class:
		case Dec:
		case MainMethod:
		case Method:
			opt_Method(v)
		case Operand:
			opt_Operand(v)
		case Program:
			opt_Program(v)
		case Stm:
			opt_Stm(v)
		case Transfer:
		case Type:
		case Vtable:
		default:
			panic("impossible")
		}
	}

	opt(prog)
	return prog
}
