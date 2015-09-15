package codegen_c

import (
	"fmt"
)

type Ftuple struct {
	classname string
	ret_type  Type
	args      []Dec
	id        string
}

func Ftuple_new(classname string, tp Type, args []Dec, id string) *Ftuple {
	return &Ftuple{classname, tp, args, id}
}

func Ftuple_equals(f1 *Ftuple, f2 *Ftuple) bool {
	if f1 == nil || f2 == nil {
		return false
	}
	return f1.id == f2.id
}

func (this *Ftuple) dump() {
	fmt.Printf("%s  %s", this.ret_type, this.id)
	fmt.Printf("(")
	for idx, dec := range this.args {
		if idx != 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", dec)
	}
	fmt.Printf(")")
}
