package util

import (
	"strconv"
)

var label_count int

type Label struct {
	i int
}

func Label_new() *Label {
	o := new(Label)
	o.i = label_count
	label_count++

	return o
}

func (this *Label) String() string {
	return "L_" + strconv.Itoa(this.i)
}
