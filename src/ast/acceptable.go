package ast

type Acceptable interface {
	accept(v Visitor)
}
