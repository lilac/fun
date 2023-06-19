package compiler

import (
	"github.com/lilac/fun-lang/pkg/ast"
	"github.com/lilac/fun-lang/pkg/ir"
	"github.com/lilac/fun-lang/pkg/typing"
)

func lowerAst(module *ast.Module, env typing.TypeEnv) *ir.Module {
	return &ir.Module{Decs: nil}
}
