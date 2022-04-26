// Package typing does type inference/reconstruction
package typing

import (
	"fmt"
	merror "github.com/hashicorp/go-multierror"
	"github.com/lilac/fun-lang/ast"
	"github.com/lilac/fun-lang/types"
)

type Env = map[string]types.Type

type VarSet = map[types.Var]bool

type TypeInference struct {
	nextVarId types.VarId
}

func (ti *TypeInference) generateVar() *types.Var {
	newVar := types.NewVar(ti.nextVarId)
	ti.nextVarId++
	return &newVar
}

func (ti *TypeInference) Infer(module *ast.Module) (Env, error) {
	var errors *merror.Error
	env := Env{}
	nonGenericVars := VarSet{}
	for _, dec := range module.Decs {
		err := ti.inferDec(env, nonGenericVars, dec)
		errors = merror.Append(errors, err)
	}
	return env, errors.ErrorOrNil()
}

func (ti *TypeInference) inferDec(env Env, nonGenericVars VarSet, dec ast.Dec) error {
	var errors error
	switch decl := dec.(type) {
	case *ast.ValDec:
		t, err := ti.inferExp(env, nonGenericVars, decl.Body)
		errors = merror.Append(errors, err)
		if annotatedType := decl.Arg.Type; annotatedType != nil {
			err := unify(annotatedType, t)
			errors = merror.Append(errors, err)
		}
		name := decl.Arg.Id.String()
		env[name] = t
	case *ast.FunDec:
	default:
		panic("unexpected ast.Dec type")
	}
	return errors
}

func (ti *TypeInference) inferExp(env Env, nonGenericVars VarSet, exp ast.Exp) (types.Type, error) {
	var errors error = nil
	switch node := exp.(type) {
	case ast.Constant:
		return node.Type(), nil
	case *ast.Var:
		return ti.typeOfVar(env, nonGenericVars, node.String())
	case *ast.Not:
		t, err := ti.inferExp(env, nonGenericVars, node.Child)
		errors = merror.Append(errors, err)
		err = unify(types.BoolType, t)
		errors = merror.Append(errors, err)
		return types.BoolType, errors
	case *ast.Neg:
		t, err := ti.inferExp(env, nonGenericVars, node.Child)
		errors = merror.Append(errors, err)
		if t != types.IntType && t != types.FloatType {
			err = fmt.Errorf("negation operator can only be applied to a number, but got %s", t)
			errors = merror.Append(errors, err)
		}
		return t, errors
	case *ast.InfixApp:
		at, err := ti.inferExp(env, nonGenericVars, node.Left)
		errors = merror.Append(errors, err)
		bt, err := ti.inferExp(env, nonGenericVars, node.Right)
		errors = merror.Append(errors, err)
		if !at.Equal(bt) {
			err = fmt.Errorf("mismatch type: %s != %s", at, bt)
			errors = merror.Append(errors, err)
		}
		switch node.Op.String() {
		case ast.Add, ast.Minus, ast.Mul, ast.Div, ast.Mod:
			if !at.Equal(types.IntType) && !at.Equal(types.FloatType) {
				err = fmt.Errorf("arithmethic operator can only be applied to a number, but got %s", at)
				errors = merror.Append(errors, err)
			}
			return at, errors
		case ast.Eq, ast.NotEq, ast.Less, ast.LessEq, ast.Greater, ast.GreaterEq:
			if !at.Equal(types.IntType) && !at.Equal(types.FloatType) {
				err = fmt.Errorf("arithmethic operator can only be applied to a number, but got %s", at)
				errors = merror.Append(errors, err)
			}
			return types.BoolType, errors
		case ast.And, ast.Or:
			if !at.Equal(types.BoolType) {
				err = fmt.Errorf("logical operator can only be applied to a boolean value, but got %s", at)
				errors = merror.Append(errors, err)
			}
			return types.BoolType, errors
		default:
			panic("Bug: unknown operator")
		}
	case *ast.Apply:
		resultType := ti.generateVar()
		argType, err := ti.inferExp(env, nonGenericVars, node.Arg)
		errors = merror.Append(errors, err)
		funType, err := ti.inferExp(env, nonGenericVars, node.Fun)
		errors = merror.Append(errors, err)
		expectedType := types.Arrow(argType, resultType)
		err = unify(funType, expectedType)
		errors = merror.Append(errors, err)
		return resultType, errors
	default:
		panic("Bug: unexpected expression type.")
	}
}

func (ti *TypeInference) typeOfVar(env Env, nonGenericVars VarSet, name string) (types.Type, error) {
	if t, ok := env[name]; ok {
		return ti.fresh(nonGenericVars, t), nil
	} else {
		err := fmt.Errorf("undefined symbol '%s'", name)
		return nil, err
	}
}

// fresh returns a type with all generic type variables substituted with new variables.
func (ti *TypeInference) fresh(nonGenericVars VarSet, t types.Type) types.Type {
	var varMap = map[*types.Var]*types.Var{}
	return ti.freshType(nonGenericVars, t, varMap)
}

func (ti *TypeInference) freshType(nonGenericVars VarSet, t types.Type, varMap map[*types.Var]*types.Var) types.Type {
	switch ty := prune(t).(type) {
	case *types.Var:
		if isGeneric(nonGenericVars, ty) {
			if v, ok := varMap[ty]; ok {
				return v
			} else {
				newVar := ti.generateVar()
				varMap[ty] = newVar
				return newVar
			}
		}
	case *types.CtorType:
		var newTypes []types.Type
		for _, arg := range ty.Args {
			freshType := ti.freshType(nonGenericVars, arg, varMap)
			newTypes = append(newTypes, freshType)
		}
		return &types.CtorType{
			Ctor: ty.Ctor,
			Args: newTypes,
		}
	case types.Var, types.CtorType:
		panic("Bug: a pointer to types.Type expected.")
	}
	return t
}

func unify(a, b types.Type) error {
	bt := prune(b)
	switch at := prune(a).(type) {
	case *types.Var:
		if at != bt {
			if occursInType(at, bt) {
				err := fmt.Errorf("recursive type unification")
				return err
			} else {
				at.Ref = bt
			}
		}
		// else ignore since they are equal.
	case *types.CtorType:
		switch bt := bt.(type) {
		case *types.Var:
			return unify(bt, at)
		case *types.CtorType:
			if at.Ctor != bt.Ctor || len(at.Args) != len(bt.Args) {
				err := fmt.Errorf("type mismatch: %s != %s", at.Ctor, bt.Ctor)
				return err
			} else if len(at.Args) > 0 {
				var errors error = nil
				for i, t := range at.Args {
					err := unify(t, bt.Args[i])
					errors = merror.Append(errors, err)
				}
				return errors
			} // else can be skipped since args are empty.
		default:
			panic("Bug: unexpected types.")
		}
	default:
		panic("Bug: unexpected types.")
	}
	return nil
}

// prune visits the type reference chain to get the ultimate type.
// As a side effect, all the type references are collapsed (flattened).
// todo: consider benefits of making it mutable.
func prune(t types.Type) types.Type {
	switch ty := t.(type) {
	case *types.Var:
		if ty.Ref != nil {
			otherType := prune(ty.Ref)
			ty.Ref = otherType
			return otherType
		}
	case types.Var, types.CtorType:
		panic("Bug: a pointer to types.Type expected.")
	}
	return t
}

func isGeneric(nonGenericVars VarSet, v *types.Var) bool {
	var ts []types.Type = nil
	for t := range nonGenericVars {
		ts = append(ts, t)
	}
	return !occursInTypes(v, ts)
}

// occursInType checks if the type variable occurs inside the other type.
//
// Note: the type var v must be pruned.
func occursInType(v *types.Var, t types.Type) bool {
	switch ty := prune(t).(type) {
	case *types.Var:
		// both vars are now pruned, so an equality check suffices.
		return v == ty
	case *types.CtorType:
		return occursInTypes(v, ty.Args)
	default:
		return false
	}
}

func occursInTypes(v *types.Var, ts []types.Type) bool {
	var result = false
	for _, t := range ts {
		result = result || occursInType(v, t)
	}
	return result
}