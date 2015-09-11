package elaborator

import (
    "fmt"
    "os"
)
type Error_Kind int
const (
    MISTYPE = iota
    UNDECL
    RET
)

func error_mistype() {
    fmt.Printf("error> type mismatch at line: %d\n", linenum)
    fmt.Printf("need type: %s\n", current_type)
    os.Exit(0)
}

func error_undecl() {
    fmt.Printf("error> un decl var at line: %d\n", linenum)
    os.Exit(0)
}

func error_ret() {
    fmt.Printf("error> return val miss at line: %d\n", linenum)
    fmt.Printf("return type must be %s\n", current_type)
    os.Exit(0)
}

func elab_error(kind Error_Kind){
    switch kind {
    case MISTYPE:
        error_mistype()
    case UNDECL:
        error_undecl()
    case RET:
        error_ret()
    default:
        panic("impossible")
    }
}
