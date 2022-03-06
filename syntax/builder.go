package syntax

import (
	"fmt"
	"github.com/lilac/funlang/ast"
	"strconv"
)
import "github.com/lilac/funlang/token"

func NewUnit(tok *token.Token) *ast.Unit {
	node := ast.Unit{ast.HasToken{tok}}
	return &node
}

func NewBool(tok *token.Token) *ast.Bool {
	node := ast.Bool{ast.HasToken{tok}, tok.Value() == "true"}
	return &node
}

type ErrorFun = func(format string /*, args ...interface{}*/)

func NewInt(tok *token.Token, handler ErrorFun) *ast.Int {
	i, err := strconv.ParseInt(tok.Value(), 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Parse error at int literal: %s", err.Error())
		handler(msg)
		return &ast.Int{}
	} else {
		return &ast.Int{ast.HasToken{tok}, i}
	}
}

func NewFloat(tok *token.Token, handler ErrorFun) *ast.Float {
	i, err := strconv.ParseFloat(tok.Value(), 64)
	if err != nil {
		msg := fmt.Sprintf("Parse error at float literal: %s", err.Error())
		handler(msg)
		return &ast.Float{}
	} else {
		return &ast.Float{ast.HasToken{tok}, i}
	}
}

func NewVar(tok *token.Token) *ast.Var {
	return &ast.Var{ast.HasToken{tok}, ast.Identifier{Name: tok.Value()}}
}

func NewNot(tok *token.Token, child ast.Exp) *ast.Not {
	return &ast.Not{ast.HasToken{tok}, child}
}

func NewInfixApp(left ast.Exp, tok *token.Token, right ast.Exp) ast.Exp {
	return &ast.InfixApp{
		Left:  left,
		Op:    ast.Op{Name: tok.Value()},
		Right: right,
	}
}

func NewValDec(tok *token.Token, body ast.Exp) *ast.ValDec {
	return &ast.ValDec{
		Vars: []ast.Var{},
		Arg:  ast.Arg{Id: ast.Identifier{Name: tok.Value()}, Type: nil},
		Body: body,
	}
}