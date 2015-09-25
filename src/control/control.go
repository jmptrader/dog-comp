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

const (
	None = iota
	Pdf
	Ps
	Jpg
	Svg
)

var Visualize_format int = None

var Optimization_Level int = 6
