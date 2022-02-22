package syntax

import (
	"fmt"
	"github.com/lilac/funlang/token"
	"github.com/rhysd/locerr"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

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
				s, err := locerr.NewSourceFromFile(n)
				if err != nil {
					panic(err)
				}
				l := NewLexer(s)
				go l.Lex()
				for {
					select {
					case tok := <-l.Tokens:
						switch tok.Kind {
						case token.Illegal:
							t.Fatal(tok.String())
						case token.Eof:
							return
						}
					}
				}
			})
		}
	}
}

// List literal can be lexed but parser should complain that it is not implemented yet.
// This behavior is implemented because array literal resembles to list literal.
func TestLexingListLiteral(t *testing.T) {
	s := locerr.NewDummySource("[1, 2, 3]")
	l := NewLexer(s)
	go l.Lex()
lexing:
	for {
		select {
		case tok := <-l.Tokens:
			switch tok.Kind {
			case token.Illegal:
				t.Fatal(tok.String())
			case token.Eof:
				break lexing
			}
		}
	}
}

func TestSampleProgram(t *testing.T) {
	program := `fun fib(n) = if n > 2 then fib(n-1) + fib(n-2) else 1
	val id = fn x => x
`
	s := locerr.NewDummySource(program)
	l := NewLexer(s)
	go l.Lex()
	for {
		select {
		case tok := <-l.Tokens:
			switch tok.Kind {
			case token.Illegal:
				t.Fatal(tok.String())
			case token.Eof:
				return
			}
		}
	}
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
			s, err := locerr.NewSourceFromFile(n)
			if err != nil {
				panic(err)
			}
			errorOccurred := false
			l := NewLexer(s)
			l.Error = func(_ string, _ locerr.Pos) {
				errorOccurred = true
			}
			go l.Lex()
			for {
				select {
				case tok := <-l.Tokens:
					switch tok.Kind {
					case token.Illegal:
						if !errorOccurred {
							t.Fatalf("Illegal token was emitted but no error occurred")
						}
						return
					case token.Eof:
						t.Fatalf("Lexing successfully done unexpectedly")
						return
					}
				}
			}
		})
	}
}
