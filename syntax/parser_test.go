package syntax

import (
	"github.com/rhysd/locerr"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTooLargeIntLiteral(t *testing.T) {
	src := locerr.NewDummySource("val a = 123456789123456789123456789123456789123456789")

	r, err := Parse(src)
	if err == nil {
		t.Fatalf("Invalid int literal must raise an error but got %v", r)
	}
	if !strings.Contains(err.Error(), "value out of range") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestInvalidStringLiteral(t *testing.T) {
	src := locerr.NewDummySource("val s = \"a\nb\"\n")
	r, err := Parse(src)
	if err == nil {
		t.Fatalf("Invalid string literal must raise an error but got %v", r)
	}
}

func TestTooLargeFloatLiteral(t *testing.T) {
	src := locerr.NewDummySource("val f = 1.7976931348623159e308")

	r, err := Parse(src)
	if err == nil {
		t.Fatalf("Invalid int literal must raise an error but got %v", r)
	}
	if !strings.Contains(err.Error(), "value out of range") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestLexFailed(t *testing.T) {
	src := locerr.NewDummySource("(* comment is not closed")
	_, err := Parse(src)
	if err == nil {
		t.Fatal("Lex error was not reported")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Lexing source into tokens failed") {
		t.Fatal("Unexpected error message:", msg)
	}
}

func TestParseDecs(t *testing.T) {
	lines := []string{
		"val a = 1",
		//"val b = true",
		//"val t = not false",
		//"val u = ()",
		//"val s = \"abc\"",
		//"val x = s",
	}
	src := locerr.NewDummySource(strings.Join(lines, "\n"))
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
