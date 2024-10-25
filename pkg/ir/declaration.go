package ir

import (
	"github.com/lilac/fun-lang/pkg/ast"
	"github.com/lilac/fun-lang/pkg/types"
)

type decTag int8

const (
	valDecTag = iota
	funDecTag
)

type Dec interface {
	tag() decTag
}

type ValDec struct {
	//Vars []Var // type variables
	Id   ast.Identifier
	Type types.Type
	Body Exp
}

func (v ValDec) tag() decTag {
	return valDecTag
}

type FunDec struct {
	//Vars []Var // type variables
	Id   ast.Identifier
	Type types.Type
	Args []Arg
	Body Exp
}

func (f FunDec) tag() decTag {
	return funDecTag
}

type Module struct {
	Decs []Dec
}
