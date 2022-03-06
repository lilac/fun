package ast

import "github.com/lilac/funlang/types"

type Dec interface {
	Kind() string
}

type Arg struct {
	Id   Identifier
	Type types.Type
}

type ValDec struct {
	Vars []Var // type variables
	Arg  Arg
	Body Exp
}

type FunDec struct {
	Vars       []Var // type variables
	Id         Identifier
	Args       []Arg
	ResultType types.Type
	Body       Exp
}

type Module struct {
	Decs []Dec
}

func (v ValDec) Kind() string {
	return "val"
}

func (f FunDec) Kind() string {
	return "fun"
}
