package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrimitiveTypes(t *testing.T) {
	assert.Equal(t, IntType.Name(), "int")
	assert.Equal(t, UnitType.Name(), "unit")
}

func TestArrow(t *testing.T) {
	ft := Arrow(IntType, FloatType)
	assert.Equal(t, "int -> float", ft.Name())
}

func TestTupleType(t *testing.T) {
	tt := TupleType([]Type{IntType, StringType})
	assert.Equal(t, "int * string", tt.Name())
}

func TestCtor(t *testing.T) {
	ct := CtorType{
		Ctor: "map",
		Args: []Type{IntType, StringType},
	}
	assert.Equal(t, "(int, string) map", ct.Name())
}

func TestVar(t *testing.T) {
	v := NewVar(0)
	assert.Equal(t, "'a", v.Name())
}
