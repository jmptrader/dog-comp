package cfg_opt

import (
    . "../../cfg"
    "../../util"
)

func Opt(prog Program)Program{
    Ast := prog

    Liveness(prog)

    Ast = DeadCode(prog)

    util.Assert(Ast != nil, func(){panic("impossible")})

    return Ast
}

