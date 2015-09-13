package elaborator

import (
	"../ast"
	"fmt"
)

type MethodType struct {
	retType  ast.Type
	argsType []ast.Dec
}

func methodType_dump(m *MethodType) {
	fmt.Printf("        retType: ")
	fmt.Println(m.retType)
	for _, dec := range m.argsType {
		fmt.Printf("        ")
		fmt.Println(dec)
	}
}
