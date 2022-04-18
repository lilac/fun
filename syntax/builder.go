package syntax

import (
	"fmt"
	"github.com/lilac/fun-lang/ast"
	"github.com/lilac/fun-lang/token"
	"github.com/lilac/fun-lang/types"
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

func NewInfixApp(left ast.Exp, tok *token.Token, right ast.Exp) *ast.InfixApp {
	return &ast.InfixApp{
		Left:  left,
		Op:    ast.Op{Name: tok.Value},
		Right: right,
	}
}

func NewTuple(left, right ast.Exp) *ast.Tuple {
	var result []ast.Exp
	switch left.(type) {
	case *ast.Tuple:
		result = append(left.(*ast.Tuple).Elements, right)
	default:
		result = []ast.Exp{left, right}
	}
	return &ast.Tuple{Elements: result}
}

func NewSequence(left, right ast.Exp) *ast.Sequence {
	var result []ast.Exp
	switch left.(type) {
	case *ast.Sequence:
		result = append(left.(*ast.Sequence).Elements, right)
	default:
		result = []ast.Exp{left, right}
	}
	return &ast.Sequence{Elements: result}
}

func NewIfThen(tok *token.Token, cond, then, els ast.Exp) *ast.IfThen {
	return &ast.IfThen{
		HasToken: ast.HasToken{Token: tok},
		Cond:     cond,
		Then:     then,
		Else:     els,
	}
}

func NewLet(tok *token.Token, dec []ast.Dec, exp ast.Exp) *ast.LetIn {
	return &ast.LetIn{
		HasToken: ast.HasToken{Token: tok},
		Decs:     dec,
		Body:     exp,
	}
}

func NewVarPattern(tok *token.Token) *ast.VarPattern {
	return &ast.VarPattern{HasToken: ast.HasToken{Token: tok}, Id: ast.Identifier{Name: tok.Value}}
}

func NewConstPattern(exp ast.Exp) *ast.ConstPattern {
	return &ast.ConstPattern{Constant: exp.(ast.Constant)}
}

func NewMatch(pattern ast.Pattern, exp ast.Exp) *ast.Match {
	return &ast.Match{
		Pattern: pattern,
		Exp:     exp,
	}
}

func NewFn(tok *token.Token, matches []ast.Match) *ast.Fn {
	return &ast.Fn{
		HasToken: ast.HasToken{Token: tok},
		Matches:  matches,
	}
}

func NewValDec(tok *token.Token, body ast.Exp) ast.Dec {
	return &ast.ValDec{
		Vars: []ast.Var{},
		Arg:  ast.Arg{Id: ast.Identifier{Name: tok.Value}, Type: nil},
		Body: body,
	}
}

func NewFunBind(tok *token.Token, patterns []ast.Pattern, ty types.Type, body ast.Exp) *ast.FunBind {
	return &ast.FunBind{
		HasToken:   ast.HasToken{Token: tok},
		Id:         ast.Identifier{Name: tok.Value},
		Patterns:   patterns,
		ResultType: ty,
		Exp:        body,
	}
}

func NewFunDec(binds []ast.FunBind) *ast.FunDec {
	return &ast.FunDec{
		Vars:  []ast.Var{},
		Binds: binds,
	}
}
