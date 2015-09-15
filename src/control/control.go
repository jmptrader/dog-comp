package control

var Lexer_test bool = false
var Lexer_dump bool = false

var CodeGen_fileName string = ""
var CodeGen_outputName string = ""

type CodeGen_Kind int

const (
	C = iota
	Bytecode
	Dalvik
	X86
)

var CodeGen_codegen CodeGen_Kind = C
var CodeGen_dump bool = false

var Ast_test bool = false
var Ast_dumpAst bool = false

var Elab_classTable bool = false
var Elab_methodTable bool = false
