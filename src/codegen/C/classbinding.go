package codegen_c

type ClassBinding struct {
	extends string
	visited bool
	fields  []*Tuple
	methods []*Ftuple
}

func ClassBinding_new(extends string) *ClassBinding {
	c := new(ClassBinding)
	c.extends = extends
	c.visited = false
	c.fields = make([]*Tuple, 0)
	c.methods = make([]*Ftuple, 0)

	return c
}

func (this *ClassBinding) putField(c string, t Type, id string) {
	this.fields = append(this.fields, &Tuple{c, t, id})
}

func (this *ClassBinding) updateFields(fields []*Tuple) {
	this.fields = fields
}

func (this *ClassBinding) updateMethods(m []*Ftuple) {
	this.methods = m
}

func (this *ClassBinding) putMethod(c string, ret Type, args []Dec, mtd_name string) {
	this.methods = append(this.methods, &Ftuple{c, ret, args, mtd_name})
}
