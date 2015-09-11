package elaborator
import (
    "../ast"
)

var class_table map[string]*ClassBinding

func initClassTable() {
    class_table = make(map[string]*ClassBinding)
}

func ct_put(id string, cb *ClassBinding) {
    if class_table[id] != nil {
        panic("duplicated class: "+id)
    }
    class_table[id] = cb
}

func ct_get(id string)*ClassBinding {
    return class_table[id]
}

func ct_getFieldType(class_name string, id string)ast.Type{
    cb := class_table[class_name]
    tp := cb.fields[id]
    if tp != nil {
        return tp
    }
    if cb.extends == "" {
        return tp
    }else {
        return ct_getFieldType(cb.extends, id)
    }
}

func ct_putFieldType(c string, id string, tp ast.Type) {
    cb := class_table[c]
    cb.put_FieldType(id, tp)
}

func ct_getMethodType(class_name string, mid string)*MethodType {
    cb := class_table[class_name]
    tp := cb.methods[mid]
    if tp != nil {
        return tp
    }
    if cb.extends == "" {
        return tp
    }else {
        return ct_getMethodType(cb.extends, mid)
    }
}

func ct_putMethodType(class_name string, mid string, tp *MethodType) {
    cb := class_table[class_name]
    cb.put_MethodType(mid, tp)
}
