package syntax

import (
	"fmt"
	"github.com/lilac/fun-lang/ast"
	"github.com/lilac/fun-lang/token"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func newTestToken(kind token.Kind, str string) *token.Token {
	start := token.Position{
		Line:   1,
		Column: 1,
	}
	end := token.Position{
		Line:   1,
		Column: len(str),
	}
	return &token.Token{
		Kind:  kind,
		Value: str,
		Location: token.Location{
			Start: start,
			End:   end,
			Path:  "<stdin>",
		},
	}
}

func TestNewBool(t *testing.T) {
	tests := []string{
		"true",
		"false",
	}
	for _, tt := range tests {
		name := fmt.Sprintf("parse the bool %s", tt)
		tok := newTestToken(Bool, tt)
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
		tok := newTestToken(Int, tt)
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
