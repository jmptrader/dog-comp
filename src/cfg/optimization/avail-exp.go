package cfg_opt

import (
	. "../../cfg"
	"../../control"
	"../../util"
	"fmt"
)

var stmExpIn map[Stm]map[Stm]bool

func AvailExp(prog Program) Program {
	var changed bool
	var f_prev map[*util.Node]bool
	var all map[Stm]bool
	var exps map[string]map[Stm]bool
	var avail func(Acceptable)
	var stmGen map[Stm]map[Stm]bool
	var stmKill map[Stm]map[Stm]bool
	var blockGen map[Block]map[Stm]bool
	var blockKill map[Block]map[Stm]bool
	var oneStmGen map[Stm]bool
	var oneStmKill map[Stm]bool
	var blockIn map[Block]map[Stm]bool
	var blockOut map[Block]map[Stm]bool
	var stmExpOut map[Stm]map[Stm]bool
	var current_method *MethodSingle

	all = make(map[Stm]bool)
	exps = make(map[string]map[Stm]bool)
	stmGen = make(map[Stm]map[Stm]bool)
	stmKill = make(map[Stm]map[Stm]bool)
	blockGen = make(map[Block]map[Stm]bool)
	blockKill = make(map[Block]map[Stm]bool)
	blockIn = make(map[Block]map[Stm]bool)
	blockOut = make(map[Block]map[Stm]bool)
	stmExpIn = make(map[Stm]map[Stm]bool)
	stmExpOut = make(map[Stm]map[Stm]bool)
	const (
		StmGenKill = iota
		BlockGenKill
		BlockInOut
		StmInOut
		InitExp
	)
	var AvailExp_Kind int

	avail_Stm := func(ss Stm) {
		switch s := ss.(type) {
		case *Add:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		case *And:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		case *ArraySelect:
			kill := exps[s.Name]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		case *AssignArray:
			//kill all form of M[x]
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
		case *InvokeVirtual:
			//kill all
			for _, set := range exps {
				for x, _ := range set {
					oneStmKill[x] = true
				}
			}
		case *Length:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
		case *Lt:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		case *Move:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		case *NewIntArray:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
		case *NewObject:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
		case *Not:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		case *Print:
		case *Sub:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		case *Times:
			kill := exps[s.Dst]
			for x, _ := range kill {
				oneStmKill[x] = true
			}
			oneStmGen[s] = true
			for x, _ := range oneStmKill {
				delete(oneStmGen, x)
			}
		default:
			panic("impossible")
		}
	}

	avail_calculateStmGenKill := func(b *BlockSingle) {
		if control.Trace_contains("avail.step1") {
			fmt.Println(current_method.Name + " " + b.Label_id.String())
		}
		for _, s := range b.Stms {
			oneStmGen = make(map[Stm]bool)
			oneStmKill = make(map[Stm]bool)
			avail(s)
			stmGen[s] = oneStmGen
			stmKill[s] = oneStmKill

			if control.Trace_contains("avail.step1") {
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

	avail_calculateBlockGenKill := func(b *BlockSingle) {
		oneBlockGen := make(map[Stm]bool)
		oneBlockKill := make(map[Stm]bool)
		for _, s := range b.Stms {
			for x, _ := range stmGen[s] {
				oneBlockGen[x] = true
			}
			for x, _ := range stmKill[s] {
				oneBlockKill[x] = true
			}
		}
		blockGen[b] = oneBlockGen
		blockKill[b] = oneBlockKill

		if control.Trace_contains("avail.step2") {
			fmt.Println(current_method.Name + " " + b.Label_id.String())
			fmt.Println("  gen:")
			for x, _ := range oneBlockGen {
				fmt.Printf("    ")
				fmt.Println(x)
			}
			fmt.Println("")
			fmt.Println("  kill:")
			for x, _ := range oneBlockKill {
				fmt.Printf("    ")
				fmt.Println(x)
			}
			fmt.Println("")
		}
	}

	avail_calculateBlockInOut := func(b *BlockSingle) bool {
		oneBlockGen := blockGen[b]
		oneBlockKill := blockKill[b]
		oneBlockIn := blockIn[b]
		oneBlockOut := blockOut[b]

		in_size := len(oneBlockIn)
		out_size := len(oneBlockOut)

		//in[n] = ^out[p] p belong to prev[n]
		for n, _ := range f_prev {
			if p, ok := n.GetData().(Block); ok {
				out := blockOut[p]
				for s, _ := range oneBlockIn {
					if out[s] != true {
						delete(oneBlockIn, s)
					}
				}
			} else {
				panic("impossible")
			}
		}
		blockIn[b] = oneBlockIn
		//out[n] = gen[n]U(in[n]-kill[n])
		for x, _ := range oneBlockIn {
			oneBlockOut[x] = true
		}
		for x, _ := range oneBlockKill {
			delete(oneBlockOut, x)
		}
		for x, _ := range oneBlockGen {
			oneBlockOut[x] = true
		}
		blockOut[b] = oneBlockOut

		if len(oneBlockIn) != in_size || len(oneBlockOut) != out_size {
			return true
		} else {
			return false
		}
	}

	avail_calculateStmInOut := func(b *BlockSingle) {
		if control.Trace_contains("avail.step4") {
			fmt.Println(current_method.Name + " " + b.Label_id.String())
		}
		prev_out := make(map[Stm]bool)
		for x, _ := range blockIn[b] {
			prev_out[x] = true
		}
		for _, s := range b.Stms {
			oneStmGen := stmGen[s]
			oneStmKill := stmKill[s]
			oneStmIn := make(map[Stm]bool)
			oneStmOut := make(map[Stm]bool)
			for x, _ := range prev_out {
				oneStmIn[x] = true
			}
			stmExpIn[s] = oneStmIn
			for x, _ := range oneStmIn {
				oneStmOut[x] = true
			}
			for x, _ := range oneStmKill {
				delete(oneStmOut, x)
			}
			for x, _ := range oneStmGen {
				oneStmOut[x] = true
			}
			stmExpOut[s] = oneStmOut
			prev_out = oneStmOut

			if control.Trace_contains("avail.step4") {
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

	avail_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			switch AvailExp_Kind {
			case StmGenKill:
				avail_calculateStmGenKill(b)
			case BlockGenKill:
				avail_calculateBlockGenKill(b)
			case BlockInOut:
				changed = changed || avail_calculateBlockInOut(b)
			case StmInOut:
				avail_calculateStmInOut(b)
			default:
				panic("impossible")
			}
		default:
			panic("impossible")
		}
	}

	initExp := func(m *MethodSingle) {
		for _, bb := range m.Blocks {
			switch b := bb.(type) {
			case *BlockSingle:
				for _, ss := range b.Stms {
					//init the all
					all[ss] = true
					init_Operand := func(oo Operand) {
						switch o := oo.(type) {
						case *Int:
						case *Var:
							if exps[o.Name] == nil {
								exps[o.Name] = make(map[Stm]bool)
							}
							set := exps[o.Name]
							set[ss] = true
						default:
							panic("impossible")
						}
					}
					switch s := ss.(type) {
					case *Add:
						init_Operand(s.Left)
						init_Operand(s.Right)
					case *And:
						init_Operand(s.Left)
						init_Operand(s.Right)
					case *ArraySelect:
						init_Operand(s.Arrayref)
						init_Operand(s.Index)
					case *AssignArray:
						init_Operand(s.Index)
						init_Operand(s.E)
					case *InvokeVirtual:
						for _, arg := range s.Args {
							init_Operand(arg)
						}
					case *Length:
						init_Operand(s.Arrayref)
					case *Lt:
						init_Operand(s.Left)
						init_Operand(s.Right)
					case *Move:
						init_Operand(s.Src)
					case *NewIntArray:
						init_Operand(s.E)
					case *NewObject:
					case *Not:
						init_Operand(s.E)
					case *Print:
						init_Operand(s.Args)
					case *Sub:
						init_Operand(s.Left)
						init_Operand(s.Right)
					case *Times:
						init_Operand(s.Left)
						init_Operand(s.Right)
					default:
						panic("impossible")
					}
				}
			default:
				panic("impossible")
			}
		}
	}

	avail_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			current_method = m
			initExp(m)
			//step1
			AvailExp_Kind = StmGenKill
			for _, b := range m.Blocks {
				avail(b)
			}
			//step2
			AvailExp_Kind = BlockGenKill
			for _, b := range m.Blocks {
				avail(b)
			}

			//step3
			AvailExp_Kind = BlockInOut
			graph := GenGraph(m)
			nodes := graph.GetNodes()
			//init set
			for idx, node := range nodes {
				if b, ok := node.GetData().(Block); ok {
					oneBlockOut := make(map[Stm]bool)
					oneBlockIn := make(map[Stm]bool)
					//init universal set
					for x, _ := range all {
						oneBlockOut[x] = true
						oneBlockIn[x] = true
					}
					//the first block's in init empty
					if idx == 0 {
						oneBlockIn = make(map[Stm]bool)
					}
					blockIn[b] = oneBlockIn
					blockOut[b] = oneBlockOut
				}
			}
			//fix-point
			times := 0
			changed = true
			for changed {
				times++
				changed = false
				for _, node := range nodes {
					f_prev = node.GetPrev()
					if b, ok := node.GetData().(Block); ok {
						avail(b)
					} else {
						panic("impossible")
					}
				}
			}
			if control.Trace_contains("avail.step3") {
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
						for x, _ := range blockOut[b] {
							fmt.Printf("    ")
							fmt.Println(x)
						}
						fmt.Println("")
					}
				}
			}

			//step4
			AvailExp_Kind = StmInOut
			for _, b := range m.Blocks {
				avail(b)
			}
		default:
			panic("impossible")
		}
	}

	avail_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			for _, m := range p.Methods {
				avail(m)
			}
		default:
			panic("impossible")
		}
	}

	avail = func(e Acceptable) {
		switch v := e.(type) {
		case Block:
			avail_Block(v)
		case Class:
		case Dec:
		case MainMethod:
		case Method:
			avail_Method(v)
		case Operand:
		case Program:
			avail_Program(v)
		case Stm:
			avail_Stm(v)
		case Transfer:
		case Type:
		case Vtable:
		default:
			panic("impossible")
		}
	}

	avail(prog)

	/*
	   for id, set := range exps{
	       fmt.Println(id)
	       for s, _ := range set{
	           fmt.Printf("  ")
	           fmt.Println(s)
	       }
	   }
	*/

	return prog
}
