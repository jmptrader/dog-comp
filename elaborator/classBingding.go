package elaborator

import (
    "../ast"
)

type ClassBinding struct{
    Fields map[string]ast.Type
    Methods map[string]ast.Func
}

func NewClassBinding() *ClassBinding{
    cb := new(ClassBinding)
    cb.Fields = make(map[string]ast.Type)
    cb.Methods = make(map[string]ast.Func)

    return cb
}
