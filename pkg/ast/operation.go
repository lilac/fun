package ast

type Op = Identifier

const (
	Add   = "+"
	Minus = "-"
	Mul   = "*"
	Div   = "/"
	Mod   = "%"

	Eq        = "="
	NotEq     = "<>"
	Less      = "<"
	LessEq    = "<="
	Greater   = ">"
	GreaterEq = ">="

	And = "&&"
	Or  = "||"
)
