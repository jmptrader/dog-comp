package codegen_c

import (
	"fmt"
)

type Tuple struct {
	Classname  string
	Tp         Type
	Field_name string
}

func Tuple_new(c string, t Type, name string) *Tuple {
	return &Tuple{c, t, name}
}

func Tuple_equals(t1 *Tuple, t2 *Tuple) bool {
	if t1 == nil || t2 == nil {
		return false
	}
	if t1.Field_name == t2.Field_name {
		return true
	}
	return false
}

func (this *Tuple) dump() {
	fmt.Printf("%s  %s  :%s", this.Field_name, this.Tp, this.Classname)
}
