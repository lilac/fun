## Notes
The types are borrowed from the simple type system of the following type inference algorithm.

## Definitions
A type constructor `tycon` is made of an identifier. It's like the type level function.

A type `t` is one of
- type var
- concrete type of the form `tycon t_1 .. t_n`
  where n equals the arity of the type constructor
- arrow type `t -> t`
- tuple type `t * t`

Primitive types can be denoted by type constructors of 0-arity. In addition, arrow types can also be treated as concrete types of the special type constructor `->`.

## Examples
- Primitive types
  - `unit`
  - `bool`
  - `int`
  - `float`
  - `char`
  - `string`
- Type var
  - `'a * 'b`
- Tuples
  - `int * bool * string`
- Arrows
  - `'a -> int -> bool`

## References
- The HindleyMilner algorithm implemented in [Scala](http://dysphoria.net/code/hindley-milner/HindleyMilner.scala).
