package util

import (
	"os"
)

func Assert(cond bool, f func()) {
	if !cond {
		if f != nil {
			f()
		}
		os.Exit(1)
	}
}
