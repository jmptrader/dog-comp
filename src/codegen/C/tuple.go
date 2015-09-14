package codegen_c

type Tuple struct {
	classname  string
	tp         Type
	field_name string
}

func Tuple_new(c string, t Type, name string) *Tuple {
	return &Tuple{c, t, name}
}

func Tuple_equals(t1 *Tuple, t2 *Tuple) bool {
	if t1 == nil || t2 == nil {
		return false
	}
	if t1.field_name == t2.field_name {
		return true
	}
	return false
}
