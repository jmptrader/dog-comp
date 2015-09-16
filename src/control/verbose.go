package control

import (
	"fmt"
)

const (
	VERBOSE_SILENCE = iota
	VERBOSE_PASS
	VERBOSE_SUBPASS
	VERBOSE_DETAIL
)

var Verbose_Kind int = VERBOSE_SILENCE

func order(l int) bool {
	if l <= Verbose_Kind {
		return true
	} else {
		return false
	}
}

func Verbose(s string, f func(), level int) {
	if order(level) {
		Trace_spaces()
		fmt.Println(s + " starting")
		Trace_indent()
	}
	f()
	if order(level) {
		Trace_unIndent()
		Trace_spaces()
		fmt.Println(s + " ending")
	}
}
