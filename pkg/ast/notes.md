## Notes
There are two kinds of nodes in the _abstract syntax tree_ (AST): declaration and expression. The root is the tree is always a `module`.

- All expression nodes implement the interface [Exp](expression.go).
- All declaration nodes implement the interface [Dec](declaration.go).

## References

[Standard ML grammar](https://people.mpi-sws.org/~rossberg/sml.html)