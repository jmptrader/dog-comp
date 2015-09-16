package util

import (
	"fmt"
)

const (
	STEP = 2
)

var indent int = 0
var traceSet map[string]bool = make(map[string]bool)
var skipedpass map[string]bool = make(map[string]bool)

func Trace_Skip_add(name string) {
	skipedpass[name] = true
}

func Trace_indent() {
	indent += STEP
}

func Trace_unIndent() {
	indent -= STEP
}

func Trace_spaces() {
	i := indent
	for i > 0 {
		fmt.Print(" ")
		i--
	}
}

func Trace_contains(s string) bool {
	return traceSet[s]
}

func Trace_add(name string) {
	traceSet[name] = true
}

//FIXME
func Trace(name string,
	f func(x interface{}) interface{},
	x interface{},
	dox func(c interface{}),
	r interface{},
	dor func(c interface{})) {
	if Trace_contains(name) == true {
		dox(x)
	}
	r = f(x)
	if Trace_contains(name) == true {
		dor(r)
	}
}
