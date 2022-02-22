package ast

import (
	"github.com/lilac/funlang/token"
	"github.com/rhysd/locerr"
)

// Constants
// con	::=
//     int	integer
//     word	word
//     float	floating point
//     char	character
//     string	string

type HasToken struct {
	Token *token.Token
}

type Unit struct {
	HasToken
}

type Bool struct {
	HasToken
	Value bool
}

type Int struct {
	HasToken
	Value int64
}

type Float struct {
	HasToken
	Value float64
}

type String struct {
	HasToken
	Value string
}

type Char struct {
	HasToken
	Value rune
}

func (h *HasToken) Start() locerr.Pos {
	return h.Token.Start
}

func (h *HasToken) End() locerr.Pos {
	return h.Token.End
}

func (u *Unit) Name() string {
	return "Unit"
}

func (b *Bool) Name() string {
	return "Bool"
}

func (i Int) Name() string {
	return "Int"
}

func (f Float) Name() string {
	return "Float"
}

func (s String) Name() string {
	return "String"
}

func (c Char) Name() string {
	return "Char"
}
