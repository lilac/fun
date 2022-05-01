package alpha

import (
	"fmt"
	merror "github.com/hashicorp/go-multierror"
	"github.com/lilac/fun-lang/ast"
	. "github.com/lilac/fun-lang/common"
	"github.com/lilac/fun-lang/syntax"
	"github.com/rhysd/locerr"
)

// rename identifiers to make them unique (alpha conversion)

type Transformer struct {
	time  uint // a monotonically increasing number to make names unique
	error error
}

type NameEnv = Env[string, string]

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t Transformer) Error() error {
	return t.error
}

func (t *Transformer) transformExp(env *NameEnv, exp ast.Exp) ast.Exp {
	// todo: make the transformation in-place.
	switch node := exp.(type) {
	case ast.Constant:
		return exp
	//case ast.Var:
	//	return t.transformVar(env, &node)
	case *ast.Var:
		return t.transformVar(env, node)
	case *ast.Fn:
		for i, m := range node.Matches {
			var matchEnv = NewEnv(env)
			_ = t.transformPattern(matchEnv, m.Pattern)
			e := t.transformExp(matchEnv, m.Exp)
			node.Matches[i].Exp = e
		}
		return node
	case *ast.LetIn:
		letEnv := NewEnv(env)
		for i, d := range node.Decs {
			node.Decs[i] = t.transformDec(letEnv, d)
		}
		node.Body = t.transformExp(letEnv, node.Body)
		return node
	// Other cases are just recursive top-down transformations
	case *ast.Not:
		e := t.transformExp(env, node.Child)
		node.Child = e
		return node
	case *ast.Neg:
		e := t.transformExp(env, node.Child)
		node.Child = e
		return node
	case *ast.Tuple:
		for i, e := range node.Elements {
			x := t.transformExp(env, e)
			node.Elements[i] = x
		}
		return node
	case *ast.Sequence:
		for i, e := range node.Elements {
			x := t.transformExp(env, e)
			node.Elements[i] = x
		}
		return node
	case *ast.Apply:
		fun := t.transformExp(env, node.Fun)
		arg := t.transformExp(env, node.Arg)
		return &ast.Apply{
			Fun: fun,
			Arg: arg,
		}
	case *ast.InfixApp:
		left := t.transformExp(env, node.Left)
		right := t.transformExp(env, node.Right)
		return &ast.InfixApp{
			Left:  left,
			Op:    node.Op,
			Right: right,
		}
	case *ast.IfThen:
		cond := t.transformExp(env, node.Cond)
		then := t.transformExp(env, node.Then)
		els := t.transformExp(env, node.Else)
		return syntax.NewIfThen(node.Token, cond, then, els)
	case ast.Var, ast.Apply, ast.Fn, ast.IfThen, ast.InfixApp:
		panic("bug: an Exp value is supplied to transformExp")
	default:
		return exp
	}
}

func (t *Transformer) transformVar(env *Env[string, string], v *ast.Var) ast.Exp {
	if v.Id.Name == "_" {
		t.errorfIn(v, "Cannot use '_' in variable reference")
		return v
	}
	name, ok := env.LookUp(v.Id.Name)
	if ok {
		v.Id.Value = name
	} else {
		t.errorfIn(v, "Undefined variable '%s'", v.Id.Name)
	}
	return v
}

func (t *Transformer) errorfIn(exp ast.Exp, format string, args ...interface{}) {
	e := locerr.ErrorfIn(exp.Start(), exp.End(), format, args...)
	t.error = merror.Append(t.error, e)
}

func (t *Transformer) transformPattern(env *Env[string, string], pattern ast.Pattern) ast.Pattern {
	switch node := pattern.(type) {
	case *ast.VarPattern:
		// check duplicate id in the same pattern list.
		if env.Contain(node.Id.Name) {
			t.errorfIn(pattern, "Duplicate identifier '%s' in pattern", node.Id.Name)
		} else {
			t.bind(env, &node.Id)
		}
	}
	return pattern
}

func (t *Transformer) newUniqueId(name string) string {
	t.time++
	uniqueName := fmt.Sprintf("%s$%d", name, t.time)
	return uniqueName
}

func (t *Transformer) bind(env *Env[string, string], id *ast.Identifier) {
	uid := t.newUniqueId(id.Name)
	id.Value = uid
	env.Add(id.Name, uid)
}

func (t *Transformer) transformDec(env *NameEnv, dec ast.Dec) ast.Dec {
	switch node := dec.(type) {
	case *ast.ValDec:
		e := t.transformExp(env, node.Body)
		node.Body = e
		t.bind(env, &node.Arg.Id)
	case *ast.FunDec:
		id := &node.Binds[0].Id
		arity := len(node.Binds[0].Patterns)
		/*
			We don't need the check since the grammar has dictated that.
			if arity == 0 {
				t.errorfIn(node.Binds[0], "A function should have at least one argument: %s", id.Name)
			}
		*/
		t.bind(env, id)
		for _, bind := range node.Binds {
			if bind.Id.Name != id.Name {
				t.errorfIn(bind, "Function name is not consistent: %s", bind.Id.Name)
			}
			if len(bind.Patterns) != arity {
				t.errorfIn(bind, "Function arity is not consistent: the arity of \"%s\" is %d", id.Name, arity)
			}
			bind.Id.Value = id.Value
			// a new environment for each bind
			bindEnv := NewEnv(env)
			for _, pattern := range bind.Patterns {
				_ = t.transformPattern(bindEnv, pattern)
			}
			e := t.transformExp(bindEnv, bind.Exp)
			bind.Exp = e
		}
	}
	return dec
}

func (t *Transformer) Transform(module *ast.Module) {
	env := NewEnv[string, string](nil)
	for i, dec := range module.Decs {
		module.Decs[i] = t.transformDec(env, dec)
	}
}
