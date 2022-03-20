package token

import "github.com/rhysd/locerr"

type Token struct {
	Kind  int
	value string
	// Line number.
	Line int
	// Column number.
	Column int
}

func NewToken(text string) *Token {
	return &Token{value: text}
}

func (t Token) Value() string {
	return t.value
}

func (t Token) Start() locerr.Pos {
	return locerr.Pos{Line: t.Line, Column: t.Column}
}

func (t Token) End() locerr.Pos {
	return locerr.Pos{Line: t.Line, Column: t.Column + len(t.value)}
}
