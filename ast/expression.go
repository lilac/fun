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
	// Repr() string
}

type Apply struct {
	Fun Exp
	Arg Exp
}

type Not struct {
	HasToken
	Child Exp
}

type Neg struct {
	HasToken
	Child Exp
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

func (i IfThen) Start() locerr.Pos {
	return i.HasToken.Start()
}

func (i IfThen) End() locerr.Pos {
	return i.Else.End()
}

func (l LetIn) Name() string {
	return "Let"
}

func (l LetIn) Start() locerr.Pos {
	return l.HasToken.Start()
}

func (l LetIn) End() locerr.Pos {
	last := len(l.Body) - 1
	return l.Body[last].End()
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

func (n Not) Name() string {
	return "Not"
}

func (n Not) Start() locerr.Pos {
	return n.HasToken.Start()
}

func (n Not) End() locerr.Pos {
	return n.Child.End()
}

func (n Neg) Name() string {
	return "Neg"
}

func (n Neg) Start() locerr.Pos {
	return n.HasToken.Start()
}

func (n Neg) End() locerr.Pos {
	return n.Child.End()
}

func (n InfixApp) Name() string {
	return "InfixOp"
}

func (n InfixApp) Start() locerr.Pos {
	return n.Left.Start()
}

func (n InfixApp) End() locerr.Pos {
	return n.Right.End()
}
