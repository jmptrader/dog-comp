package util

import (
	"../control"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Dot struct {
	list []*DotElement
}

func Dot_new() *Dot {
	o := new(Dot)
	o.list = make([]*DotElement, 0)
	return o
}

func (this *Dot) Insert(from, to string) {
	this.list = append(this.list, this.DotElement_new(from, to, ""))
}

func (this *Dot) String() string {
	var buf bytes.Buffer
	for _, e := range this.list {
		buf.Write([]byte(e.String()))
	}
	return buf.String()
}

func (this *Dot) toDot(name string) {
	fname := name + ".dot"
	fd, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	var buf bytes.Buffer
	buf.Write([]byte("digraph G{\n"))
	buf.Write([]byte("\tsize = \"10,10\";\n"))
	buf.Write([]byte("\tnode [color=lightblue2, style=filled];\n"))
	buf.Write([]byte(this.String()))
	buf.Write([]byte("}\n"))
	fd.WriteString(buf.String())
}

func (this *Dot) Visualize(name string) {
	this.toDot(name)
	format := ""
	postfix := ""
	switch control.Visualize_format {
	case control.None:
		format = "-Tsvg"
		postfix = "svg"
	case control.Pdf:
		format = "-Tpdf"
		postfix = "pdf"
	case control.Ps:
		format = "-Tps"
		postfix = "ps"
	case control.Jpg:
		format = "-Tjpg"
		postfix = "jpg"
	case control.Svg:
		format = "-Tsvg"
		postfix = "svg"
	default:
		panic("impossible")
	}
	cmd := exec.Command("dot", format, name+".dot", "-o", name+"."+postfix)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	fmt.Print(stdout.String())
	fmt.Print(stderr.String())

	stdout.Reset() //reset the buffer
	stderr.Reset()
	cmd_delet := exec.Command("rm", name+".dot")
	cmd_delet.Stdout = &stdout
	cmd_delet.Stderr = &stderr
	cmd_delet.Run()
	fmt.Print(stdout.String())
	fmt.Print(stderr.String())

}

type DotElement struct {
	X string
	Y string
	Z string
}

func (this *Dot) DotElement_new(x, y, z string) *DotElement {
	return &DotElement{x, y, z}
}

//FIXME x,y,z need assert
func (this *DotElement) String() string {
	s := ""
	if this.Z != "" {
		s = this.Z
	}
	return "\"" + this.X + "\"" +
		"->" +
		"\"" + this.Y + "\"" +
		s + ";\n"
}
