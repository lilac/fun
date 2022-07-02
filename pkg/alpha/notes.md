## Notes
This package includes a transformer that renames all identifiers of a program, and also does some semantic checking.

A core data structure used is [Env](../common/env.go), representing the context scope of expressions. `Env` is a generic type. For alpha transformer, the context data is a mapping from the original variable names to unique names; so we define `NameEnv` as `Env[string, string]`.

During the traversal of the abstract semantic tree, we create a new `NameEnv` when a new scope is entered.

## References
- The [alpha](https://github.com/esumii/min-caml/blob/master/alpha.ml#L7) module of min-caml.
- The [alpha_transform](https://github.com/rhysd/gocaml/blob/master/sema/alpha_transform.go#L41) package of gocaml.
- The [ASTChangeNames](https://github.com/j-c-w/mlc/blob/master/src/main/scala/ast_change_names/ASTChangeNames.scala#L10) object of mlc compiler.