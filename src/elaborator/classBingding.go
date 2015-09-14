package elaborator

import (
	"../ast"
	"fmt"
)

type ClassBinding struct {
	extends string
	fields  map[string]ast.Type
	methods map[string]*MethodType
}

func ClassBinding_new(s string) *ClassBinding {
	cb := new(ClassBinding)
	cb.extends = s
	cb.fields = make(map[string]ast.Type)
	cb.methods = make(map[string]*MethodType)
	return cb
}

func (this *ClassBinding) put_FieldType(id string, tp ast.Type) {
	if this.fields[id] != nil {
		panic("duplicated class field: " + id + " at line:")
	}
	this.fields[id] = tp
}

func (this *ClassBinding) put_MethodType(mid string, tp *MethodType) {
	if this.methods[mid] != nil {
		panic("duplicated class method: " + mid)
	}
	this.methods[mid] = tp
}

func classBinding_dump(c *ClassBinding) {
	if c.extends != "" {
		fmt.Println(" extends " + c.extends)
	} else {
		fmt.Println("")
	}
	fmt.Println("   ---Field---")
	for name, t := range c.fields {
		fmt.Printf("    %s  %s\n", name, t)
	}
	fmt.Println("")
	fmt.Println("   ---Method---")
	for name, m := range c.methods {
		fmt.Printf("    %s\n", name)
		methodType_dump(m)
	}
}
