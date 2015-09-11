package elaborator

import (
    "../ast"
)

type MethodType struct {
    retType ast.Type
    argsType []ast.Dec
}
