package syntax

import (
	"fmt"
	"github.com/lilac/fun-lang/ast"
	"github.com/lilac/fun-lang/token"
	"strconv"
)

func NewUnit(tok *token.Token) *ast.Unit {
	return &ast.Unit{HasToken: ast.HasToken{Token: tok}}
}

func NewBool(tok *token.Token) *ast.Bool {
	return &ast.Bool{HasToken: ast.HasToken{Token: tok}, Value: tok.Value == "true"}
}

type ErrorFun = func(format string /*, args ...interface{}*/)

func NewInt(tok *token.Token, handler ErrorFun) *ast.Int {
	i, err := strconv.ParseInt(tok.Value, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Parse error at int literal: %s", err.Error())
		handler(msg)
		return &ast.Int{}
	} else {
		return &ast.Int{ast.HasToken{tok}, i}
	}
}

func NewFloat(tok *token.Token, handler ErrorFun) *ast.Float {
	i, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		msg := fmt.Sprintf("Parse error at float literal: %s", err.Error())
		handler(msg)
		return &ast.Float{}
	} else {
		return &ast.Float{ast.HasToken{tok}, i}
	}
}

func NewString(tok *token.Token, handler ErrorFun) *ast.String {
	s, err := strconv.Unquote(tok.Value)
	if err != nil {
		msg := fmt.Sprintf("Parse error at string literal %s: %v", tok.Value, err)
		handler(msg)
		return &ast.String{}
	}
	return &ast.String{ast.HasToken{tok}, s}
}

func NewVar(tok *token.Token) *ast.Var {
	return &ast.Var{ast.HasToken{tok}, ast.Identifier{Name: tok.Value}}
}

func NewNot(tok *token.Token, child ast.Exp) *ast.Not {
	return &ast.Not{HasToken: ast.HasToken{Token: tok}, Child: child}
}

func NewNeg(tok *token.Token, child ast.Exp) *ast.Neg {
	return &ast.Neg{HasToken: ast.HasToken{Token: tok}, Child: child}
}

func NewInfixApp(left ast.Exp, tok *token.Token, right ast.Exp) ast.Exp {
	return &ast.InfixApp{
		Left:  left,
		Op:    ast.Op{Name: tok.Value},
		Right: right,
	}
}

func NewValDec(tok *token.Token, body ast.Exp) ast.Dec {
	return &ast.ValDec{
		Vars: []ast.Var{},
		Arg:  ast.Arg{Id: ast.Identifier{Name: tok.Value}, Type: nil},
		Body: body,
	}
}
