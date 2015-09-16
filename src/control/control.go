package control

var Lexer_test bool = false
var Lexer_dump bool = false

var CodeGen_fileName string = ""
var CodeGen_outputName string = ""

var CodeGen_dump bool = false

type CodeGen_Kind int

const (
	C = iota
	Bytecode
	Dalvik
	X86
)

var CodeGen_codegen CodeGen_Kind = C

var Ast_test bool = false
var Ast_dumpAst bool = false

var Elab_classTable bool = false
var Elab_methodTable bool = false

var skipedpass map[string]bool = make(map[string]bool)
var tracedpass map[string]bool = make(map[string]bool)

func Trace_add(name string){
    tracedpass[name] = true
}

func Skip_add(name string){
    skipedpass[name] = true
}

type Verbose_Kind int
const (
	VERBOSE_SILENCE = iota
	VERBOSE_PASS
	VERBOSE_SUBPASS
	VERBOSE_DETAIL
)
var Verbose Verbose_Kind = VERBOSE_SILENCE
