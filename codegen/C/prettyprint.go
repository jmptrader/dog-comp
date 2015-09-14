package codegen_c

import (
	"../../control"
	"../../util"
	"fmt"
	"os"
	"strconv"
)

var outputName string
var fd *os.File
var indentLevel int
var redec *util.HashSet
var classLocal map[string][]*Tuple

func say(s string) {
	fd.WriteString(s)
}

func sayln(s string) {
	say(s)
	fd.WriteString("\n")
}

func indent() {
	indentLevel += 2
}

func unIndent() {
	indentLevel -= 2
}

func printSpeaces() {
	i := indentLevel
	for i != 0 {
		say(" ")
		i--
	}
}

func outputVtable(v Vtable) {
	var vt *VtableSingle
	if vv, ok := v.(*VtableSingle); ok {
		vt = vv
	} else {
		panic("impossible")
	}
	sayln("struct " + vt.id + "_vtable " + vt.id + "_vtable_ =")
	sayln("{")
	locals := classLocal[vt.id]
	printSpeaces()
	say("\"")
	for _, t := range locals {
		_, ok := t.tp.(*ClassType)
		_, ok2 := t.tp.(*ClassType)
		if ok || ok2 {
			say("1")
		} else {
			say("0")
		}
	}
	sayln("\",")
	for _, f := range vt.methods {
		say("  ")
		sayln(f.classname + "_" + f.id + ",")
	}
	sayln("};\n")
}

func outputMainGcStack(mm MainMethod) {
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

	for _, dec := range m.locals {
		if d, ok := dec.(*DecSingle); ok {
			if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
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

func outputGcStack(mm Method) {
	var m *MethodSingle
	if v, ok := mm.(*MethodSingle); ok {
		m = v
	} else {
		panic("impossible")
	}
	sayln("struct " + m.classId + "_" + m.id + "_gc_frame")
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
	for _, dec := range m.locals {
		if d, ok := dec.(*DecSingle); ok {
			if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
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

func outputGcMap(method Method) {
	var m *MethodSingle
	if v, ok := method.(*MethodSingle); ok {
		m = v
	} else {
		panic("impossible")
	}
	say("char * " + m.classId + "_" + m.id + "_arguments_gc_map = \"")
	for _, dec := range m.formals {
		if d, ok := dec.(*DecSingle); ok {
			if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
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
	for _, dec := range m.locals {
		if d, ok := dec.(*DecSingle); ok {
			if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
				i++
			}
		} else {
			panic("impossible")
		}
	}
	sayln("int " + m.classId + "_" + m.id + "_locals_gc_map= " + strconv.Itoa(i) + ";")
	sayln("")
}

func pp_Exp_Add(e *Add) {
	pp(e.left)
	say(" + ")
	pp(e.right)
}

func pp_Exp_And(e *And) {
	pp(e.left)
	say(" && ")
	pp(e.right)
}

func pp_Exp_ArraySelect(e *ArraySelect) {
	pp(e.array)
	say("[")
	pp(e.index)
	say("+4]")
}

func pp_Exp_Call(e *Call) {
	if redec.Contains(e.assign) == false {
		say("(" + e.assign + "=")
	} else {
		say("(frame." + e.assign + "=")
	}
	pp(e.e)
	say(", ")
	if redec.Contains(e.assign) == false {
		say(e.assign + "->vptr->" + e.name + "(" + e.assign)
	} else {
		say("frame." + e.assign + "->vptr->" + e.name + "(frame." + e.assign)
	}
	for _, x := range e.args {
		say(", ") //XXX
		pp(x)
	}
	say("))")
}

func pp_Exp_Id(e *Id) {
	if e.isField == false {
		if redec.Contains(e.name) == true {
			say("frame." + e.name)
		} else {
			say(e.name)
		}
	} else {
		say("this->" + e.name)
	}
}

func pp_Exp_Length(e *Length) {
	pp(e.array)
	say("[2]")
}

func pp_Exp_Lt(e *Lt) {
	pp(e.left)
	say(" < ")
	pp(e.right)
}

func pp_Exp_NewIntArray(e *NewIntArray) {
	say("*(int*)Tiger_new_array(")
	pp(e.e)
	say(")")
}

func pp_Exp_NewObject(e *NewObject) {
	say("((struct " + e.id + "*)(Tiger_new(&" +
		e.id + "_vtable_, sizeof(struct " + e.id + "))))")
}

func pp_Exp_Not(e *Not) {
	say("!(")
	pp(e.e)
	say(")")
}

func pp_Exp_Num(e *Num) {
	say(strconv.Itoa(e.num))
}

func pp_Exp_Sub(e *Sub) {
	pp(e.left)
	say(" - ")
	pp(e.right)
}

func pp_Exp_This(e *This) {
	say("this")
}

func pp_Exp_Times(e *Times) {
	pp(e.left)
	say(" * ")
	pp(e.right)
}

func pp_Exp(exp Exp) {
	switch e := exp.(type) {
	case *Add:
		pp_Exp_Add(e)
	case *And:
		pp_Exp_And(e)
	case *ArraySelect:
		pp_Exp_ArraySelect(e)
	case *Call:
		pp_Exp_Call(e)
	case *Id:
		pp_Exp_Id(e)
	case *Length:
		pp_Exp_Length(e)
	case *Lt:
		pp_Exp_Lt(e)
	case *NewIntArray:
		pp_Exp_NewIntArray(e)
	case *NewObject:
		pp_Exp_NewObject(e)
	case *Not:
		pp_Exp_Not(e)
	case *Num:
		pp_Exp_Num(e)
	case *Sub:
		pp_Exp_Sub(e)
	case *This:
		pp_Exp_This(e)
	case *Times:
		pp_Exp_Times(e)
	default:
		panic("impossible")
	}
}

func pp_Stm_Assign(s *Assign) {
	printSpeaces()
	if s.isField == false {
		if redec.Contains(s.id) == true {
			say("frame." + s.id + " = ")
		} else {
			say(s.id + " = ")
		}
	} else {
		say("this->" + s.id + " = ")
	}
	pp(s.e)
	sayln(";")
}

func pp_Stm_AssignArray(s *AssignArray) {
	printSpeaces()
	if s.isField == false {
		if redec.Contains(s.id) == true {
			say("frame." + s.id + "[")
		} else {
			say(s.id + "[")
		}
	} else {
		say("this->" + s.id + "[")
	}
	pp(s.index)
	say("+4]")
	say(" = ")
	pp(s.e)
	sayln(";")
}

func pp_Stm_Block(s *Block) {
	printSpeaces()
	sayln("{")
	indent()
	for _, ss := range s.stms {
		pp(ss)
	}
	unIndent()
	printSpeaces()
	sayln("}")
}

func pp_Stm_If(s *If) {
	printSpeaces()
	say("if (")
	pp(s.cond)
	sayln(")")
	indent()
	pp(s.thenn)
	unIndent()
	sayln("")
	printSpeaces()
	sayln("else")
	indent()
	pp(s.elsee)
	sayln("")
	unIndent()
}

func pp_Stm_Print(s *Print) {
	printSpeaces()
	say("System_out_println(")
	pp(s.e)
	sayln(");")
}

func pp_Stm_While(s *While) {
	printSpeaces()
	say("while (")
	pp(s.cond)
	sayln(")")
	indent()
	pp(s.body)
	unIndent()
	printSpeaces()
}

func pp_Stm(stm Stm) {
	switch s := stm.(type) {
	case *Assign:
		pp_Stm_Assign(s)
	case *AssignArray:
		pp_Stm_AssignArray(s)
	case *Block:
		pp_Stm_Block(s)
	case *If:
		pp_Stm_If(s)
	case *Print:
		pp_Stm_Print(s)
	case *While:
		pp_Stm_While(s)
	default:
		panic("impossible")
	}
}

func pp_Type_ClassType(t *ClassType) {
	say("struct " + t.id + "*")
}

func pp_Type_Int(t *Int) {
	say("int")
}

func pp_Type_IntArray(t *IntArray) {
	say("int*")
}

func pp_Type(tp Type) {
	switch t := tp.(type) {
	case *Int:
		pp_Type_Int(t)
	case *IntArray:
		pp_Type_IntArray(t)
	case *ClassType:
		pp_Type_ClassType(t)
	default:
		panic("impossible")
	}
}

func pp_Dec_DecSingle(d *DecSingle) {
	pp(d.tp)
	say(" ")
	say(d.id)
}

func pp_Dec(dec Dec) {
	switch d := dec.(type) {
	case *DecSingle:
		pp_Dec_DecSingle(d)
	default:
		panic("impossible")
	}
}

func pp_Method_MethodSingle(m *MethodSingle) {
	redec.Clear()
	pp(m.retType)
	say(" " + m.classId + "_" + m.id + "(")
	for idx, dec := range m.formals {
		if idx != 0 {
			say(", ")
		}
		pp(dec)
	}
	sayln(")")

	sayln("{")
	printSpeaces()
	sayln("struct " + m.classId + "_" + m.id + "_gc_frame frame;")
	printSpeaces()
	sayln("frame.prev_ = previous;")
	printSpeaces()
	sayln("previous = &frame;")
	printSpeaces()
	sayln("frame.arguments_gc_map = " + m.classId + "_" + m.id + "_arguments_gc_map;")
	printSpeaces()
	sayln("frame.arguments_base_address = (int*)&this;")
	printSpeaces()
	sayln("frame.locals_gc_map = " + m.classId + "_" + m.id + "_locals_gc_map;")

	for _, dec := range m.locals {
		if d, ok := dec.(*DecSingle); ok {
			t := d.tp.GetType()
			printSpeaces()
			if t != TYPE_INTARRAY && t != TYPE_CLASSTYPE {
				pp(dec)
				sayln(";")
			} else {
				redec.Add(d.id)
				sayln("frame." + d.id + "=0;")
			}
		} else {
			panic("impossible")
		}
	}

	sayln("")
	for _, stm := range m.stms {
		pp(stm)
	}
	sayln("")
	printSpeaces()
	sayln("previous = frame.prev_;")

	say("  return ")
	pp(m.retExp)
	sayln(";")
	sayln("}")

}

func pp_Method(method Method) {
	switch m := method.(type) {
	case *MethodSingle:
		pp_Method_MethodSingle(m)
	default:
		panic("impossible")
	}
}

func pp_Method_MainMethodSingle(m *MainMethodSingle) {
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

	for _, dec := range m.locals {
		if d, ok := dec.(*DecSingle); ok {
			say("  ")
			t := d.tp.GetType()
			if t != TYPE_INTARRAY && t != TYPE_CLASSTYPE {
				pp(d)
				sayln(";")
			} else {
				redec.Add(d.id)
				sayln("frame." + d.id + "=0;")
			}
		} else {
			panic("impossible")
		}
	}
	indent()
	pp(m.stm)
	sayln("return 0;\n}\n")
}

func pp_MainMethod(method MainMethod) {
	switch m := method.(type) {
	case *MainMethodSingle:
		pp_Method_MainMethodSingle(m)
	default:
		panic("impossible")
	}
}

func pp_Vtable_VtableSingle(v *VtableSingle) {
	sayln("struct " + v.id + "_vtable")
	sayln("{")
	printSpeaces()
	sayln("char* " + v.id + "_gc_map;")
	for _, f := range v.methods {
		say("  ")
		pp(f.ret_type)
		say(" (*" + f.id + ")(")
		for idx, dec := range f.args {
			if idx != 0 {
				say(", ")
			}
			pp(dec)
		}
		sayln(");")
	}

	sayln("};\n")
}

func pp_Vtable(v Vtable) {
	switch vv := v.(type) {
	case *VtableSingle:
		pp_Vtable_VtableSingle(vv)
	default:
		panic("impossible")
	}
}

func pp_Class(cc Class) {
	var c *ClassSingle
	if v, ok := cc.(*ClassSingle); ok {
		c = v
	} else {
		panic("impossible")
	}
	locals := make([]*Tuple, 0)
	sayln("struct " + c.id)
	sayln("{")
	sayln("  struct " + c.id + "_vtable *vptr;")

	printSpeaces()
	sayln("int isObjOrArray;")
	printSpeaces()
	sayln("int length;")
	printSpeaces()
	sayln("void *forwarding;")

	for _, t := range c.decs {
		say("  ")
		pp(t.tp)
		say("  ")
		sayln(t.field_name + ";")
		locals = append(locals, t)
	}
	classLocal[c.id] = locals
	sayln("};")
}

func pp_Program(p Program) {
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
	for _, c := range pc.classes {
		pp(c)
	}
	sayln("// vtables")
	for _, v := range pc.vtables {
		pp(v)
	}
	sayln("")

	sayln("// method decls")
	for _, m := range pc.methods {
		if mm, ok := m.(*MethodSingle); ok {
			pp(mm.retType)
			say(" " + mm.classId + "_" + mm.id + "(")
			for idx, d := range mm.formals {
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
	for _, v := range pc.vtables {
		outputVtable(v)
	}
	sayln("")

	sayln("// GC stack frames")
	outputMainGcStack(pc.mainMethod)
	for _, method := range pc.methods {
		outputGcStack(method)
	}

	sayln("// memory GC maps")
	sayln("int Tiger_main_locals_gc_map = 1;\n")
	for _, m := range pc.methods {
		outputGcMap(m)
	}

	sayln("// methods")
	for _, m := range pc.methods {
		pp(m)
	}
	sayln("")
	sayln("// main")
	pp(pc.mainMethod)
	sayln("")
	say("\n\n")
}

func pp(e Acceptable) {
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

func pp_init() {
	indentLevel = 2
	redec = util.HashSet_new()
	classLocal = make(map[string][]*Tuple)
	if control.Control_CodeGen_outputName != "" {
		outputName = control.Control_CodeGen_outputName
	} else if control.Control_CodeGen_fileName != "" {
		outputName = control.Control_CodeGen_fileName + ".c"
	} else {
		outputName = "a.c"
	}
	f, err := os.Create(outputName)
	if err != nil {
		panic("Error> create output file error")
	}
	fd = f
}

func CodegenC(e Acceptable) {
	pp_init()
	defer fd.Close()
	pp(e)
}
