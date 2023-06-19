// Package ir contains the types for intermediate representation.

package ir

import (
	"github.com/lilac/fun-lang/pkg/ast"
	"github.com/lilac/fun-lang/pkg/types"
)

type expTag int

const (
	// constants
	unitTag = iota
	boolTag
	charTag
	intTag
	floatTag
	stringTag
	// unary operation
	notTag
	negTag
	// infix operation
	binaryOpTag

	tupleTag
	SequenceTag

	varTag
	appTag
	ifThenTag
	letInTag
	funTag
)

type Exp interface {
	//Start() locerr.Pos
	//End() locerr.Pos
	tag() expTag
}

type Tuple struct {
	Elements []Exp
}

func (t Tuple) tag() expTag {
	return tupleTag
}

type Sequence struct {
	Elements []Exp
}

func (s Sequence) tag() expTag {
	return SequenceTag
}

type Var struct {
	Id ast.Identifier
}

func (v Var) tag() expTag {
	return varTag
}

type App struct {
	Id  ast.Identifier
	Arg Exp
}

func (a App) tag() expTag {
	return appTag
}

type IfThen struct {
	Cond Exp
	Then Exp
	Else Exp
}

func (i IfThen) tag() expTag {
	return ifThenTag
}

type LetIn struct {
	Decs []Dec
	Body Exp
}

func (l LetIn) tag() expTag {
	return letInTag
}

type Arg struct {
	Id   ast.Identifier
	Type types.Type
}

type Fn struct {
	Id   ast.Identifier
	Type types.Type
	Body Exp
}

func (f Fn) tag() expTag {
	return funTag
}
