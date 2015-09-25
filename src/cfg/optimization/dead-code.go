package cfg_opt

import (
	. "../../cfg"
)

func DeadCode(prog Program) Program {

	var f_block Block
	var f_method Method
	var f_stm Stm
	var opt func(Acceptable)

	opt_Stm := func(ss Stm) {
		out := stmLiveOut[ss]
		elim := func(dst string) {
			if out[dst] != true {
				f_stm = nil
			} else {
				f_stm = ss
			}
		}

		switch s := ss.(type) {
		case *Add:
			elim(s.Dst)
		case *And:
			elim(s.Dst)
		case *ArraySelect:
			elim(s.Name)
		case *AssignArray:
			elim(s.Dst)
			if s.IsField == true {
				f_stm = s
			}
		case *InvokeVirtual:
			//XXX invoke method may generate side-effect
			f_stm = s
		case *Length:
			elim(s.Dst)
		case *Lt:
			elim(s.Dst)
		case *Move:
			elim(s.Dst)
			if s.IsField == true {
				f_stm = s
				return
			}
			if t := s.Tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
				f_stm = s
			}
		case *NewIntArray:
			elim(s.Dst)
		case *NewObject:
			elim(s.Dst)
		case *Not:
			elim(s.Dst)
		case *Print:
			f_stm = s
		case *Sub:
			elim(s.Dst)
		case *Times:
			elim(s.Dst)
		default:
			panic("impossible")
		}
	}

	opt_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			stms := make([]Stm, 0)
			for _, s := range b.Stms {
				opt(s)
				if f_stm != nil {
					stms = append(stms, f_stm)
				}
			}
			//XXX
			//b.Stms = stms
			f_block = &BlockSingle{b.Label_id, stms, b.Trans}

		default:
			panic("impossible")
		}
	}

	opt_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			blocks := make([]Block, 0)
			for _, b := range m.Blocks {
				opt(b)
				blocks = append(blocks, f_block)
			}
			f_method = &MethodSingle{m.Ret_type, m.Name, m.ClassId,
				m.Formals, m.Locals, blocks, m.Entry}
		default:
			panic("impossible")
		}
	}

	opt_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			methods := make([]Method, 0)
			for _, m := range p.Methods {
				opt(m)
				methods = append(methods, f_method)
			}
			p.Methods = methods
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
