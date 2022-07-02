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

## Example
For this [sample program](../../examples/id.fun), the final type information after type inference is as follows.
```sml
val id$1 : 'a -> 'a
val x$2 : 'a
val res$3 : int * bool
```

## References
- A tutorial explaining the [basic polynomial typechecking](http://lucacardelli.name/Papers/BasicTypechecking.pdf).
- The HindleyMilner algorithm implemented in [Scala](http://dysphoria.net/code/hindley-milner/HindleyMilner.scala).
- A type inference implementation in [miniml](https://github.com/cmaes/miniml/blob/master/typing.ml#L74).
- [Hindley–Milner type system](https://en.wikipedia.org/wiki/Hindley%E2%80%93Milner_type_system)