package typing

import (
	"github.com/lilac/fun-lang/alpha"
	"github.com/lilac/fun-lang/syntax"
	"github.com/lilac/fun-lang/types"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTypeInference_Infer(t *testing.T) {
	lines := []string{
		"val a = (1 + 2) * 3",
		"val b = not (a > 0) && false",
		"val s = \"abc\"",
		"val u = ()",
		"val f = -1.2",
		"val c = if f > 0. then a else 0",
	}
	env, _ := run(t, lines)
	assert.Equal(t, types.IntType, env["a$1"])
	assert.Equal(t, types.BoolType, env["b$2"])
	assert.Equal(t, types.StringType, env["s$3"])
	assert.Equal(t, types.UnitType, env["u$4"])
	assert.Equal(t, types.FloatType, env["f$5"])
	assert.Equal(t, types.IntType, env["c$6"])
}

func TestFnInference(t *testing.T) {
	lines := []string{
		"val id = fn x => x",
		"val a = id 3",
		"val s = id true",
		"val i = (fn f => fn x => f x) id 1.0",
	}
	env, _ := run(t, lines)
	//fmt.Println(env)
	aVar := types.NewVar(0)
	assert.Equal(t, types.Arrow(aVar, aVar).String(), env["id$2"].String())
	assert.Equal(t, types.IntType, env["a$3"].(*types.Var).Ref)
	assert.Equal(t, types.Arrow(types.FloatType, types.FloatType).String(), env["f$5"].String())
}

func run(t *testing.T, lines []string) (TypeEnv, error) {
	ti := TypeInference{}
	src := syntax.NewDummySource(strings.Join(lines, "\n"))
	module, err := syntax.Parse(src)
	assert.NoError(t, err, "parsing error")

	transformer := alpha.NewTransformer()
	transformer.Transform(module)
	assert.NoError(t, transformer.Error())

	env, err := ti.Infer(module)
	assert.NoError(t, err, "type inference error")

	return env, err
}
