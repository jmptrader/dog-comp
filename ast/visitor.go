package ast

type Visitor interface {
	visit(e Acceptable)
}
