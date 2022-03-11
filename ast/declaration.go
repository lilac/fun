package ast

import (
	"fmt"
	"github.com/lilac/fun-lang/types"
	"strings"
)

type Dec interface {
	Kind() string
	String() string
}

type Arg struct {
	Id   Identifier
	Type types.Type
}

type ValDec struct {
	Vars []Var // type variables
	Arg  Arg
	Body Exp
}

type FunDec struct {
	Vars       []Var // type variables
	Id         Identifier
	Args       []Arg
	ResultType types.Type
	Body       Exp
}

type Module struct {
	Decs []Dec
}

func (a Arg) String() string {
	if a.Type != nil {
		return fmt.Sprintf("%s : %s", a.Id.Name, a.Type.String())
	} else {
		return a.Id.Name
	}
}

func (v ValDec) Kind() string {
	return "val"
}

func (v ValDec) String() string {
	return fmt.Sprintf("val %s = %s", v.Arg, v.Body)
}

func (f FunDec) Kind() string {
	return "fun"
}

func (f FunDec) String() string {
	args := make([]string, len(f.Args))
	for i, arg := range f.Args {
		args[i] = "(" + arg.String() + ")"
	}
	argStr := strings.Join(args, " ")
	return fmt.Sprintf("fun %s %s = %s", f.Id, argStr, f.Body.String())
}
