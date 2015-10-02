package cfg_opt

import (
	. "../../cfg"
	"../../control"
	"../../util"
)

func Opt(prog Program) Program {
	Ast := prog

	//dead-code
	times := 5
	for times > 0 {
		if control.Optimization_Level >= 3 {
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
		}

		if control.Optimization_Level >= 4 {
			Ast = ReachingDef(Ast)

			Ast = ConstProp(Ast)
		}

		if control.Optimization_Level >= 5 {
			Ast = CopyProp(Ast)
		}

		if control.Optimization_Level >= 6 {
			Ast = AvailExp(Ast)

			Ast = Cse(Ast)
		}

		times--
	}

	return Ast
}
