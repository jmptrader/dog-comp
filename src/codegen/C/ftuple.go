package codegen_c

import (
	"fmt"
)

type Ftuple struct {
	Classname string
	RetType   Type
	Args      []Dec
	Name      string
}

func Ftuple_new(classname string, tp Type, args []Dec, id string) *Ftuple {
	return &Ftuple{classname, tp, args, id}
}

func Ftuple_equals(f1 *Ftuple, f2 *Ftuple) bool {
	if f1 == nil || f2 == nil {
		return false
	}
	return f1.Name == f2.Name
}

func (this *Ftuple) dump() {
	fmt.Printf("%s  %s", this.RetType, this.Name)
	fmt.Printf("(")
	for idx, dec := range this.Args {
		if idx != 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", dec)
	}
	fmt.Printf(")")
}
