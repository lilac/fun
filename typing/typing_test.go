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
	ti := TypeInference{}

	lines := []string{
		"val a = (1 + 2) * 3",
		"val b = not (a > 0) && false",
		"val s = \"abc\"",
		"val u = ()",
		"val f = -1.2",
		"val c = if f > 0. then a else 0",
	}
	src := syntax.NewDummySource(strings.Join(lines, "\n"))
	module, err := syntax.Parse(src)
	assert.NoError(t, err, "parsing error")

	transformer := alpha.NewTransformer()
	transformer.Transform(module)
	assert.NoError(t, transformer.Error())

	env, err := ti.Infer(module)
	assert.NoError(t, err, "type inference error")
	assert.Equal(t, types.IntType, env["a$1"])
	assert.Equal(t, types.BoolType, env["b$2"])
	assert.Equal(t, types.StringType, env["s$3"])
	assert.Equal(t, types.UnitType, env["u$4"])
	assert.Equal(t, types.FloatType, env["f$5"])
	assert.Equal(t, types.IntType, env["c$6"])
}
