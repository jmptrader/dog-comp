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
			emit("goto " + t.label.String() + ";\n")
		case *If:
			emit("if (")
			pp(t.cond)
			emit(")\n")
			emit(" goto " + t.truee.String() + ";\n")
			emit("else\n")
			emit(" goto " + t.falsee.String() + ";\n")
		case *Return:
			emit("return ")
			pp(t.op)
		default:
			panic("impossible")
		}
	}

	pp_Operand := func(oo Operand) {
		switch o := oo.(type) {
		case *Int:
			emit(strconv.Itoa(o.value))
		case *Var:
			if o.isField == false {
				emit(o.id)
			} else {
				emit("this->" + o.id)
			}
		default:
			panic("impossible")
		}
	}

	pp_Stm := func(ss Stm) {
		switch s := ss.(type) {
		case *Add:
			emit(s.dst + " = ")
			pp(s.left)
			emit(" + ")
			pp(s.right)
			emit(";")
		case *And:
			emit(s.dst + " = ")
			pp(s.left)
			emit(" && ")
			pp(s.right)
			emit(";")
		case *ArraySelect:
			emit(s.id + " = ")
			pp(s.array)
			emit("[")
			pp(s.index)
			emit("+4]")
		case *AssignArray:
			if s.isField == false {
				emit(s.dst + "[")
			} else {
				emit("this->" + s.dst + "[")
			}
			pp(s.index)
		case *InvokeVirtual:
			emit(s.dst + " = " + s.obj)
			emit("->vptr->" + s.f + "(" + s.obj)
			for _, x := range s.args {
				emit(", ")
				pp(x)
			}
			emit(");")
		case *Length:
			emit(s.dst + " = ")
			pp(s.array)
			emit("[2]")
		case *Lt:
			emit(s.dst + " = ")
			pp(s.left)
			emit(" < ")
			pp(s.right)
			emit(";")
		case *Move:
			emit(s.dst + " = ")
			pp(s.src)
			emit(";")
		case *NewIntArray:
			emit(s.dst + " = (int*)Tiger_new_array(")
			pp(s.exp)
			emit(")")
		case *NewObject:
			emit(s.dst + " =((struct " + s.c + "*)(Tiger_new(&" + s.c +
				"_vtable_, sizeof(struct " + s.c + "))));")
		case *Not:
			emit(s.dst + " = ")
			pp(s.exp)
		case *Print:
			emit("system_out_println(")
			pp(s.arg)
			emit(");")
		case *Sub:
			emit(s.dst + " = ")
			pp(s.left)
			emit(" - ")
			pp(s.right)
			emit(";")
		case *Times:
			emit(s.dst + " = ")
			pp(s.left)
			emit(" * ")
			pp(s.right)
		default:
			panic("impossible")
		}
	}

	pp_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			buf.Write([]byte(b.label.String() + ":\\n"))
			for _, s := range b.stms {
				pp(s)
				emit("\n")
			}
			pp(b.transfer)
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
			f_blocks[b.label] = b
		default:
			panic("impossible")
		}
	}

	vv_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			f_blocks = make(map[*util.Label]Block)
			for _, b := range m.blocks {
				vv(b)
			}

			gname := m.classId + "_" + m.name
			graph := util.Graph_new(gname)
			graph.Node_String = PP
			for _, b := range m.blocks {
				graph.AddNode(b)
			}
			for _, b := range m.blocks {
				if b_1, ok := b.(*BlockSingle); ok {
					switch x := b_1.transfer.(type) {
					case *Goto:
						graph.AddEdge(b, f_blocks[x.label])
					case *If:
						graph.AddEdge(b, f_blocks[x.truee])
						graph.AddEdge(b, f_blocks[x.falsee])
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
			for _, b := range m.blocks {
				vv(b)
			}
			gname := "Dog_main"
			graph := util.Graph_new(gname)
			graph.Node_String = PP

			for _, b := range m.blocks {
				graph.AddNode(b)
			}

			for _, b := range m.blocks {
				if b_1, ok := b.(*BlockSingle); ok {
					switch x := b_1.transfer.(type) {
					case *Goto:
						graph.AddEdge(b, x.label)
					case *If:
						graph.AddEdge(b, f_blocks[x.truee])
						graph.AddEdge(b, f_blocks[x.falsee])
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
			for _, m := range p.methods {
				vv(m)
			}
			vv(p.main_method)
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
			panic("impossible")
		}
	}

	vv(p)
}
