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

func ParserError(expect string, current string, linenum int){
    fmt.Printf("Expect: <%s>, but got <%s> at line:%d\n", expect, current, linenum)
    os.Exit(0)
}


