package util

import (
	"fmt"
	"os"
	"path"
)

func Bug(info string, filename string, linenum int) {
	fmt.Printf("ERROR>%s:%d:%s\n", path.Base(filename), linenum, info)
	os.Exit(0)
}
