package cfg

import (
	"../codegen/C"
	"../util"
)

func TransCfg(prog codegen_c.Program) Program {
	var f_additional_locals []Dec
	var f_stm_transfer []interface{}
	var f_operand Operand
	var f_tp Type
	var f_dec Dec
	var f_vtable Vtable
	var f_method Method
	var f_main_method MainMethod
	var f_class Class

	var trans func(codegen_c.Acceptable)

	emit := func(e interface{}) {
		f_stm_transfer = append(f_stm_transfer, e)
	}

	genVar := func() string {
		fresh := util.Temp_next()
		dec := &DecSingle{&IntType{}, fresh}
		f_additional_locals = append(f_additional_locals, dec)
		return fresh
	}

	genVarT := func(t Type) string {
		fresh := util.Temp_next()
		dec := &DecSingle{t, fresh}
		f_additional_locals = append(f_additional_locals, dec)
		return fresh
	}

	//XXX enhancement
	cookBlocks := func() []Block {
		blocks := make([]Block, 0)
		stms := make([]Stm, 0)
		var transfer Transfer
		var label *util.Label
		for _, bb := range f_stm_transfer {
			switch b := bb.(type) {
			case *util.Label:
				label = b
			case Stm:
				stms = append(stms, b)
			case Transfer:
				transfer = b
				block := &BlockSingle{label, stms, transfer}
				blocks = append(blocks, block)
				stms = make([]Stm, 0) //fresh
			default:
				panic("impossible")
			}
		}
		f_stm_transfer = make([]interface{}, 0)
		return blocks
	}

	trans_Dec := func(dd codegen_c.Dec) {
		switch d := dd.(type) {
		case *codegen_c.DecSingle:
			trans(d.Tp)
			f_dec = &DecSingle{f_tp, d.Name}
		default:
			panic("impossible")
		}
	}

	trans_Type := func(tt codegen_c.Type) {
		switch t := tt.(type) {
		case *codegen_c.Int:
			f_tp = &IntType{}
		case *codegen_c.IntArray:
			f_tp = &IntArrayType{}
		case *codegen_c.ClassType:
			f_tp = &ClassType{t.Name}
		default:
			panic("impossible")
		}
	}

	trans_Exp := func(ee codegen_c.Exp) {
		switch e := ee.(type) {
		case *codegen_c.Add:
			dst := genVar()
			trans(e.Left)
			left := f_operand
			trans(e.Right)
			right := f_operand
			emit(&Add{dst, nil, left, right})
			f_operand = &Var{dst, false}
		case *codegen_c.And:
			dst := genVar()
			trans(e.Left)
			left := f_operand
			trans(e.Right)
			right := f_operand
			emit(&And{dst, left, right})
			f_operand = &Var{dst, false} //XXX IsField now unknow
		case *codegen_c.ArraySelect:
			dst := genVar()
			trans(e.Arrayref)
			array := f_operand
			trans(e.Index)
			index := f_operand
			emit(&ArraySelect{dst, array, index})
			f_operand = &Var{dst, false}
		case *codegen_c.Call:
			trans(e.RetType)
			dst := genVarT(f_tp)
			trans(e.E)
			var obj string
			operand := f_operand
			if v, ok := operand.(*Var); ok {
				if v.IsField == false {
					obj = v.Name
				} else {
					obj = "this->" + v.Name
				}
			} else {
				panic("impossible")
			}
			new_args := make([]Operand, 0)
			for _, x := range e.Args {
				trans(x)
				new_args = append(new_args, f_operand)
			}
			emit(&InvokeVirtual{dst, obj, e.Name, new_args})
			f_operand = &Var{dst, false}
		case *codegen_c.Id:
			f_operand = &Var{e.Name, e.IsField}
		case *codegen_c.Length:
			dst := genVar()
			trans(e.Arrayref)
			array := f_operand
			emit(&Length{dst, array})
			f_operand = &Var{dst, false}
		case *codegen_c.Lt:
			dst := genVar()
			trans(e.Left)
			left := f_operand
			trans(e.Right)
			right := f_operand
			emit(&Lt{dst, nil, left, right})
			f_operand = &Var{dst, false}
		case *codegen_c.NewIntArray:
			trans(e.E) //new int[E] -> a=E, b=new int[a];
			size := f_operand
			dst := genVarT(&IntArrayType{}) //int *array;
			emit(&NewIntArray{dst, size})
			f_operand = &Var{dst, false}
		case *codegen_c.NewObject:
			dst := genVarT(&ClassType{e.Class_name})
			//emit(&NewObject{"frame." + dst, e.Class_name})
			//f_operand = &Var{"frame." + dst, false}
			emit(&NewObject{dst, e.Class_name})
			f_operand = &Var{dst, false}
		case *codegen_c.Not:
			dst := genVar()
			trans(e.E)
			exp := f_operand
			emit(&Not{dst, exp})
			f_operand = &Var{dst, false}
		case *codegen_c.Num:
			f_operand = &Int{e.Value}
		case *codegen_c.Sub:
			dst := genVar()
			trans(e.Left)
			left := f_operand
			trans(e.Right)
			right := f_operand
			emit(&Sub{dst, nil, left, right})
			f_operand = &Var{dst, false}
		case *codegen_c.This:
			f_operand = &Var{"this", false}
		case *codegen_c.Times:
			dst := genVar()
			trans(e.Left)
			left := f_operand
			trans(e.Right)
			right := f_operand
			emit(&Times{dst, nil, left, right})
			f_operand = &Var{dst, false}
		default:
			panic("impossible")
		}
	}

	trans_Stm := func(ss codegen_c.Stm) {
		switch s := ss.(type) {
		case *codegen_c.Assign:
			trans(s.E)
			emit(&Move{s.Name, nil, f_operand, s.IsField})
		case *codegen_c.AssignArray:
			trans(s.Index)
			index := f_operand
			trans(s.E)
			exp := f_operand
			emit(&AssignArray{s.Name, index, exp, s.IsField})
		case *codegen_c.Block:
			for _, t := range s.Stms {
				trans(t)
			}
		case *codegen_c.If:
			t := util.Label_new()
			f := util.Label_new()
			e := util.Label_new() //exit labed
			trans(s.Cond)
			emit(&If{f_operand, t, f})
			emit(f)
			trans(s.Elsee)
			emit(&Goto{e})
			emit(t)
			trans(s.Thenn)
			emit(&Goto{e})
			emit(e)
		case *codegen_c.Print:
			trans(s.E)
			emit(&Print{f_operand})
		case *codegen_c.While:
			start := util.Label_new()
			end := util.Label_new()
			body := util.Label_new()
			emit(&Goto{start}) //XXX a hack
			emit(start)
			trans(s.Cond)
			emit(&If{f_operand, body, end})
			emit(body)
			trans(s.Body)
			emit(&Goto{start})
			emit(end)
		default:
			panic("impossible")
		}
	}

	trans_Vtable := func(vv codegen_c.Vtable) {
		switch v := vv.(type) {
		case *codegen_c.VtableSingle:
			new_ftuple := make([]*Ftuple, 0)
			for _, f := range v.Methods {
				trans(f.RetType)
				ret_type := f_tp
				args := make([]Dec, 0)
				for _, d := range f.Args {
					trans(d)
					args = append(args, f_dec)
				}
				new_ftuple = append(new_ftuple, &Ftuple{f.Classname, ret_type, args, f.Name})
			}
			f_vtable = &VtableSingle{v.Name, new_ftuple}
		default:
			panic("impossible")
		}
	}

	trans_Method := func(mm codegen_c.Method) {
		switch m := mm.(type) {
		case *codegen_c.MethodSingle:
			f_additional_locals = make([]Dec, 0)
			trans(m.RetType)
			ret_type := f_tp
			new_formals := make([]Dec, 0)
			for _, d := range m.Formals {
				trans(d)
				new_formals = append(new_formals, f_dec)
			}
			new_locals := make([]Dec, 0)
			for _, d := range m.Locals {
				trans(d)
				new_locals = append(new_locals, f_dec)
			}
			//XXX a junk label before the fiest block
			entry := util.Label_new()
			emit(entry)
			for _, s := range m.Stms {
				trans(s)
			}
			trans(m.RetExp)
			emit(&Return{f_operand})

			//TODO cookblock
			//util.Todo()
			blocks := cookBlocks()
			for _, d := range f_additional_locals {
				new_locals = append(new_locals, d)
			}
			f_method = &MethodSingle{ret_type,
				m.Name,
				m.ClassId,
				new_formals,
				new_locals,
				blocks,
				entry}

		default:
			panic("impossible")
		}
	}

	trans_MainMethod := func(mm codegen_c.MainMethod) {
		switch m := mm.(type) {
		case *codegen_c.MainMethodSingle:
			f_additional_locals = make([]Dec, 0) //fresh
			locals := make([]Dec, 0)
			for _, d := range m.Locals {
				trans(d)
				locals = append(locals, f_dec)
			}
			entry := util.Label_new()
			emit(entry)
			trans(m.Stms)
			emit(&Return{&Int{0}})
			blocks := cookBlocks()
			for _, d := range f_additional_locals {
				locals = append(locals, d)
			}
			f_main_method = &MainMethodSingle{locals, blocks}
		default:
			panic("impossible")
		}
	}

	trans_Class := func(cc codegen_c.Class) {
		switch c := cc.(type) {
		case *codegen_c.ClassSingle:
			new_tuple := make([]*Tuple, 0)
			for _, t := range c.Decs {
				trans(t.Tp)
				new_tuple = append(new_tuple,
					&Tuple{t.Classname, f_tp, t.Field_name})
			}
			f_class = &ClassSingle{c.Name, new_tuple}
		default:
			panic("impossible")
		}
	}

	trans = func(e codegen_c.Acceptable) {
		switch v := e.(type) {
		case codegen_c.Class:
			trans_Class(v)
		case codegen_c.Dec:
			trans_Dec(v)
		case codegen_c.Exp:
			trans_Exp(v)
		case codegen_c.MainMethod:
			trans_MainMethod(v)
		case codegen_c.Method:
			trans_Method(v)
		case codegen_c.Stm:
			trans_Stm(v)
		case codegen_c.Type:
			trans_Type(v)
		case codegen_c.Vtable:
			trans_Vtable(v)
		default:
			panic("impossible")
		}
	}

	var Ast *ProgramSingle
	f_stm_transfer = make([]interface{}, 0)
	if p, ok := prog.(*codegen_c.ProgramC); ok {
		new_classes := make([]Class, 0)
		for _, c := range p.Classes {
			trans(c)
			new_classes = append(new_classes, f_class)
		}
		new_vtable := make([]Vtable, 0)
		for _, v := range p.Vtables {
			trans(v)
			new_vtable = append(new_vtable, f_vtable)
		}
		new_methods := make([]Method, 0)
		for _, m := range p.Methods {
			trans(m)
			new_methods = append(new_methods, f_method)
		}
		trans(p.Mainmethod)
		new_main := f_main_method
		Ast = &ProgramSingle{new_classes, new_vtable, new_methods, new_main}
	} else {
		panic("impossible")
	}

	return Ast
}
