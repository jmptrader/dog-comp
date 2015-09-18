package ast_opt

import (
	"../../ast"
	"../../control"
)

func Opt(prog ast.Program) ast.Program {
	Ast := prog
	before_opt := Statistics_Ast(Ast)
	if !control.Trace_skipPass("deadclass") {
		control.Verbose("DeadClass-Opt", func() {
			control.Trace("deadclass", func() {
				Ast = DeadClass_new().DeadClass_Opt(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			})
		}, control.VERBOSE_SUBPASS)
	}
	if !control.Trace_skipPass("deadcode") {
		control.Verbose("DeadCode-Opt", func() {
			control.Trace("deadcode", func() {
				Ast = DeadCode_new().DeadCode_Opt(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			})
		}, control.VERBOSE_SUBPASS)
	}
	if !control.Trace_skipPass("algsimp") {
		control.Verbose("AlgSimp-Opt", func() {
			control.Trace("algsimp", func() {
				Ast = AlgSimp(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			})
		}, control.VERBOSE_SUBPASS)
	}
	if !control.Trace_skipPass("constfold") {
		control.Verbose("ConstFold-Opt", func() {
			control.Trace("constfold", func() {
				Ast = ConstFold(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			}, func() {
				ast.NewPP().DumpProg(Ast)
			})
		}, control.VERBOSE_SUBPASS)
	}
	after_opt := Statistics_Ast(Ast)
	if after_opt < before_opt {
		Ast = Opt(Ast)
	}
	return Ast
}
