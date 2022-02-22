package ast

import (
	"fmt"
	"github.com/lilac/funlang/token"
	"github.com/lilac/funlang/types"
	"github.com/rhysd/locerr"
)

type Exp interface {
	Start() locerr.Pos
	End() locerr.Pos
	Name() string
}

type Apply struct {
	Fun Exp
	Arg Exp
}

type Op Identifier

type InfixApp struct {
	Left  Exp
	Op    Op
	Right Exp
}

type Tuple struct {
	Elements []Exp
}

type Sequence struct {
	Elements []Exp
}

type IfThen struct {
	HasToken
	Cond Exp
	Then Exp
	Else Exp
}

type LetIn struct {
	HasToken
	Decs []Dec
	Body []Exp
}

type TypeAnnotation struct {
	Exp      Exp
	Type     types.Type
	EndToken token.Token
}

func (a *Apply) Start() locerr.Pos {
	return a.Fun.Start()
}

func (a *Apply) End() locerr.Pos {
	return a.Arg.End()
}

func (*Apply) Name() string {
	return "Apply"
}

func (t Tuple) Start() locerr.Pos {
	return t.Elements[0].Start()
}

func (t Tuple) End() locerr.Pos {
	size := len(t.Elements)
	return t.Elements[size-1].End()
}

func (t Tuple) Name() string {
	return "Tuple"
}

func (s Sequence) Start() locerr.Pos {
	return s.Elements[0].Start()
}

func (s Sequence) End() locerr.Pos {
	size := len(s.Elements)
	return s.Elements[size-1].End()
}

func (s Sequence) Name() string {
	return "Sequence"
}

func (i IfThen) Name() string {
	return "If"
}

func (l LetIn) Name() string {
	return "Let"
}

func (t TypeAnnotation) Start() locerr.Pos {
	return t.Exp.Start()
}

func (t TypeAnnotation) End() locerr.Pos {
	return t.EndToken.End
}

func (t TypeAnnotation) Name() string {
	return fmt.Sprintf("%s : %s", t.Exp, t.Type)
}
