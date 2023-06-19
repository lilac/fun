package compiler

import (
	"fmt"
	"github.com/lilac/fun-lang/pkg/alpha"
	"github.com/lilac/fun-lang/pkg/codegen"
	"github.com/lilac/fun-lang/pkg/syntax"
	"github.com/lilac/fun-lang/pkg/typing"
	"go/printer"
	"go/token"
	"os"
)

func Compile(source *syntax.Source) error {
	ti := typing.TypeInference{}
	module, err := syntax.Parse(source)
	if err != nil {
		return err
	}

	transformer := alpha.NewTransformer()
	transformer.Transform(module)

	env, err := ti.Infer(module)
	if err != nil {
		return err
	}
	dumpTypeEnv(env)

	irModule := lowerAst(module, env)
	fmt.Println(irModule)

	// code generation
	fun := codegen.GenFunction(nil)
	printer.Fprint(os.Stdout, token.NewFileSet(), fun)
	return nil
}

func dumpTypeEnv(env typing.TypeEnv) {
	for name, t := range env {
		typ := t.String()
		fmt.Printf("val %s : %s\n", name, typ)
	}
}
