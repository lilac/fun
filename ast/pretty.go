package ast

import "fmt"

// precedence returns the precedence of the root operation of the expression.
// Note: a lower number has higher precedence
func precedence(exp Exp) uint8 {
	switch exp.(type) {
	case Unit, *Unit, Bool, *Bool, Int, *Int, Float, *Float, String, *String, Char, *Char, Var, *Var, LetIn, *LetIn:
		// these expressions' starting and ending positions are clear, so they never need a parenthesis.
		return 1
	case Not, *Not, Neg, *Neg:
		return 2
	case Apply, *Apply:
		return 3
	case InfixApp:
		op := exp.(InfixApp).Op.String()
		return operatorPrecedence(op)
	case *InfixApp:
		op := exp.(*InfixApp).Op.String()
		return operatorPrecedence(op)
	case IfThen, *IfThen:
		return 8
	case Tuple, *Tuple, Sequence, *Sequence:
		return 10
	}
	return 0
}

func operatorPrecedence(op string) uint8 {
	switch op {
	case Mul, Div, Mod:
		return 4
	case Add, Minus:
		return 5
	case Eq, NotEq, Less, LessEq, Greater, GreaterEq:
		return 6
	case And, Or:
		return 7
	}
	panic(fmt.Sprintf("unknown operator %s", op))
}

func parenthesis(parent, child Exp) string {
	if precedence(parent) < precedence(child) {
		return fmt.Sprintf("(%v)", child)
	} else {
		return child.String()
	}
}

func parenthesisRight(parent, child Exp) string {
	if precedence(parent) <= precedence(child) {
		return fmt.Sprintf("(%v)", child)
	} else {
		return child.String()
	}
}
