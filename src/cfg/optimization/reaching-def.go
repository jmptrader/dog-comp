package cfg_opt

import (
	. "../../cfg"
	"../../control"
	"../../util"
	"fmt"
)

var stmDefIn map[Stm]map[Stm]bool
var stmDefOut map[Stm]map[Stm]bool
var defs map[string]map[Stm]bool

func ReachingDef(prog Program) Program {

	var f_prev map[*util.Node]bool
	var changed bool
	var current_method *MethodSingle
	//gen[s] = {d}
	//kill[s] = defs(t)-{d}
	var stmGen map[Stm]map[Stm]bool
	var stmKill map[Stm]map[Stm]bool
	var blockGen map[Block]map[Stm]bool
	var blockKill map[Block]map[Stm]bool
	var blockIn map[Block]map[Stm]bool
	var blockOut map[Block]map[Stm]bool
	var oneStmGen map[Stm]bool
	var oneStmKill map[Stm]bool

	stmGen = make(map[Stm]map[Stm]bool)
	stmKill = make(map[Stm]map[Stm]bool)
	blockGen = make(map[Block]map[Stm]bool)
	blockKill = make(map[Block]map[Stm]bool)
	blockIn = make(map[Block]map[Stm]bool)
	blockOut = make(map[Block]map[Stm]bool)
	stmDefIn = make(map[Stm]map[Stm]bool)
	stmDefOut = make(map[Stm]map[Stm]bool)
	defs = make(map[string]map[Stm]bool)

	var do func(Acceptable)
	const (
		StmGenKill = iota
		BlockGenKill
		BlockInOut
		StmInOut
		InitDefs
	)
	var ReachingDef_Kind int

	do_Transfer := func(tt Transfer) {
	}
	refresh_oneStmGenKill := func() {
		oneStmGen = make(map[Stm]bool)
		oneStmKill = make(map[Stm]bool)
	}

	do_Stm := func(ss Stm) {
		switch s := ss.(type) {
		case *Add:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			//addAll defs(t)
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			//removeAll
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *And:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *AssignArray:
		case *ArraySelect:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Name] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *InvokeVirtual:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *Length:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *Lt:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *Move:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *NewIntArray:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *NewObject:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *Not:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *Print:
		case *Sub:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		case *Times:
			oneStmGen[s] = true
			stmGen[s] = oneStmGen //get gen
			for x, _ := range defs[s.Dst] {
				oneStmKill[x] = true
			}
			for x, _ := range oneStmGen {
				delete(oneStmKill, x)
			}
			stmKill[s] = oneStmKill
		default:
			panic("impossible")
		}
	}

	do_calculateStmGenKill := func(b *BlockSingle) {
		if control.Trace_contains("reaching.step1") {
			fmt.Println(current_method.Name + " " + b.Label_id.String())
		}

		for _, s := range b.Stms {
			refresh_oneStmGenKill()
			do(s)
			if control.Trace_contains("reaching.step1") {
				fmt.Printf("  ")
				fmt.Println(s)
				fmt.Println("  gen:")
				for x, _ := range oneStmGen {
					fmt.Printf("    ")
					fmt.Println(x)
				}
				fmt.Println("")
				fmt.Println("  kill:")
				for x, _ := range oneStmKill {
					fmt.Printf("    ")
					fmt.Println(x)
				}
				fmt.Println("")
			}

		}

	}

	do_calculateBlockGenKill := func(b *BlockSingle) {
		oneBlockGen := make(map[Stm]bool)
		oneBlockKill := make(map[Stm]bool)
		for _, s := range b.Stms {
			oneStmGen := stmGen[s]
			oneStmKill := stmKill[s]
			for x, _ := range oneStmGen {
				oneBlockGen[x] = true
			}
			for x, _ := range oneStmKill {
				oneBlockKill[x] = true
			}
		}
		blockGen[b] = oneBlockGen
		blockKill[b] = oneBlockKill
	}

	do_calculateBlockInOut := func(b *BlockSingle) bool {
		oneBlockGen := blockGen[b]
		oneBlockKill := blockKill[b]
		oneBlockIn := make(map[Stm]bool)
		oneBlockOut := make(map[Stm]bool)

		in_size := 0
		out_size := 0

		if x := blockIn[b]; x != nil {
			in_size = len(x)
		}
		if x := blockOut[b]; x != nil {
			out_size = len(x)
		}

		//in[n] = U out[p]      p is pred[n]
		for n, _ := range f_prev {
			bb := n.GetData()
			if x, ok := bb.(Block); ok {
				if blockOut[x] != nil {
					for s, _ := range blockOut[x] {
						oneBlockIn[s] = true
					}
				}
			} else {
				panic("impossible")
			}
		}
		blockIn[b] = oneBlockIn

		//out[n] = gen[n]U(in[n]-kill[n])
		for s, _ := range oneBlockIn {
			oneBlockOut[s] = true
		}
		for s, _ := range oneBlockKill {
			delete(oneBlockOut, s)
		}
		for s, _ := range oneBlockGen {
			oneBlockOut[s] = true
		}
		blockOut[b] = oneBlockOut

		if len(oneBlockIn) != in_size || len(oneBlockOut) != out_size {
			return true
		} else {
			return false
		}
	}

	do_calculateStmInOut := func(b *BlockSingle) {
		if control.Trace_contains("reaching.step4") {
			fmt.Println(current_method.Name + " " + b.Label_id.String())
		}
		prev_out := make(map[Stm]bool)
		for s, _ := range blockIn[b] {
			prev_out[s] = true
		}

		for _, s := range b.Stms {
			oneStmIn := make(map[Stm]bool)
			oneStmOut := make(map[Stm]bool)
			oneStmGen := stmGen[s]
			oneStmKill := stmKill[s]

			//in[n] = U out[p]      p is pred[n]
			for x, _ := range prev_out {
				oneStmIn[x] = true
			}
			stmDefIn[s] = oneStmIn

			//out[n] = gen[n]U(in[n]-kill[n])
			for x, _ := range oneStmIn {
				oneStmOut[x] = true
			}
			for x, _ := range oneStmKill {
				delete(oneStmOut, x)
			}
			for x, _ := range oneStmGen {
				oneStmOut[x] = true
			}
			stmDefOut[s] = oneStmOut
			prev_out = oneStmOut
			if control.Trace_contains("reaching.step4") {
				fmt.Printf("  ")
				fmt.Println(s)
				fmt.Println("  In:")
				for x, _ := range oneStmIn {
					fmt.Printf("    ")
					fmt.Println(x)
				}
				fmt.Println("  Out:")
				for x, _ := range oneStmOut {
					fmt.Printf("    ")
					fmt.Println(x)
				}
				fmt.Println("")
			}

		}
	}

	do_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			switch ReachingDef_Kind {
			case StmGenKill:
				do_calculateStmGenKill(b)
			case BlockGenKill:
				do_calculateBlockGenKill(b)
			case BlockInOut:
				changed = changed || do_calculateBlockInOut(b)
			case StmInOut:
				do_calculateStmInOut(b)
			default:
				panic("impossible")

			}
		default:
			panic("impossible")
		}
	}

	initDefs := func(m *MethodSingle) {
		for _, bb := range m.Blocks {
			switch b := bb.(type) {
			case *BlockSingle:
				for _, ss := range b.Stms {
					switch s := ss.(type) {
					case *Add:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *And:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *ArraySelect:
						if defs[s.Name] == nil {
							defs[s.Name] = make(map[Stm]bool)
						}
						set := defs[s.Name]
						set[s] = true
					case *AssignArray:
					case *InvokeVirtual:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *Length:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *Lt:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *Move:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *NewIntArray:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *NewObject:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *Not:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *Print:
					case *Sub:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					case *Times:
						if defs[s.Dst] == nil {
							defs[s.Dst] = make(map[Stm]bool)
						}
						set := defs[s.Dst]
						set[s] = true
					default:
						panic("impossible")
					}
				}
			default:
				panic("impossible")
			}
		}
	}

	do_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			current_method = m
			initDefs(m)
			//step1  stm gen kill
			ReachingDef_Kind = StmGenKill
			for _, b := range m.Blocks {
				do(b)
			}
			//step2 block gen kill
			ReachingDef_Kind = BlockGenKill
			for _, b := range m.Blocks {
				do(b)
			}
			ReachingDef_Kind = BlockInOut
			graph := GenGraph(m)
			nodes := graph.GetNodes()
			//fix-point
			times := 0
			changed = true
			for changed {
				times++
				changed = false
				for _, node := range nodes {
					f_prev = node.GetPrev()
					if b, ok := node.GetData().(Block); ok {
						do(b)
					} else {
						panic("impossible")
					}
				}
			}
			if control.Trace_contains("reaching.step3") {
				fmt.Printf("%d times reach pix-piont\n", times)
				for _, bb := range m.Blocks {
					if b, ok := bb.(*BlockSingle); ok {
						fmt.Println(current_method.Name + " " + b.Label_id.String())
						fmt.Println("  In:")
						for x, _ := range blockIn[b] {
							fmt.Printf("    ")
							fmt.Println(x)
						}
						fmt.Println("  Out:")
						for x, _ := range blockIn[b] {
							fmt.Printf("    ")
							fmt.Println(x)
						}
						fmt.Println("")
					}
				}
			}

			//step 4
			ReachingDef_Kind = StmInOut
			for _, b := range m.Blocks {
				do(b)
			}
		default:
			panic("impossible")
		}
	}

	do_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			for _, m := range p.Methods {
				do(m)
			}
		default:
			panic("impossible")
		}
	}

	do = func(e Acceptable) {
		switch v := e.(type) {
		case Block:
			do_Block(v)
		case Class:
		case Dec:
		case MainMethod:
		case Method:
			do_Method(v)
		case Operand:
		case Program:
			do_Program(v)
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

	do(prog)

	return prog
}
