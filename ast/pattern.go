package ast

import "fmt"

type Pattern interface {
	String() string
	IsPattern() bool
}

type ConstPattern struct {
	Constant Constant
}

type VarPattern struct {
	HasToken
	Id Identifier
}

type Match struct {
	Pattern Pattern
	Exp     Exp
}

func (c ConstPattern) String() string {
	return c.Constant.String()
}

func (c ConstPattern) IsPattern() bool {
	return true
}

func (v VarPattern) String() string {
	return v.Id.String()
}

func (v VarPattern) IsPattern() bool {
	return true
}

func (m Match) String() string {
	return fmt.Sprintf("%v => %v", m.Pattern, m.Exp)
}
