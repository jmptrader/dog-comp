package cfg

import (
	"../control"
	"../util"
	"bytes"
	"fmt"
	"strconv"
)

func PP(e interface{}) string {

	var buf bytes.Buffer
	var pp func(interface{})

	emit := func(s string) {
		buf.Write([]byte(s))
	}

	pp_Transfer := func(tt Transfer) {
		switch t := tt.(type) {
		case *Goto:
			emit("goto " + t.Label_id.String() + ";\n")
		case *If:
			emit("if (")
			pp(t.Cond)
			emit(")\n")
			emit(" goto " + t.Truee.String() + ";\n")
			emit("else\n")
			emit(" goto " + t.Falsee.String() + ";\n")
		case *Return:
			emit("return ")
			pp(t.Op)
		default:
			panic("impossible")
		}
	}

	pp_Operand := func(oo Operand) {
		switch o := oo.(type) {
		case *Int:
			emit(strconv.Itoa(o.Value))
		case *Var:
			if o.IsField == false {
				emit(o.Name)
			} else {
				emit("this->" + o.Name)
			}
		default:
			panic("impossible")
		}
	}

	pp_Stm := func(ss Stm) {
		switch s := ss.(type) {
		case *Add:
			emit(s.Dst + " = ")
			pp(s.Left)
			emit(" + ")
			pp(s.Right)
			emit(";")
		case *And:
			emit(s.Dst + " = ")
			pp(s.Left)
			emit(" && ")
			pp(s.Right)
			emit(";")
		case *ArraySelect:
			emit(s.Name + " = ")
			pp(s.Arrayref)
			emit("[")
			pp(s.Index)
			emit("+4]")
		case *AssignArray:
			if s.IsField == false {
				emit(s.Dst + "[")
			} else {
				emit("this->" + s.Dst + "[")
			}
			pp(s.Index)
		case *InvokeVirtual:
			emit(s.Dst + " = " + s.Obj)
			emit("->vptr->" + s.F + "(" + s.Obj)
			for _, x := range s.Args {
				emit(", ")
				pp(x)
			}
			emit(");")
		case *Length:
			emit(s.Dst + " = ")
			pp(s.Arrayref)
			emit("[2]")
		case *Lt:
			emit(s.Dst + " = ")
			pp(s.Left)
			emit(" < ")
			pp(s.Right)
			emit(";")
		case *Move:
			emit(s.Dst + " = ")
			pp(s.Src)
			emit(";")
		case *NewIntArray:
			emit(s.Dst + " = (int*)Tiger_new_array(")
			pp(s.E)
			emit(")")
		case *NewObject:
			emit(s.Dst + " =((struct " + s.Class_name + "*)(Tiger_new(&" + s.Class_name +
				"_vtable_, sizeof(struct " + s.Class_name + "))));")
		case *Not:
			emit(s.Dst + " = ")
			pp(s.E)
		case *Print:
			emit("system_out_println(")
			pp(s.Args)
			emit(");")
		case *Sub:
			emit(s.Dst + " = ")
			pp(s.Left)
			emit(" - ")
			pp(s.Right)
			emit(";")
		case *Times:
			emit(s.Dst + " = ")
			pp(s.Left)
			emit(" * ")
			pp(s.Right)
		default:
			panic("impossible")
		}
	}

	pp_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			buf.Write([]byte(b.Label_id.String() + ":\\n"))
			for _, s := range b.Stms {
				pp(s)
				emit("\n")
			}
			pp(b.Trans)
		default:
			panic("impossible")
		}
	}

	pp = func(e interface{}) {
		switch v := e.(type) {
		case Block:
			pp_Block(v)
		case Transfer:
			pp_Transfer(v)
		case Operand:
			pp_Operand(v)
		case Stm:
			pp_Stm(v)
		default:
			fmt.Printf("%T\n", v)
			if vv, ok := v.(string); ok {
				fmt.Println(vv)
			}
			panic("impossible")
		}
	}

	pp(e)
	return buf.String()
}

func Visualize(p Program) {
	var f_blocks map[*util.Label]Block
	var vv func(e Acceptable)

	vv_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			f_blocks[b.Label_id] = b
		default:
			panic("impossible")
		}
	}

	vv_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			f_blocks = make(map[*util.Label]Block)
			for _, b := range m.Blocks {
				vv(b)
			}

			gname := m.ClassId + "_" + m.Name
			graph := util.Graph_new(gname)
			graph.Node_String = PP
			for _, b := range m.Blocks {
				graph.AddNode(b)
			}
			for _, b := range m.Blocks {
				if b_1, ok := b.(*BlockSingle); ok {
					switch x := b_1.Trans.(type) {
					case *Goto:
						graph.AddEdge(b, f_blocks[x.Label_id])
					case *If:
						graph.AddEdge(b, f_blocks[x.Truee])
						graph.AddEdge(b, f_blocks[x.Falsee])
					case *Return:
					default:
						panic("impossible")
					}

				}
			}
			if control.Visualize_format != control.None {
				fmt.Println("Visualize " + gname)
				graph.Visualize()
			}

		default:
			panic("impossible")
		}
	}

	vv_MainMethod := func(mm MainMethod) {
		switch m := mm.(type) {
		case *MainMethodSingle:
			f_blocks = make(map[*util.Label]Block)
			for _, b := range m.Blocks {
				vv(b)
			}
			gname := "Dog_main"
			graph := util.Graph_new(gname)
			graph.Node_String = PP

			for _, b := range m.Blocks {
				graph.AddNode(b)
			}

			for _, b := range m.Blocks {
				if b_1, ok := b.(*BlockSingle); ok {
					switch x := b_1.Trans.(type) {
					case *Goto:
						graph.AddEdge(b, x.Label_id)
					case *If:
						graph.AddEdge(b, f_blocks[x.Truee])
						graph.AddEdge(b, f_blocks[x.Falsee])
					case *Return:
					default:
						panic("impossible")
					}

				}
			}
			if control.Visualize_format != control.None {
				fmt.Println("Visualize " + gname)
				graph.Visualize()
			}
		default:
			panic("impossible")
		}
	}

	vv_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			for _, m := range p.Methods {
				vv(m)
			}
			vv(p.Main_method)
		default:
			panic("impossible")
		}
	}

	vv = func(e Acceptable) {
		switch v := e.(type) {
		case Block:
			vv_Block(v)
		case Class:
		case Dec:
		case MainMethod:
			vv_MainMethod(v)
		case Method:
			vv_Method(v)
		case Operand:
		case Program:
			vv_Program(v)
		case Stm:
		case Transfer:
		case Type:
		case Vtable:
		default:
			fmt.Printf("%T\n", v)
			panic("impossible")
		}
	}

	vv(p)
}
