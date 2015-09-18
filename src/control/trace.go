package control

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

func Trace_skipPass(name string) bool {
	return skipedpass[name]
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

func Trace(name string, f func(), dox func(), dor func()) {
	if Trace_contains(name) {
		dox()
	}
	f()
	if Trace_contains(name) {
		dor()
	}

}
