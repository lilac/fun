package types

import (
	"fmt"
	"strings"
	"unicode"
)

var (
	UnitType   = PrimitiveType("unit")
	IntType    = PrimitiveType("int")
	FloatType  = PrimitiveType("float")
	StringType = PrimitiveType("string")
)

type Type interface {
	Name() string
}

type VarId int

// Var denotes a type variable
type Var struct {
	Id  VarId
	Ref Type
}

// CtorType denotes a type derived from a type constructor
type CtorType struct {
	// Ctor is the type constructor name
	Ctor string
	Args []Type
}

func NewVar(id VarId) Var {
	return Var{Id: id, Ref: nil}
}

func PrimitiveType(name string) Type {
	return CtorType{
		Ctor: name,
		Args: nil,
	}
}

func Arrow(a, b Type) Type {
	return CtorType{
		Ctor: "->",
		Args: []Type{a, b},
	}
}

func TupleType(ts []Type) Type {
	return CtorType{
		Ctor: "*",
		Args: ts,
	}
}

func (v Var) Name() string {
	c := 'a' + int(v.Id)
	return "'" + string(rune(c))
}

func (c CtorType) Name() string {
	count := len(c.Args)
	if count == 0 {
		return c.Ctor
	} else if count == 1 {
		return fmt.Sprintf("%s %s", c.Args[0].Name(), c.Ctor)
	} else if count == 2 && !unicode.IsLetter(rune(c.Ctor[0])) {
		return fmt.Sprintf("%s %s %s", c.Args[0].Name(), c.Ctor, c.Args[1].Name())
	} else {
		args := make([]string, count)
		for i, v := range c.Args {
			args[i] = v.Name()
		}
		s := strings.Join(args, ", ")
		return fmt.Sprintf("(%s) %s", s, c.Ctor)
	}
}
