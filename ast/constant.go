package ast

import (
	"fmt"
	"github.com/lilac/fun-lang/token"
	"github.com/lilac/fun-lang/types"
	"github.com/rhysd/locerr"
	"strconv"
)

// Constants
// con	::=
//     int	integer
//     word	word
//     float	floating point
//     char	character
//     string	string

type Constant interface {
	String() string
	Type() types.Type
}

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

func (h HasToken) Start() locerr.Pos {
	return h.Token.Start()
}

func (h HasToken) End() locerr.Pos {
	return h.Token.End()
}

func (u Unit) String() string {
	return "()"
}

func (u Unit) Type() types.Type {
	return types.UnitType
}

func (b Bool) String() string {
	return strconv.FormatBool(b.Value)
}

func (b Bool) Type() types.Type {
	return types.BoolType
}

func (i Int) String() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i Int) Type() types.Type {
	return types.IntType
}

func (f Float) String() string {
	return fmt.Sprintf("%f", f.Value)
}

func (f Float) Type() types.Type {
	return types.FloatType
}

func (s String) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

func (String) Type() types.Type {
	return types.StringType
}

func (c Char) String() string {
	return string(c.Value)
}

func (Char) Type() types.Type {
	return types.CharType
}
