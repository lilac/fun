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
		"val a = not (a > 0) && false",
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
	assert.Equal(t, types.BoolType, env["a$2"])
}
