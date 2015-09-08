package control

var Control_Lexer_test bool = false
var Control_Lexer_dump bool = false

var Control_CodeGen_fileName string = ""
var Control_CodeGen_outputName string = ""

type CodeGen_Kind int

const (
	C = iota
	Bytecode
	Dalvik
	X86
)

var Control_CodeGen_codegen CodeGen_Kind = C

var Control_Ast_test bool = false
var Control_Ast_dumpAst bool = false
