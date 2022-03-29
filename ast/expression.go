package ast

import (
	"fmt"
	"github.com/lilac/fun-lang/token"
	"github.com/lilac/fun-lang/types"
	"github.com/rhysd/locerr"
	"strings"
)

type Exp interface {
	Start() locerr.Pos
	End() locerr.Pos
	String() string
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
	Body Exp
}

type Arg struct {
	Id   Identifier
	Type types.Type
}

type Fn struct {
	HasToken
	Matches []Match
}

type TypeAnnotation struct {
	Exp      Exp
	Type     types.Type
	EndToken token.Token
}

func (a Apply) Start() locerr.Pos {
	return a.Fun.Start()
}

func (a Apply) End() locerr.Pos {
	return a.Arg.End()
}

func (a Apply) String() string {
	return fmt.Sprintf("%s %s", parenthesis(a, a.Fun), parenthesis(a, a.Arg))
}

func (t Tuple) Start() locerr.Pos {
	return t.Elements[0].Start()
}

func (t Tuple) End() locerr.Pos {
	size := len(t.Elements)
	return t.Elements[size-1].End()
}

func (t Tuple) String() string {
	elements := make([]string, len(t.Elements))
	for i, t := range t.Elements {
		elements[i] = t.String()
	}
	return strings.Join(elements, ", ")
}

func (s Sequence) Start() locerr.Pos {
	return s.Elements[0].Start()
}

func (s Sequence) End() locerr.Pos {
	size := len(s.Elements)
	return s.Elements[size-1].End()
}

func (s Sequence) String() string {
	elements := make([]string, len(s.Elements))
	for i, t := range s.Elements {
		elements[i] = parenthesis(s, t)
	}
	return strings.Join(elements, "; ")
}

func (i IfThen) String() string {
	return fmt.Sprintf("if %s then %s else %s", i.Cond, i.Then, i.Else)
}

func (i IfThen) Start() locerr.Pos {
	return i.HasToken.Start()
}

func (i IfThen) End() locerr.Pos {
	return i.Else.End()
}

func (l LetIn) String() string {
	elements := make([]string, len(l.Decs))
	for i, dec := range l.Decs {
		elements[i] = dec.String()
	}
	decs := strings.Join(elements, " ")
	return fmt.Sprintf("let %s in %v end", decs, l.Body)
}

func (l LetIn) Start() locerr.Pos {
	return l.HasToken.Start()
}

func (l LetIn) End() locerr.Pos {
	return l.Body.End()
}

func (t TypeAnnotation) Start() locerr.Pos {
	return t.Exp.Start()
}

func (t TypeAnnotation) End() locerr.Pos {
	return t.EndToken.End()
}

func (t TypeAnnotation) String() string {
	if t.Type != nil {
		return fmt.Sprintf("%s : %s", t.Exp, t.Type)
	} else {
		return t.Exp.String()
	}
}

func (n Not) String() string {
	return fmt.Sprintf("not %s", parenthesis(n, n.Child))
}

func (n Not) Start() locerr.Pos {
	return n.HasToken.Start()
}

func (n Not) End() locerr.Pos {
	return n.Child.End()
}

func (n Neg) String() string {
	return fmt.Sprintf("-%s", parenthesis(n, n.Child))
}

func (n Neg) Start() locerr.Pos {
	return n.HasToken.Start()
}

func (n Neg) End() locerr.Pos {
	return n.Child.End()
}

func (n InfixApp) String() string {
	return fmt.Sprintf("%s %s %s", parenthesis(n, n.Left), n.Op, parenthesis(n, n.Right))
}

func (n InfixApp) Start() locerr.Pos {
	return n.Left.Start()
}

func (n InfixApp) End() locerr.Pos {
	return n.Right.End()
}

func (a Arg) String() string {
	if a.Type != nil {
		return fmt.Sprintf("(%v: %v)", a.Id, a.Type)
	}
	return a.Id.String()
}

func (f Fn) End() locerr.Pos {
	l := len(f.Matches)
	return f.Matches[l-1].Exp.End()
}

func (f Fn) String() string {
	elements := make([]string, len(f.Matches))
	for i, m := range f.Matches {
		elements[i] = m.String()
	}
	s := strings.Join(elements, " | ")
	return fmt.Sprintf("fn %s", s)
}
