package alpha

import (
	"github.com/lilac/fun-lang/syntax"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func assertErrorContains(t *testing.T, err error, msg string) {
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), msg)
	}
}

func run(t *testing.T, lines []string) Transformer {
	src := syntax.NewDummySource(strings.Join(lines, "\n"))
	module, err := syntax.Parse(src)
	assert.NoError(t, err, "parsing error")

	transformer := Transformer{}
	transformer.Transform(module)
	return transformer
}

func TestUndefinedVariable(t *testing.T) {
	lines := []string{
		"val a = f 1",
	}
	transformer := run(t, lines)
	assertErrorContains(t, transformer.error, "Undefined variable 'f'")
}

func TestDuplicateId(t *testing.T) {
	lines := []string{
		"fun f x x = 1",
	}
	transformer := run(t, lines)
	assertErrorContains(t, transformer.error, "Duplicate identifier 'x' in pattern")
}

func TestConsistentFunName(t *testing.T) {
	lines := []string{
		"fun fib 1 = 1 | f x = x",
	}
	transformer := run(t, lines)

	assertErrorContains(t, transformer.error, "Function name is not consistent: f")
}

func TestConsistentFunArity(t *testing.T) {
	lines := []string{
		"fun fib 1 = 1 | fib x y = x",
	}
	transformer := run(t, lines)

	assertErrorContains(t, transformer.error, "Function arity is not consistent")
}

func TestUnderscoreVar(t *testing.T) {
	lines := []string{
		"fun f _ = _",
	}
	transformer := run(t, lines)

	assertErrorContains(t, transformer.error, "Cannot use '_' in variable reference")
}

func TestUnderscoreInPattern(t *testing.T) {
	lines := []string{
		"fun fib 1 = 1 | fib _ = 0",
		"val _ = 0",
	}
	src := syntax.NewDummySource(strings.Join(lines, "\n"))
	module, err := syntax.Parse(src)
	assert.NoError(t, err, "parsing error")
	expectedLines := []string{
		"fun fib$1 1 = 1 | fib _$2 = 0",
		"val _$3 = 0",
	}
	transformer := Transformer{}
	transformer.Transform(module)
	assert.NoError(t, transformer.error)
	assert.Equal(t, strings.Join(expectedLines, "\n"), module.String())
}

func TestTransform(t *testing.T) {
	lines := []string{
		"val a = (1 + 2) * 3",
		"val a = let val x = 1 val y = \"ab\" in x > 0; 1 end",
		"val a = if 1 > 0 then 1 else fn true => 1 | x => 0",
		"fun fib 0 = 0 | fib 1 = 1 | fib x = fib (x - 1) + fib (x - 2)",
	}
	src := syntax.NewDummySource(strings.Join(lines, "\n"))
	module, err := syntax.Parse(src)
	assert.NoError(t, err, "parsing error")
	transformer := Transformer{time: 0}
	transformer.Transform(module)
	s := module.String()
	expectedLines := []string{
		"val a$1 = (1 + 2) * 3",
		"val a$4 = let val x$2 = 1 val y$3 = \"ab\" in x$2 > 0; 1 end",
		"val a$6 = if 1 > 0 then 1 else fn true => 1 | x$5 => 0",
		"fun fib$7 0 = 0 | fib 1 = 1 | fib x$8 = fib$7 (x$8 - 1) + fib$7 (x$8 - 2)",
	}
	assert.NoError(t, transformer.error)
	assert.Equal(t, strings.Join(expectedLines, "\n"), s)
}
