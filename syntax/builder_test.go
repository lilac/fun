package syntax

import (
	"fmt"
	"github.com/lilac/funlang/ast"
	"github.com/lilac/funlang/token"
	"github.com/rhysd/locerr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func newTestToken(kind token.Kind, str string) *token.Token {
	s := locerr.NewDummySource(str)
	start := locerr.Pos{
		Offset: 0,
		Line:   1,
		Column: 1,
		File:   s,
	}
	end := locerr.Pos{
		Offset: len(str),
		Line:   1,
		Column: len(str),
		File:   s,
	}
	return &token.Token{
		Kind:  kind,
		Start: start,
		End:   end,
		File:  s,
	}
}

func TestNewBool(t *testing.T) {
	tests := []string{
		"true",
		"false",
	}
	for _, tt := range tests {
		name := fmt.Sprintf("parse the bool %s", tt)
		tok := newTestToken(token.Bool, tt)
		want := &ast.Bool{
			ast.HasToken{tok},
			tt == "true",
		}
		t.Run(name, func(t *testing.T) {
			actual := NewBool(tok)
			assert.Equal(t, want, actual)
		})
	}
}

func TestNewInt(t *testing.T) {
	tests := []string{
		"122",
		"012",
		"-12",
	}
	handler := func(s string) {
		t.Error(s)
	}
	for _, tt := range tests {
		name := fmt.Sprintf("parse the int %s", tt)
		tok := newTestToken(token.Int, tt)
		num, err := strconv.ParseInt(tt, 10, 64)
		want := &ast.Int{
			ast.HasToken{tok},
			num,
		}
		if err != nil {
			t.Error(err)
		}
		t.Run(name, func(t *testing.T) {
			actual := NewInt(tok, handler)
			assert.Equal(t, want, actual)
		})
	}
}
