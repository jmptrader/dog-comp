package codegen_c

import (
	"../../control"
	"../../util"
	"fmt"
	"os"
	"strconv"
)

func CodegenC(e Acceptable) {
	var outputName string
	var fd *os.File
	var indentLevel int
	var redec *util.HashSet
	var classLocal map[string][]*Tuple

	var pp func(e Acceptable)

	say := func(s string) {
		fd.WriteString(s)
	}

	sayln := func(s string) {
		say(s)
		fd.WriteString("\n")
	}

	indent := func() {
		indentLevel += 2
	}

	unIndent := func() {
		indentLevel -= 2
	}

	printSpeaces := func() {
		i := indentLevel
		for i != 0 {
			say(" ")
			i--
		}
	}

	outputVtable := func(v Vtable) {
		var vt *VtableSingle
		if vv, ok := v.(*VtableSingle); ok {
			vt = vv
		} else {
			panic("impossible")
		}
		sayln("struct " + vt.Name + "_vtable " + vt.Name + "_vtable_ =")
		sayln("{")
		locals := classLocal[vt.Name]
		printSpeaces()
		say("\"")
		for _, t := range locals {
			_, ok := t.Tp.(*ClassType)
			_, ok2 := t.Tp.(*ClassType)
			if ok || ok2 {
				say("1")
			} else {
				say("0")
			}
		}
		sayln("\",")
		for _, f := range vt.Methods {
			say("  ")
			sayln(f.Classname + "_" + f.Name + ",")
		}
		sayln("};\n")
	}

	outputMainGcStack := func(mm MainMethod) {
		var m *MainMethodSingle
		if v, ok := mm.(*MainMethodSingle); ok {
			m = v
		} else {
			panic("impossible")
		}
		sayln("struct Tiger_main_gc_frame")
		sayln("{")
		indent()
		printSpeaces()
		sayln("void *prev_;")
		printSpeaces()
		sayln("char *arguments_gc_map;")
		printSpeaces()
		sayln("int *arguments_base_address;")
		printSpeaces()
		sayln("int locals_gc_map;")

		for _, dec := range m.Locals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.Tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					printSpeaces()
					pp(d)
					sayln(";")
				}
			} else {
				panic("impossible")
			}
		}
		unIndent()
		sayln("};\n")
	}

	outputGcStack := func(mm Method) {
		var m *MethodSingle
		if v, ok := mm.(*MethodSingle); ok {
			m = v
		} else {
			panic("impossible")
		}
		sayln("struct " + m.ClassId + "_" + m.Name + "_gc_frame")
		sayln("{")
		indent()
		printSpeaces()
		sayln("void *prev_;")
		printSpeaces()
		sayln("char *arguments_gc_map;")
		printSpeaces()
		sayln("int *arguments_base_address;")
		printSpeaces()
		sayln("int locals_gc_map;")
		for _, dec := range m.Locals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.Tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					printSpeaces()
					pp(d)
					sayln(";")
				}
			} else {
				panic("impossible")
			}
		}
		unIndent()
		sayln("};\n")
	}

	outputGcMap := func(method Method) {
		var m *MethodSingle
		if v, ok := method.(*MethodSingle); ok {
			m = v
		} else {
			panic("impossible")
		}
		say("char * " + m.ClassId + "_" + m.Name + "_arguments_gc_map = \"")
		for _, dec := range m.Formals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.Tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					say("1")
				} else {
					say("0")
				}
			} else {
				panic("impossible")
			}
		}
		sayln("\";")
		//locals map
		i := 0
		for _, dec := range m.Locals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.Tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					i++
				}
			} else {
				panic("impossible")
			}
		}
		sayln("int " + m.ClassId + "_" + m.Name + "_locals_gc_map= " + strconv.Itoa(i) + ";")
		sayln("")
	}

	pp_Exp := func(exp Exp) {
		switch e := exp.(type) {
		case *Add:
			pp(e.Left)
			say(" + ")
			pp(e.Right)
		case *And:
			pp(e.Left)
			say(" && ")
			pp(e.Right)
		case *ArraySelect:
			pp(e.Arrayref)
			say("[")
			pp(e.Index)
			say("+4]")
		case *Call:
			if redec.Contains(e.New_id) == false {
				say("(" + e.New_id + "=")
			} else {
				say("(frame." + e.New_id + "=")
			}
			pp(e.E)
			say(", ")
			if redec.Contains(e.New_id) == false {
				say(e.New_id + "->vptr->" + e.Name + "(" + e.New_id)
			} else {
				say("frame." + e.New_id + "->vptr->" + e.Name + "(frame." + e.New_id)
			}
			for _, x := range e.Args {
				say(", ") //XXX
				pp(x)
			}
			say("))")
		case *Id:
			if e.IsField == false {
				if redec.Contains(e.Name) == true {
					say("frame." + e.Name)
				} else {
					say(e.Name)
				}
			} else {
				say("this->" + e.Name)
			}

		case *Length:
			pp(e.Arrayref)
			say("[2]")
		case *Lt:
			pp(e.Left)
			say(" < ")
			pp(e.Right)
		case *NewIntArray:
			say("(int*)Tiger_new_array(")
			pp(e.E)
			say(")")
		case *NewObject:
			say("((struct " + e.Class_name + "*)(Tiger_new(&" +
				e.Class_name + "_vtable_, sizeof(struct " + e.Class_name + "))))")
		case *Not:
			say("!(")
			pp(e.E)
			say(")")
		case *Num:
			say(strconv.Itoa(e.Value))
		case *Sub:
			pp(e.Left)
			say(" - ")
			pp(e.Right)
		case *This:
			say("this")
		case *Times:
			pp(e.Left)
			say(" * ")
			pp(e.Right)
		default:
			panic("impossible")
		}
	}

	pp_Stm := func(stm Stm) {
		switch s := stm.(type) {
		case *Assign:
			printSpeaces()
			if s.IsField == false {
				if redec.Contains(s.Name) == true {
					say("frame." + s.Name + " = ")
				} else {
					say(s.Name + " = ")
				}
			} else {
				say("this->" + s.Name + " = ")
			}
			pp(s.E)
			sayln(";")
		case *AssignArray:
			printSpeaces()
			if s.IsField == false {
				if redec.Contains(s.Name) == true {
					say("frame." + s.Name + "[")
				} else {
					say(s.Name + "[")
				}
			} else {
				say("this->" + s.Name + "[")
			}
			pp(s.Index)
			say("+4]")
			say(" = ")
			pp(s.E)
			sayln(";")

		case *Block:
			printSpeaces()
			sayln("{")
			indent()
			for _, ss := range s.Stms {
				pp(ss)
			}
			unIndent()
			printSpeaces()
			sayln("}")
		case *If:
			printSpeaces()
			say("if (")
			pp(s.Cond)
			sayln(")")
			indent()
			pp(s.Thenn)
			unIndent()
			sayln("")
			printSpeaces()
			sayln("else")
			indent()
			pp(s.Elsee)
			sayln("")
			unIndent()

		case *Print:
			printSpeaces()
			say("System_out_println(")
			pp(s.E)
			sayln(");")

		case *While:
			printSpeaces()
			say("while (")
			pp(s.Cond)
			sayln(")")
			indent()
			pp(s.Body)
			unIndent()
			printSpeaces()

		default:
			panic("impossible")
		}
	}

	pp_Type := func(tp Type) {
		switch t := tp.(type) {
		case *Int:
			say("int")
		case *IntArray:
			say("int*")
		case *ClassType:
			say("struct " + t.Name + "*")
		default:
			panic("impossible")
		}
	}

	pp_Dec := func(dec Dec) {
		switch d := dec.(type) {
		case *DecSingle:
			pp(d.Tp)
			say(" ")
			say(d.Name)
		default:
			panic("impossible")
		}
	}

	pp_Method := func(method Method) {
		switch m := method.(type) {
		case *MethodSingle:
			redec.Clear()
			pp(m.RetType)
			say(" " + m.ClassId + "_" + m.Name + "(")
			for idx, dec := range m.Formals {
				if idx != 0 {
					say(", ")
				}
				pp(dec)
			}
			sayln(")")

			sayln("{")
			printSpeaces()
			sayln("struct " + m.ClassId + "_" + m.Name + "_gc_frame frame;")
			printSpeaces()
			sayln("frame.prev_ = previous;")
			printSpeaces()
			sayln("previous = &frame;")
			printSpeaces()
			sayln("frame.arguments_gc_map = " + m.ClassId + "_" + m.Name + "_arguments_gc_map;")
			printSpeaces()
			sayln("frame.arguments_base_address = (int*)&this;")
			printSpeaces()
			sayln("frame.locals_gc_map = " + m.ClassId + "_" + m.Name + "_locals_gc_map;")

			for _, dec := range m.Locals {
				if d, ok := dec.(*DecSingle); ok {
					t := d.Tp.GetType()
					printSpeaces()
					if t != TYPE_INTARRAY && t != TYPE_CLASSTYPE {
						pp(dec)
						sayln(";")
					} else {
						redec.Add(d.Name)
						sayln("frame." + d.Name + "=0;")
					}
				} else {
					panic("impossible")
				}
			}

			sayln("")
			for _, stm := range m.Stms {
				pp(stm)
			}
			sayln("")
			printSpeaces()
			sayln("previous = frame.prev_;")

			say("  return ")
			pp(m.RetExp)
			sayln(";")
			sayln("}")

		default:
			panic("impossible")
		}
	}

	pp_MainMethod := func(method MainMethod) {
		switch m := method.(type) {
		case *MainMethodSingle:
			redec.Clear()
			sayln("int Tiger_main()")
			sayln("{")

			indent()
			printSpeaces()
			sayln("struct Tiger_main_gc_frame frame;")
			printSpeaces()
			sayln("frame.prev_ = previous;")
			printSpeaces()
			sayln("previous = &frame;")
			printSpeaces()
			sayln("frame.arguments_gc_map = 0;")
			printSpeaces()
			sayln("frame.arguments_base_address = 0;")
			printSpeaces()
			sayln("frame.locals_gc_map = Tiger_main_locals_gc_map;")
			unIndent()

			for _, dec := range m.Locals {
				if d, ok := dec.(*DecSingle); ok {
					say("  ")
					t := d.Tp.GetType()
					if t != TYPE_INTARRAY && t != TYPE_CLASSTYPE {
						pp(d)
						sayln(";")
					} else {
						redec.Add(d.Name)
						sayln("frame." + d.Name + "=0;")
					}
				} else {
					panic("impossible")
				}
			}
			indent()
			pp(m.Stms)
			sayln("return 0;\n}\n")

		default:
			panic("impossible")
		}
	}

	pp_Vtable := func(v Vtable) {
		switch vv := v.(type) {
		case *VtableSingle:
			sayln("struct " + vv.Name + "_vtable")
			sayln("{")
			printSpeaces()
			sayln("char* " + vv.Name + "_gc_map;")
			for _, f := range vv.Methods {
				say("  ")
				pp(f.RetType)
				say(" (*" + f.Name + ")(")
				for idx, dec := range f.Args {
					if idx != 0 {
						say(", ")
					}
					pp(dec)
				}
				sayln(");")
			}

			sayln("};\n")

		default:
			panic("impossible")
		}
	}

	pp_Class := func(cc Class) {
		var c *ClassSingle
		if v, ok := cc.(*ClassSingle); ok {
			c = v
		} else {
			panic("impossible")
		}
		locals := make([]*Tuple, 0)
		sayln("struct " + c.Name)
		sayln("{")
		sayln("  struct " + c.Name + "_vtable *vptr;")

		printSpeaces()
		sayln("int isObjOrArray;")
		printSpeaces()
		sayln("int length;")
		printSpeaces()
		sayln("void *forwarding;")

		for _, t := range c.Decs {
			say("  ")
			pp(t.Tp)
			say("  ")
			sayln(t.Field_name + ";")
			locals = append(locals, t)
		}
		classLocal[c.Name] = locals
		sayln("};")
	}

	pp_Program := func(p Program) {
		sayln("// This is automatically generated by the Dog compiler.")
		sayln("// Do NOT modify!\n")
		sayln("extern void *previous;")
		sayln("extern void *Tiger_new_array(int);")
		sayln("extern void *Tiger_new(void*, int);")
		sayln("extern int System_out_println(int);")

		var pc *ProgramC
		if v, ok := p.(*ProgramC); ok {
			pc = v
		} else {
			panic("impossible")
		}

		sayln("// strutures")
		for _, c := range pc.Classes {
			pp(c)
		}
		sayln("// vtables")
		for _, v := range pc.Vtables {
			pp(v)
		}
		sayln("")

		sayln("// method decls")
		for _, m := range pc.Methods {
			if mm, ok := m.(*MethodSingle); ok {
				pp(mm.RetType)
				say(" " + mm.ClassId + "_" + mm.Name + "(")
				for idx, d := range mm.Formals {
					if idx != 0 {
						say(", ")
					}
					pp(d)
				}
				sayln(");")
			} else {
				panic("impossible")
			}
		}

		sayln("// vtables")
		for _, v := range pc.Vtables {
			outputVtable(v)
		}
		sayln("")

		sayln("// GC stack frames")
		outputMainGcStack(pc.Mainmethod)
		for _, method := range pc.Methods {
			outputGcStack(method)
		}

		sayln("// memory GC maps")
		sayln("int Tiger_main_locals_gc_map = 1;\n")
		for _, m := range pc.Methods {
			outputGcMap(m)
		}

		sayln("// methods")
		for _, m := range pc.Methods {
			pp(m)
		}
		sayln("")
		sayln("// main")
		pp(pc.Mainmethod)
		sayln("")
		say("\n\n")
	}

	pp = func(e Acceptable) {
		switch v := e.(type) {
		case Class:
			pp_Class(v)
		case Dec:
			pp_Dec(v)
		case Exp:
			pp_Exp(v)
		case MainMethod:
			pp_MainMethod(v)
		case Method:
			pp_Method(v)
		case Program:
			pp_Program(v)
		case Stm:
			pp_Stm(v)
		case Type:
			pp_Type(v)
		case Vtable:
			pp_Vtable(v)
		default:
			fmt.Println(v)
			panic("impossible")
		}
	}

	pp_init := func() {
		indentLevel = 2
		redec = util.HashSet_new()
		classLocal = make(map[string][]*Tuple)
		if control.CodeGen_outputName != "" {
			outputName = control.CodeGen_outputName
		} else if control.CodeGen_fileName != "" {
			outputName = control.CodeGen_fileName + ".c"
		} else {
			outputName = "a.c"
		}
		f, err := os.Create(outputName)
		if err != nil {
			panic("Error> create output file error")
		}
		fd = f
	}

	pp_init()
	defer fd.Close()
	pp(e)
}
