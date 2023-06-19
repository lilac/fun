package types

import (
	"fmt"
	"golang.org/x/exp/maps"
	"reflect"
	"strings"
	"unicode"
)

var (
	UnitType   = PrimitiveType("unit")
	BoolType   = PrimitiveType("bool")
	IntType    = PrimitiveType("int")
	FloatType  = PrimitiveType("float")
	StringType = PrimitiveType("string")
	CharType   = PrimitiveType("char")
)

type Type interface {
	String() string
	Equal(t Type) bool
	Prune() Type
	// VarSet returns the set of generic variables in this type
	VarSet() map[VarId]struct{}
}

type VarId int

// Var denotes a type variable
type Var struct {
	Id  VarId
	Ref Type
}

func (v Var) VarSet() map[VarId]struct{} {
	if v.Ref == nil {
		return map[VarId]struct{}{v.Id: {}}
	} else {
		return v.Ref.VarSet()
	}
}

// CtorType denotes a type derived from a type constructor
type CtorType struct {
	// Ctor is the type constructor name
	Ctor string
	Args []Type
}

func (c CtorType) VarSet() map[VarId]struct{} {
	var vars map[VarId]struct{}
	for _, arg := range c.Args {
		v := arg.VarSet()
		maps.Copy(vars, v)
	}
	return vars
}

func NewVar(id VarId) *Var {
	return &Var{Id: id, Ref: nil}
}

func PrimitiveType(name string) Type {
	return &CtorType{
		Ctor: name,
		Args: nil,
	}
}

func Arrow(a, b Type) Type {
	return &CtorType{
		Ctor: "->",
		Args: []Type{a, b},
	}
}

func TupleType(ts []Type) Type {
	return &CtorType{
		Ctor: "*",
		Args: ts,
	}
}

func (v Var) String() string {
	if v.Ref != nil {
		return v.Ref.String()
	}
	c := 'a' + int(v.Id)
	return "'" + string(rune(c))
}

func (v *Var) Equal(t Type) bool {
	switch vt := v.Prune().(type) {
	case *Var:
		switch ot := t.Prune().(type) {
		case *Var:
			return vt.Id == ot.Id || vt.Ref == ot.Ref || vt.Ref == ot || ot.Ref == vt
		default:
			return ot.Equal(vt)
		}
	default:
		return vt.Equal(t)
	}
}

func (c CtorType) String() string {
	count := len(c.Args)
	if count == 0 {
		return c.Ctor
	} else if count == 1 {
		return fmt.Sprintf("%s %s", c.Args[0].String(), c.Ctor)
	} else if count == 2 && !unicode.IsLetter(rune(c.Ctor[0])) {
		return fmt.Sprintf("%s %s %s", c.Args[0].String(), c.Ctor, c.Args[1].String())
	} else {
		args := make([]string, count)
		for i, v := range c.Args {
			args[i] = v.String()
		}
		s := strings.Join(args, ", ")
		return fmt.Sprintf("(%s) %s", s, c.Ctor)
	}
}

func (c CtorType) Equal(t Type) bool {
	switch ot := t.Prune().(type) {
	case *CtorType:
		return reflect.DeepEqual(c, *ot)
	default:
		return false
	}
}

// Prune visits the type reference chain to get the ultimate type.
// As a side effect, all the type references are collapsed (flattened).
// todo: consider benefits of making it mutable.
func (v *Var) Prune() Type {
	if v.Ref != nil {
		otherType := v.Ref.Prune()
		v.Ref = otherType
		return otherType
	}
	return v
}

func (c *CtorType) Prune() Type {
	return c
}
