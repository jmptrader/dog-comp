package cfg_opt

import (
	. "../../cfg"
	"../../control"
	"../../util"
	"fmt"
)

var stmGen map[Stm]map[string]bool
var stmKill map[Stm]map[string]bool
var transGen map[Transfer]map[string]bool
var transKill map[Transfer]map[string]bool

//setp4
var stmLiveIn map[Stm]map[string]bool
var stmLiveOut map[Stm]map[string]bool
var transLiveIn map[Transfer]map[string]bool
var transLiveOut map[Transfer]map[string]bool

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
	//setp1
	var oneStmGen map[string]bool
	var oneStmKill map[string]bool
	var oneTransferGen map[string]bool
	var oneTransferKill map[string]bool
	//setp2
	var blockGen map[Block]map[string]bool
	var blockKill map[Block]map[string]bool
	//setp3
	var blockLiveIn map[Block]map[string]bool
	var blockLiveOut map[Block]map[string]bool

	var f_succ map[*util.Node]bool

	stmGen = make(map[Stm]map[string]bool)
	stmKill = make(map[Stm]map[string]bool)
	transGen = make(map[Transfer]map[string]bool)
	transKill = make(map[Transfer]map[string]bool)
	blockGen = make(map[Block]map[string]bool)
	blockKill = make(map[Block]map[string]bool)
	blockLiveIn = make(map[Block]map[string]bool)
	blockLiveOut = make(map[Block]map[string]bool)
	stmLiveIn = make(map[Stm]map[string]bool)
	stmLiveOut = make(map[Stm]map[string]bool)
	transLiveIn = make(map[Transfer]map[string]bool)
	transLiveOut = make(map[Transfer]map[string]bool)

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
			oneTransferGen[t.Cond.String()] = true
		case *Goto:
			//no need
		case *Return:
			if v, ok := t.Op.(*Var); ok {
				oneTransferGen[v.String()] = true
			}
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

	do_calculateBlockInOut := func(b *BlockSingle) bool {
		oneBlockGen := blockGen[b]
		oneBlockKill := blockKill[b]
		oneBlockIn := make(map[string]bool)
		oneBlockOut := make(map[string]bool)
		tempOut := make(map[string]bool)

		var in_size int
		var out_size int

		if b := blockLiveIn[b]; b != nil {
			in_size = len(b)
		} else {
			in_size = 0
		}
		if b := blockLiveOut[b]; b != nil {
			out_size = len(b)
		} else {
			out_size = 0
		}

		//out[n] = U in[s] s belong to the n.cucc
		for n, _ := range f_succ {
			bb := n.GetData()
			if b, ok := bb.(Block); ok {
				if blockLiveIn[b] != nil {
					//addall
					for s, _ := range blockLiveIn[b] {
						oneBlockOut[s] = true
					}
				}
			} else {
				panic("impossible")
			}
		}
		blockLiveOut[b] = oneBlockOut
		//in[n] = use[n] U (out[n]-def[n])

		//out[n]-def[n]
		for s, _ := range oneBlockOut {
			tempOut[s] = true
		}
		for s, _ := range oneBlockKill {
			delete(tempOut, s)
		}
		//use[n] U (out[n]-def[n])
		for s, _ := range oneBlockGen {
			oneBlockIn[s] = true
		}
		for s, _ := range tempOut {
			oneBlockIn[s] = true
		}
		blockLiveIn[b] = oneBlockIn

		if len(oneBlockOut) != out_size || len(oneBlockIn) != in_size {
			return true
		} else {
			//trace
			if control.Trace_contains("liveness.step3") {
				fmt.Println(current_method + " " + b.Label_id.String())
				fmt.Printf("  In: ")
				travers_GenKill(oneBlockIn)
				fmt.Printf("  Out: ")
				travers_GenKill(oneBlockOut)
				fmt.Println("")
			}
			return false
		}

	}

	do_calculateStmInOut := func(b *BlockSingle) {
		if control.Trace_contains("liveness.step4") {
			fmt.Println(current_method + " " + b.Label_id.String())
		}
		//transfer in out
		oneTransferGen = transGen[b.Trans]
		oneTransferKill = transKill[b.Trans]
		oneTransferOut := make(map[string]bool)
		oneTransferIn := make(map[string]bool)
		for s, _ := range blockLiveOut[b] {
			oneTransferOut[s] = true
		}
		transLiveOut[b.Trans] = oneTransferOut
		//in[n] = use[n] U (out[n] - def[n])
		for s, _ := range oneTransferOut {
			oneTransferIn[s] = true
		}
		for s, _ := range oneTransferKill {
			delete(oneTransferIn, s)
		}
		for s, _ := range oneTransferGen {
			oneTransferIn[s] = true
		}
		transLiveIn[b.Trans] = oneTransferIn
		if control.Trace_contains("liveness.step4") {
			fmt.Printf("  ")
			fmt.Println(b.Trans)
			fmt.Printf("    In: ")
			travers_GenKill(oneTransferIn)
			fmt.Printf("    Out: ")
			travers_GenKill(oneTransferOut)
			fmt.Println("")
		}

		//Stm in out
		prevIn := make(map[string]bool)
		//prevOut := make(map[string]bool)
		prevIn = oneTransferIn
		//prevOut = oneTransferOut
		for i := len(b.Stms) - 1; i >= 0; i-- {
			oneStmIn := make(map[string]bool)
			oneStmOut := make(map[string]bool)
			oneStmGen = stmGen[b.Stms[i]]
			oneStmKill = stmKill[b.Stms[i]]
			//out
			for s, _ := range prevIn {
				oneStmOut[s] = true
			}
			stmLiveOut[b.Stms[i]] = oneStmOut
			//in
			for s, _ := range oneStmOut {
				oneStmIn[s] = true
			}
			for s, _ := range oneStmKill {
				delete(oneStmIn, s)
			}
			for s, _ := range oneStmGen {
				oneStmIn[s] = true
			}
			stmLiveIn[b.Stms[i]] = oneStmIn
			prevIn = oneStmIn
			// prevOut = oneStmOut
			if control.Trace_contains("liveness.step4") {
				fmt.Printf("  ")
				fmt.Println(b.Stms[i])
				fmt.Printf("    In: ")
				travers_GenKill(oneStmIn)
				fmt.Printf("    Out: ")
				travers_GenKill(oneStmOut)
				fmt.Println("")
			}

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
				changed := true
				times := 1
				for changed {
					changed = do_calculateBlockInOut(b)
					times++
					util.Assert(times < 20, func() { panic("out of time") })
				}
			case StmInOut:
				do_calculateStmInOut(b)
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
			//step1
			Liveness_Kind = StmGenKill
			for _, b := range m.Blocks {
				do(b)
			}
			//step2
			Liveness_Kind = BlockGenKill
			for _, b := range m.Blocks {
				do(b)
			}
			//setp3
			Liveness_Kind = BlockInOut
			graph := GenGraph(m)
			retop_nodes := graph.Quasi_reverse()
			//check
			util.Assert(len(m.Blocks) == len(retop_nodes), func() { panic("assert fault") })
			// check the quasi-top
			/*
			   fmt.Println("\n"+ current_method)
			   for _, n := range retop_nodes{
			       b := n.GetData()
			       if v , ok := b.(*BlockSingle); ok{
			           fmt.Println(v.Label_id.String())
			       }else{
			           panic("impossible")
			       }
			   }
			*/
			for _, node := range retop_nodes {
				f_succ = node.GetSucc()
				if b, ok := node.GetData().(Block); ok {
					do(b)
				} else {
					panic("impossible")
				}
			}

			//step4
			Liveness_Kind = StmInOut
			for _, b := range m.Blocks {
				do_Block(b)
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
