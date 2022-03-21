package token

import (
	"fmt"
	"github.com/rhysd/locerr"
)

type Kind = int

type Position struct {
	// Line number.
	Line int
	// Column number.
	Column int
}

// Location denotes a range in a file
type Location struct {
	Start Position
	End   Position
	Path  string
}

type Token struct {
	Kind     Kind
	Value    string
	Location Location
}

func NewToken(text string) *Token {
	return &Token{Value: text}
}

func (t Token) Start() locerr.Pos {
	pos := t.Location.Start
	return locerr.Pos{Line: pos.Line, Column: pos.Column, File: &locerr.Source{Path: t.Location.Path}}
}

func (t Token) End() locerr.Pos {
	pos := t.Location.End
	return locerr.Pos{Line: pos.Line, Column: pos.Column, File: &locerr.Source{Path: t.Location.Path}}
}

func (t Token) String() string {
	return fmt.Sprintf("%v:%s", t.Kind, t.Value)
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}
