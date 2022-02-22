// Package token defines tokens of source codes.
package token

import (
	"fmt"
	"github.com/rhysd/locerr"
)

type Kind int

const (
	Illegal Kind = iota
	Comment
	Lparen
	Rparen
	Ident
	Bool
	Not
	Int
	Float
	Minus
	Plus
	Equal
	LessGreater
	LessEqual
	Less
	Greater
	GreaterEqual
	If
	Then
	Else
	Let
	In
	End
	Val
	Rec
	Comma
	Dot
	LessMinus
	Semicolon
	Star
	Slash
	BarBar
	AndAnd
	StringLiteral
	Percent
	Match
	With
	Bar
	MinusGreater
	Arrow
	Fn
	Fun
	Colon
	Type
	Lbracket
	Rbracket
	Eof
)

var tokenTable = [...]string{
	Illegal:       "Illegal",
	Eof:           "Eof",
	Comment:       "Comment",
	Lparen:        "(",
	Rparen:        ")",
	Ident:         "Ident",
	Bool:          "Bool",
	Not:           "Not",
	Int:           "Int",
	Float:         "Float",
	Minus:         "-",
	Plus:          "+",
	Equal:         "=",
	LessGreater:   "<>",
	LessEqual:     "<=",
	Less:          "<",
	Greater:       ">",
	GreaterEqual:  ">=",
	If:            "if",
	Then:          "then",
	Else:          "else",
	Let:           "let",
	In:            "in",
	End:           "end",
	Val:           "val",
	Rec:           "rec",
	Comma:         ",",
	Dot:           ".",
	LessMinus:     "<-",
	Semicolon:     ";",
	Star:          "*",
	Slash:         "/",
	BarBar:        "||",
	AndAnd:        "&&",
	StringLiteral: "STRING_LITERAL",
	Percent:       "%",
	Match:         "match",
	With:          "with",
	Bar:           "|",
	MinusGreater:  "->",
	Arrow:         "=>",
	Fn:            "fn",
	Fun:           "fun",
	Colon:         ":",
	Type:          "type",
	Lbracket:      "[",
	Rbracket:      "]",
}

// Token is the parsed elements of the source code.
// It contains its location information and kind.
type Token struct {
	Kind  Kind
	Start locerr.Pos
	End   locerr.Pos
	File  *locerr.Source
}

// String returns an information of token. This method is used mainly for
// debug purpose.
func (tok *Token) String() string {
	return fmt.Sprintf(
		"<%s:%s>(%d:%d:%d-%d:%d:%d)",
		tokenTable[tok.Kind],
		tok.Value(),
		tok.Start.Line, tok.Start.Column, tok.Start.Offset,
		tok.End.Line, tok.End.Column, tok.End.Offset)
}

// Value returns the corresponding a string part of code.
func (tok *Token) Value() string {
	return string(tok.File.Code[tok.Start.Offset:tok.End.Offset])
}
