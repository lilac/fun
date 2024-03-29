## Features
[x] HindleyMilner algorithm

## Notes
The algorithm is a little tricky to understand the details. My tip is not expecting to grasp everything in one run.
You could try to understand the overall control flow of the algorithm, and then dig into the detail function implementation one by one.

## Algorithm
The type inference is a bottom-up algorithm.

### Unification
When two types are unified, we assert that type are equal in the type system.

### Type inference
The root expression is traversed from top to bottom, and the type of each sub-expression is inferred. A placeholder "type variable" is inserted when the type is unknown. In addition, type terms are unified in-place based on the typing rules.

The whole process is like solving a system of equations. At the end, some type variables are equal to concrete types, so they are called bound. Any remaining unbound type variables should be polymorphic.

### Generic variable
The definition of generic variables is:

_A type variable occurring in the type of an expression e is generic (with respect to e) iff it does not occur in the type of the binder of any fun expression enclosing e_.

### Data structures
- The type information of variables is represented by a map from identifiers (string) to `types.Type`. It's aliased as `TypeEnv`. Since type inference is executed after alpha transformation, all variables are uniquely identified, so a plain map suffices.
- The set of non-generic type variables is represented by a map from `types.Var` to `bool`, since Golang does not have set. Since this set is scope dependent, the common data structure `common.Env` is used, and aliased to `VarSet`. Upon entering a function scope, the type (variable) of each argument is bound (and not generic), so it is added to the set.
- A type variable `types.Var` includes an extra type reference. When the reference is not nil, it means the type variable equals to the referred type. In this way, type equations generated in the unification process becomes link chains of types.

### Functions
- `isGeneric(nonGenericVars common.Env[*types.Var, bool], v *types.Var) bool` returns if a type variable is generic or not. A type variable is generic if it does not occur in any one of the non-generic type variable set.
- `fresh(nonGenericVars VarSet, t types.Type) types.Type` returns a type with all generic type variables substituted with new variables.

## Example
For this [sample program](../../examples/id.fun), the final type information after type inference is as follows.
```sml
val id$1 : 'a -> 'a
val x$2 : 'a
val res$3 : int * bool
```

## References
- A tutorial explaining the [basic polynomial typechecking](http://lucacardelli.name/Papers/BasicTypechecking.pdf).
- The HindleyMilner algorithm implemented in [Scala](https://dysphoria.net/2009/06/28/hindley-milner-type-inference-in-scala/).
- A type inference implementation in [miniml](https://github.com/cmaes/miniml/blob/master/typing.ml#L74).
- [A tutorial on type inference of OCaml with step by step calculation](https://cs3110.github.io/textbook/chapters/interp/inference.html)
- [Hindley–Milner type system](https://en.wikipedia.org/wiki/Hindley%E2%80%93Milner_type_system)