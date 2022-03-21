package syntax

import (
	"fmt"
	"github.com/lilac/fun-lang/token"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func checkTokens(t *testing.T, tokens []*token.Token) {
	for _, tok := range tokens {
		switch tok.Kind {
		case Illegal:
			t.Fatalf("Illegal token %v\n", tok)
		case Eof:
			return
		default:
			fmt.Printf("token %v\n", tok)
		}
	}
}

func TestLexingOK(t *testing.T) {
	for _, testDir := range []string{
		"test-data",
	} {
		files, err := ioutil.ReadDir(filepath.FromSlash(testDir))
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			n := filepath.Join(testDir, f.Name())
			if !strings.HasSuffix(n, ".ml") {
				continue
			}

			t.Run(fmt.Sprintf("Check lexing successfully: %s", n), func(t *testing.T) {
				s, err := NewSourceFromFile(n)
				if err != nil {
					panic(err)
				}
				l := NewLexer(s)
				checkTokens(t, l.LexAll())
			})
		}
	}
}

// List literal can be lexed but parser should complain that it is not implemented yet.
// This behavior is implemented because array literal resembles to list literal.
func TestLexingListLiteral(t *testing.T) {
	s := NewDummySource("[1, 2, 3]")
	l := NewLexer(s)
	tokens := l.LexAll()
	checkTokens(t, tokens)
}

func TestSampleProgram(t *testing.T) {
	program := `fun fib(n) = if n > 2 then fib(n-1) + fib(n-2) else 1
	val id = fn x => x
`
	s := NewDummySource(program)
	l := NewLexer(s)
	checkTokens(t, l.LexAll())
}

func TestLexingIllegal(t *testing.T) {
	testDir := filepath.FromSlash("test-data/invalid")
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		n := filepath.Join(testDir, f.Name())
		if !strings.HasSuffix(n, ".ml") {
			continue
		}

		t.Run(fmt.Sprintf("Check lexing illegal input: %s", f.Name()), func(t *testing.T) {
			s, err := NewSourceFromFile(n)
			if err != nil {
				panic(err)
			}
			errorOccurred := false
			l := NewLexer(s)

			for _, tok := range l.LexAll() {
				switch tok.Kind {
				case Illegal:
					errorOccurred = true
				case Eof:
					if !errorOccurred {
						t.Fatalf("Lexing successfully done unexpectedly")
					}
					return
				}
			}
		})
	}
}
