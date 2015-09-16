package ast_opt

import (
	"../../ast"
	"../../control"
)

func Opt(prog ast.Program) ast.Program {
	Ast := prog
	if !control.Trace_skipPass("deadclass") {
		Ast = DeadClass_new().DeadClass_Opt(Ast)
	}
	if !control.Trace_skipPass("deadcode") {
		Ast = DeadCode_new().DeadCode_Opt(Ast)
	}
	return Ast
}
