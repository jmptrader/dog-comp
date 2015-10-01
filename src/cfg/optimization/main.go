package cfg_opt

import (
	. "../../cfg"
	"../../control"
	"../../util"
)

func Opt(prog Program) Program {
	Ast := prog

	//dead-code
	if control.Optimization_Level < 3 {
		return Ast
	}
	Liveness(Ast)
	in_size := len(stmLiveIn)
	out_size := len(stmLiveOut)
	Ast = DeadCode(Ast)
	for {
		Liveness(Ast)
		if in_size == len(stmLiveIn) && out_size == len(stmLiveOut) {
			break
		}
		in_size = len(stmLiveIn)
		out_size = len(stmLiveOut)
		Ast = DeadCode(Ast)
	}
	util.Assert(Ast != nil, func() { panic("impossible") })

	Ast = ReachingDef(Ast)

	Ast = ConstProp(Ast)

	Ast = CopyProp(Ast)

    Ast = AvailExp(Ast)

	return Ast
}
