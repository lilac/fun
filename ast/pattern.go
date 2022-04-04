package ast

import (
	"fmt"
	"github.com/lilac/fun-lang/types"
	"strings"
)

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

type FunBind struct {
	Id         Identifier
	Patterns   []Pattern
	ResultType types.Type
	Exp        Exp
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

func (b FunBind) String() string {
	patterns := make([]string, len(b.Patterns))
	for i, pattern := range b.Patterns {
		patterns[i] = pattern.String()
	}
	pat := strings.Join(patterns, " ")
	if b.ResultType != nil {
		return fmt.Sprintf("%v %s : %v = %v", b.Id, pat, b.ResultType, b.Exp)
	}
	return fmt.Sprintf("%v %s = %v", b.Id, pat, b.Exp)
}
