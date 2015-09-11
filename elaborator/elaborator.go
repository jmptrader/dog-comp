package elaborator

import (
    "../ast"
)

var current_class string
var current_method string
var current_type ast.Type
var linenum int

func elabMainClass(mc ast.MainClass){
    switch m := mc.(type) {
    case *ast.MainClassSingle:
        current_class = m.Name
        elaborate(m.Stms)
    default:
        panic("wrong type")
    }
}

func elabClass (class ast.Class){
    switch c := class.(type) {
    case *ast.ClassSingle:
        current_class = c.Name
        for _, m := range c.Methods {
            elaborate(m)
        }
    default:
        panic("wrong type")
    }
}

func elabMethod(mth ast.Method){
    if m, ok := mth.(*ast.MethodSingle); ok {
        initMethodTable()
        mt_put(m.Formals, m.Locals)
        for _, stm := range m.Stms {
            elaborate(stm)
        }
        cb := ct_get(current_class);
        mtd_type := cb.methods[m.Name]
        elaborate(m.RetExp)//modify current_type
        if mtd_type.retType.Gettype() != current_type.Gettype() {
            elab_error(RET)
        }
    }else {
        panic("wrong type")
    }
}

func elabStm_Assign(s *ast.Assign) {
    tp := mt_get(s.Name)
    if tp == nil {
        tp = ct_getFieldType(current_class, s.Name)
        s.IsField = true
    }
    if tp == nil {
        elab_error(UNDECL)
    }
    s.Tp = tp
    elaborate(s.E)
}

func elabStm_AssignArray(s *ast.AssignArray){
    tp := mt_get(s.Name)
    if tp == nil {
        tp = ct_getFieldType(current_class, s.Name)
        s.IsField = true
    }
    if tp == nil {
        elab_error(UNDECL)
    }
    s.Tp = tp
    elaborate(s.Index)
    if _, ok := current_type.(*ast.Int); !ok {
        elab_error(MISTYPE)
    }
    elaborate(s.E)
    //TODO
}

func elabStm_Block(s *ast.Block){
    for _, t := range s.Stms {
        elaborate(t)
    }
}

func elabStm_If(s *ast.If) {
    elaborate(s.Condition)
    if _, ok := current_type.(*ast.Boolean); !ok {
        elab_error(MISTYPE)
    }
    elaborate(s.Thenn)
    elaborate(s.Elsee)
}

func elabStm_Print(s *ast.Print) {
    elaborate(s.E)
    if _, ok := current_type.(*ast.Int); !ok {
        elab_error(MISTYPE)
    }
}

func elabStm_While(s *ast.While) {
    elaborate(s.E)
    if _, ok := current_type.(*ast.Boolean); !ok {
        elab_error(MISTYPE)
    }
    elaborate(s.Body)
}

func elabStm(stm ast.Stm){
    switch s := stm.(type) {
    case *ast.Assign:
        linenum = s.LineNum
        elabStm_Assign(s)
    case *ast.AssignArray:
        linenum = s.LineNum
        elabStm_AssignArray(s)
    case *ast.Block:
        linenum = s.LineNum
        elabStm_Block(s)
    case *ast.If:
        linenum = s.LineNum
        elabStm_If(s)
    case *ast.Print:
        linenum = s.LineNum
        elabStm_Print(s)
    case *ast.While:
        linenum = s.LineNum
        elabStm_While(s)
    default:
        panic("wrong type")
    }
}

/**
 * recursive compare
 */
func compareClass(t1 string, t2 string)bool {
    if t2 == "" {
        return false
    }
    if t1 == t2 {
        return true
    }else {
        super := ct_get(t2).extends
        return compareClass(t1, super)
    }
}

func elabExp_Call(e *ast.Call) {
    var ty *ast.ClassType
    elaborate(e.Callee)
    left_type := current_type
    if t, ok := left_type.(*ast.ClassType); ok {
        ty = t
        e.Firsttype = t.Name
    }else {
        elab_error(MISTYPE)
    }
    mtd_type := ct_getMethodType(ty.Name, e.MethodName)
    var args_ty []ast.Type
    for _, a := range e.ArgsList {
        elaborate(a)
        args_ty = append(args_ty, current_type)
    }
    //setp 1:check args number
    if len(mtd_type.argsType) != len(args_ty) {
        elab_error(MISTYPE)
    }
    //setp 2:check args type
    for i, arg := range args_ty {
        if arg.Gettype() != mtd_type.argsType[i].GetDecType() {
            elab_error(MISTYPE)
        }
        if dec, ok := mtd_type.argsType[i].(*ast.DecSingle); ok{
            if t1, ok := arg.(*ast.ClassType); ok{
                if t2, ok2 := dec.Tp.(*ast.ClassType); ok2 {
                    if !compareClass(t2.Name,t1.Name) {
                        elab_error(MISTYPE)
                    }
    //setp 3:rewirte ast.Call.ArgsType
                    args_ty[i] = t2
                }
            }
        }
    }
    current_type = mtd_type.retType
    e.ArgsType = args_ty
    e.Rt = current_type
}

func elabExp_Add(e *ast.Add) {
    elaborate(e.Left)
    left_type := current_type
    elaborate(e.Right)
    if left_type.Gettype() != ast.TYPE_INT {
        elab_error(MISTYPE)
    }
    if left_type.Gettype() != current_type.Gettype() {
        elab_error(MISTYPE)
    }
}

func elabExp_And(e *ast.And) {
    elaborate(e.Left)
    left_type := current_type
    elaborate(e.Right)
    if left_type.Gettype() != ast.TYPE_BOOLEAN {
        elab_error(MISTYPE)
    }
    if left_type.Gettype() != current_type.Gettype() {
        elab_error(MISTYPE)
    }
}

func elabExp_Sub(e *ast.Sub) {
    elaborate(e.Left)
    left_type := current_type
    elaborate(e.Right)
    if left_type.Gettype() != ast.TYPE_INT {
        elab_error(MISTYPE)
    }
    if left_type.Gettype() != current_type.Gettype() {
        elab_error(MISTYPE)
    }
}
func elabExp_Times(e *ast.Times) {
    elaborate(e.Left)
    left_type := current_type
    elaborate(e.Right)
    if left_type.Gettype() != ast.TYPE_INT {
        elab_error(MISTYPE)
    }
    if left_type.Gettype() != current_type.Gettype() {
        elab_error(MISTYPE)
    }
}
func elabExp_Lt(e *ast.Lt) {
    elaborate(e.Left)
    left_type := current_type
    elaborate(e.Right)
    if left_type.Gettype() != ast.TYPE_INT {
        elab_error(MISTYPE)
    }
    if left_type.Gettype() != current_type.Gettype() {
        elab_error(MISTYPE)
    }
    current_type = &ast.Boolean{ast.TYPE_BOOLEAN}
}

func elabExp_ArraySelect(e *ast.ArraySelect) {
    elaborate(e.Index)
    if _, ok := current_type.(*ast.Int); !ok {
        elab_error(MISTYPE)
    }
    elaborate(e.Arrayref)
    current_type = &ast.Int{ast.TYPE_INT}
}

func elabExp_Id(e *ast.Id) {
    tp := mt_get(e.Name)
    if tp == nil {
        tp = ct_getFieldType(current_class, e.Name)
        e.IsField = true
    }
    if tp == nil {
        elab_error(UNDECL)
    }
    current_type = tp
    e.Tp = tp
}

func elabExp_Length(e *ast.Length){
    elaborate(e.Arrayref)
    current_type = &ast.Int{ast.TYPE_INT}
}

func elabExp_NewIntArray(e *ast.NewIntArray){
    elaborate(e.Size)
    if _, ok := current_type.(*ast.Int); !ok {
        elab_error(MISTYPE)
    }
    current_type = &ast.IntArray{ast.TYPE_INTARRAY}
}
func elabExp(exp ast.Exp) {
    switch e := exp.(type) {
    case *ast.Add:
        linenum = e.LineNum
        elabExp_Add(e)
    case *ast.And:
        linenum = e.LineNum
        elabExp_And(e)
    case *ast.ArraySelect:
        linenum = e.LineNum
        elabExp_ArraySelect(e)
    case *ast.Call:
        linenum = e.LineNum
        elabExp_Call(e)
    case *ast.False:
        linenum = e.LineNum
        current_type = &ast.Boolean{ast.TYPE_BOOLEAN}
    case *ast.Id:
        linenum = e.LineNum
        elabExp_Id(e)
    case *ast.Length:
        linenum = e.LineNum
        elabExp_Length(e)
    case *ast.Lt:
        linenum = e.LineNum
        elabExp_Lt(e)
    case *ast.NewIntArray:
        linenum = e.LineNum
        elabExp_NewIntArray(e)
    case *ast.NewObject:
        linenum = e.LineNum
        current_type = &ast.ClassType{e.Name, ast.TYPE_CLASS}
    case *ast.Not:
        linenum = e.LineNum
        elaborate(e.E)
        current_type = &ast.Boolean{ast.TYPE_BOOLEAN}
    case *ast.Num:
        linenum = e.LineNum
        current_type = &ast.Int{ast.TYPE_INT}
    case *ast.Sub:
        linenum = e.LineNum
        elabExp_Sub(e)
    case *ast.This:
        linenum = e.LineNum
        current_type = &ast.ClassType{current_class, ast.TYPE_CLASS}
    case *ast.Times:
        linenum = e.LineNum
        elabExp_Times(e)
    case *ast.True:
        linenum = e.LineNum
        current_type = &ast.Boolean{ast.TYPE_BOOLEAN}
    default:
        panic("wrong type")
    }
}

func buildMainClass(c ast.MainClass) {
    switch v := c.(type) {
    case *ast.MainClassSingle:
        ct_put(v.Name, &ClassBinding{})
    default:
        panic("wrong type")
    }
}

func buildClass(c ast.Class) {
    switch v := c.(type){
    case *ast.ClassSingle:
        ct_put(v.Name, ClassBinding_new(v.Extends))
        for _, dec := range v.Decs {
            switch d:= dec.(type) {
            case *ast.DecSingle:
                ct_putFieldType(v.Name, d.Name, d.Tp)
            default:
                panic("wrong type")
            }
        }
        for _, mtd := range v.Methods {
            switch m:=mtd.(type) {
            case *ast.MethodSingle:
                ct_putMethodType(v.Name, m.Name, &MethodType{m.RetType, m.Formals})
            default:
                panic("wrong type")
            }
        }
    default:
        panic("wrong type")
    }
}

func elabProg(p ast.Program) {
    switch v:=p.(type){
    case *ast.ProgramSingle:
        buildMainClass(v.Mainclass)
        for _, c := range v.Classes{
            buildClass(c)
        }
        elabMainClass(v.Mainclass)
        for _, c :=range v.Classes {
            elaborate(c)
        }
        for _, c := range v.Classes {
            elaborate(c)
        }
    default:
        panic("wrong type")
    }
}

func elaborate(e ast.Acceptable){
    switch v := e.(type){
    case ast.Exp:
        elabExp(v)
    case ast.Stm:
        elabStm(v)
    case ast.Method:
        elabMethod(v)
    case ast.Class:
        elabClass(v)
    case ast.MainClass:
        elabMainClass(v)
    default:
        panic("wrong type")
    }
}

func Elaborate(e ast.Acceptable) {
    current_class = ""
    current_method = ""
    current_type = nil
    initClassTable()

    switch v := e.(type){
    case ast.Program:
        elabProg(v)
    default:
        panic("wrong type")
    }
}
