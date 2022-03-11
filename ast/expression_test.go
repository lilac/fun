package ast

import (
	"github.com/lilac/funlang/token"
	"github.com/rhysd/locerr"
	"github.com/stretchr/testify/assert"
	"testing"
)

var s = locerr.NewDummySource("")
var start = locerr.Pos{Offset: 0, Line: 1, Column: 1, File: s}
var end = locerr.Pos{Offset: 1, Line: 1, Column: 2, File: s}
var tok = &token.Token{
	Kind:  token.Illegal,
	Start: start,
	End:   end,
	File:  s,
}

func TestUnit(t *testing.T) {
	unit := &Unit{HasToken{Token: tok}}
	assert.Equal(t, unit.Start(), start)
	assert.Equal(t, unit.End(), end, "End token should match")
	assert.Equal(t, unit.String(), "()")
}

func TestBool(t *testing.T) {
	const boolValue = true
	node := &Bool{HasToken{tok}, boolValue}
	assert.Equal(t, node.Start(), start)
	assert.Equal(t, node.End(), end, "End token should match")
	assert.Equal(t, node.String(), "true")
	assert.Equal(t, node.Value, boolValue)
}
