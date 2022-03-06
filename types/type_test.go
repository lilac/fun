package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrimitiveTypes(t *testing.T) {
	assert.Equal(t, IntType.String(), "int")
	assert.Equal(t, UnitType.String(), "unit")
}

func TestArrow(t *testing.T) {
	ft := Arrow(IntType, FloatType)
	assert.Equal(t, "int -> float", ft.String())
}

func TestTupleType(t *testing.T) {
	tt := TupleType([]Type{IntType, StringType})
	assert.Equal(t, "int * string", tt.String())
}

func TestCtor(t *testing.T) {
	ct := CtorType{
		Ctor: "map",
		Args: []Type{IntType, StringType},
	}
	assert.Equal(t, "(int, string) map", ct.String())
}

func TestVar(t *testing.T) {
	v := NewVar(0)
	assert.Equal(t, "'a", v.String())
}
