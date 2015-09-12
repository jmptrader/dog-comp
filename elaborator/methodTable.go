package elaborator

import (
    "../ast"
    "fmt"
)
var method_table map[string]ast.Type

func initMethodTable() {
    method_table = make(map[string]ast.Type)
}

func mt_get(id string)ast.Type {
    return method_table[id]
}

func mt_put(formals []ast.Dec, locals []ast.Dec) {
    f := func (list []ast.Dec) {
        for _, dec :=  range list {
            switch v := dec.(type) {
            case *ast.DecSingle:
                if method_table[v.Name] != nil {
                    panic("duplicated parameter: "+ v.Name)
                }else {
                    method_table[v.Name] = v.Tp
                }
            default:
                panic("wrong type")
            }
        }
    }
    f(formals)
    f(locals)
}

func methodTable_dump() {
    for name, t := range method_table {
        fmt.Println(name)
        fmt.Print(" ")
        fmt.Println(t)
        fmt.Print("\n")
    }
}
