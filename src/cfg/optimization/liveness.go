package cfg_opt

import (
	. "../../cfg"
	"../../control"
	"fmt"
)

func Liveness(p Program) {
	const (
		StmGenKill = iota
		BlockGenKill
		BlockInOut
		StmInOut
	)
	var Liveness_Kind int
	var do func(Acceptable)

	var current_method string
	//stm
	var oneStmGen map[string]bool
	var oneStmKill map[string]bool
	var oneTransferGen map[string]bool
	var oneTransferKill map[string]bool
	var stmGen map[Stm]map[string]bool
	var stmKill map[Stm]map[string]bool
	var transGen map[Transfer]map[string]bool
	var transKill map[Transfer]map[string]bool
	//block
	var blockGen map[Block]map[string]bool
	var blockKill map[Block]map[string]bool

	stmGen = make(map[Stm]map[string]bool)
	stmKill = make(map[Stm]map[string]bool)
	transGen = make(map[Transfer]map[string]bool)
	transKill = make(map[Transfer]map[string]bool)
	blockGen = make(map[Block]map[string]bool)
	blockKill = make(map[Block]map[string]bool)

	travers_GenKill := func(m map[string]bool) {
		for x, _ := range m {
			fmt.Printf("%s ", x)
		}
		fmt.Println("")
	}

	refresh_oneStmGenKill := func() {
		oneStmGen = make(map[string]bool)
		oneStmKill = make(map[string]bool)
	}
	refresh_oneTransGenKill := func() {
		oneTransferGen = make(map[string]bool)
		oneTransferKill = make(map[string]bool)
	}

	do_Operand := func(oo Operand) {
		switch o := oo.(type) {
		case *Int: //no need
		case *Var:
			oneStmGen[o.Name] = true
		default:
			panic("impossible")
		}
	}

	do_Transfer := func(tt Transfer) {
		switch t := tt.(type) {
		case *If:
			oneTransferKill[t.Cond.String()] = true
		case *Goto:
			//no need
		case *Return:
			oneTransferGen[t.Op.String()] = true
		default:
			panic("impossible")
		}
	}

	do_Stm := func(ss Stm) {
		switch s := ss.(type) {
		case *Add:
			oneStmKill[s.Dst] = true
			do(s.Left)
			do(s.Right)
		case *And:
			oneStmKill[s.Dst] = true
			do(s.Left)
			do(s.Right)
		case *ArraySelect:
			oneStmKill[s.Name] = true
			do(s.Arrayref)
			do(s.Index)
		case *AssignArray:
			oneStmKill[s.Dst] = true
			do(s.E)
			do(s.Index)
		case *InvokeVirtual:
			oneStmKill[s.Dst] = true
			oneStmGen[s.Obj] = true
			for _, arg := range s.Args {
				do(arg)
			}
		case *Length:
			oneStmKill[s.Dst] = true
			do(s.Arrayref)
		case *Lt:
			oneStmKill[s.Dst] = true
			do(s.Left)
			do(s.Right)
		case *Move:
			oneStmKill[s.Dst] = true
			do(s.Src)
		case *NewIntArray:
			oneStmKill[s.Dst] = true
			do(s.E)
		case *NewObject:
			oneStmKill[s.Dst] = true
		case *Not:
			oneStmKill[s.Dst] = true
			do(s.E)
		case *Print:
			do(s.Args)
		case *Sub:
			oneStmKill[s.Dst] = true
			do(s.Left)
			do(s.Right)
		case *Times:
			oneStmKill[s.Dst] = true
			do(s.Left)
			do(s.Right)
		default:
			panic("impossible")
		}
	}

	do_calculateStmGenKill := func(b *BlockSingle) {
		if control.Trace_contains("liveness.step1") {
			fmt.Println(current_method + " " + b.Label_id.String())
		}
		for _, s := range b.Stms {
			refresh_oneStmGenKill()
			do(s)
			stmGen[s] = oneStmGen
			stmKill[s] = oneStmKill
			if control.Trace_contains("liveness.step1") {
				fmt.Printf("  ")
				fmt.Println(s)
				fmt.Printf("    gen: ")
				travers_GenKill(oneStmGen)
				fmt.Printf("    kill: ")
				travers_GenKill(oneStmKill)
			}
		}
		refresh_oneTransGenKill()
		do(b.Trans)
		transGen[b.Trans] = oneTransferGen
		transKill[b.Trans] = oneTransferKill
		if control.Trace_contains("liveness.step1") {
			fmt.Printf("  ")
			fmt.Println(b.Trans)
			fmt.Printf("    gen: ")
			travers_GenKill(oneTransferGen)
			fmt.Printf("    kill: ")
			travers_GenKill(oneTransferKill)
			fmt.Println("")
		}
	}

	do_calculateBlockGenKill := func(b *BlockSingle) {
		if control.Trace_contains("liveness.step2") {
			fmt.Println(current_method + " " + b.Label_id.String())
		}
		//init set
		oneBlockGen := make(map[string]bool)
		oneBlockKill := make(map[string]bool)
		//addll
		for s, _ := range transGen[b.Trans] {
			oneBlockGen[s] = true
		}
		for s, _ := range transKill[b.Trans] {
			oneBlockKill[s] = true
		}
		//revers
		for i := len(b.Stms) - 1; i >= 0; i-- {
			stm := b.Stms[i]
			//remove all
			for s, _ := range stmKill[stm] {
				delete(oneBlockGen, s)
			}
			//addAll
			for s, _ := range stmGen[stm] {
				oneBlockGen[s] = true
			}
			//addAll
			for s, _ := range stmKill[stm] {
				oneBlockKill[s] = true
			}
		}
		blockGen[b] = oneBlockGen
		blockKill[b] = oneBlockKill
		if control.Trace_contains("liveness.step2") {
			fmt.Printf("  gen: ")
			travers_GenKill(oneBlockGen)
			fmt.Printf("  kill: ")
			travers_GenKill(oneBlockKill)
			fmt.Println("")
		}
	}

	do_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			switch Liveness_Kind {
			case StmGenKill:
				do_calculateStmGenKill(b)
			case BlockGenKill:
				do_calculateBlockGenKill(b)
			case BlockInOut:
			case StmInOut:
			default:
				panic("impossible")
			}
		default:
			panic("impossible")
		}
	}

	do_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			current_method = m.ClassId + "_" + m.Name
			Liveness_Kind = StmGenKill
			for _, b := range m.Blocks {
				do(b)
			}
			Liveness_Kind = BlockGenKill
			for _, b := range m.Blocks {
				do(b)
			}
		default:
			panic("impossible")
		}
	}

	do_MainMethod := func(mm MainMethod) {
	}

	do_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			do(p.Main_method)
			for _, m := range p.Methods {
				do(m)
			}
		default:
			panic("impossible")
		}
	}

	do = func(e Acceptable) {
		switch v := e.(type) {
		case Program:
			do_Program(v)
		case Block:
			do_Block(v)
		case Class:
		case Dec:
		case MainMethod:
			do_MainMethod(v)
		case Method:
			do_Method(v)
		case Operand:
			do_Operand(v)
		case Stm:
			do_Stm(v)
		case Transfer:
			do_Transfer(v)
		case Type:
		case Vtable:
		default:
			panic("impossible")
		}
	}

	do(p)

}
