package ast

import (
	"fmt"
	"strings"
)

type Dec interface {
	Kind() string
	String() string
}

type ValDec struct {
	Vars []Var // type variables
	Arg  Arg
	Body Exp
}

type FunDec struct {
	Vars  []Var // type variables
	Binds []FunBind
}

type Module struct {
	Decs []Dec
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
	binds := make([]string, len(f.Binds))
	for i, bind := range f.Binds {
		binds[i] = bind.String()
	}
	s := strings.Join(binds, " | ")
	return fmt.Sprintf("fun %s", s)
}
