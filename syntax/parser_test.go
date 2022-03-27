package syntax

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTooLargeIntLiteral(t *testing.T) {
	src := NewDummySource("val a = 123456789123456789123456789123456789123456789")

	r, err := Parse(src)
	if err == nil {
		t.Fatalf("Invalid int literal must raise an error but got %v", r)
	}
	if !strings.Contains(err.Error(), "value out of range") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestInvalidStringLiteral(t *testing.T) {
	src := NewDummySource("val s = \"a\nb\"\n")
	r, err := Parse(src)
	if err == nil {
		t.Fatalf("Invalid string literal must raise an error but got %v", r)
	}
}

func TestTooLargeFloatLiteral(t *testing.T) {
	src := NewDummySource("val f = 1.7976931348623159e308")

	r, err := Parse(src)
	if err == nil {
		t.Fatalf("Invalid int literal must raise an error but got %v", r)
	}
	if !strings.Contains(err.Error(), "value out of range") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestLexFailed(t *testing.T) {
	src := NewDummySource("(* comment is not closed")
	_, err := Parse(src)
	if err == nil {
		t.Fatal("Lex error was not reported")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Expected '*' for closing comment") {
		t.Fatal("Unexpected error message:", msg)
	}
}

func TestParseDecs(t *testing.T) {
	lines := []string{
		"val a = 1",
		"val b = true",
		"val t = not false",
		"val u = ()",
		"val s = \"abc\"",
		"val x = s",
	}
	src := NewDummySource(strings.Join(lines, "\n"))
	module, err := Parse(src)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	actual := make([]string, len(module.Decs))
	for i, d := range module.Decs {
		actual[i] = d.String()
	}
	assert.Equal(t, lines, actual)
}

func TestPrecedence(t *testing.T) {
	lines := []string{
		"val a = (1 + 2) * 3",
		"val a = (1 + 2) * 3 - 1 > 0 && false",
		"val a = 1 + 3 * 4 / 2",
		"val a = 1 + 3 * (4 / 2)",
		"val a = not false",
		"val a = not (3 > 0)",
		"val a = -3",
		"val a = -(3 - 1)",
		//"val a = if true then 1 else 0 + 1",
	}
	src := NewDummySource(strings.Join(lines, "\n"))
	module, err := Parse(src)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	for i, d := range module.Decs {
		assert.Equal(t, lines[i], d.String())
	}
}
