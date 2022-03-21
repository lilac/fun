package ast

import (
	"github.com/lilac/fun-lang/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

var start = token.Position{Line: 1, Column: 1}
var end = token.Position{Line: 1, Column: 2}
var path = "<stdin>"
var tok = &token.Token{
	Kind:  0,
	Value: "",
	Location: token.Location{
		Start: start,
		End:   end,
		Path:  path,
	},
}

var startStr = path + ":" + start.String()
var endStr = path + ":" + end.String()

func TestUnit(t *testing.T) {
	unit := &Unit{HasToken{Token: tok}}
	assert.Equal(t, startStr, unit.Start().String())
	assert.Equal(t, endStr, unit.End().String(), "End token should match")
	assert.Equal(t, unit.String(), "()")
}

func TestBool(t *testing.T) {
	const boolValue = true
	node := &Bool{HasToken{tok}, boolValue}
	assert.Equal(t, startStr, node.Start().String())
	assert.Equal(t, endStr, node.End().String(), "End token should match")
	assert.Equal(t, node.String(), "true")
	assert.Equal(t, node.Value, boolValue)
}
