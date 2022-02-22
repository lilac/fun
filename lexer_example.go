package main

import (
	"flag"
	"fmt"
	"github.com/lilac/funlang/syntax"
	"github.com/lilac/funlang/token"
	"github.com/rhysd/locerr"
	"os"
	"path/filepath"
)

const usageDoc = `Usage: go run lexer_example [src]`

func usage() {
	fmt.Fprintln(os.Stderr, usageDoc)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var src *locerr.Source
	var err error

	if flag.NArg() == 0 {
		src, err = locerr.NewSourceFromStdin()
	} else {
		file := filepath.FromSlash(flag.Arg(0))
		src, err = locerr.NewSourceFromFile(file)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on opening source: %s\n", err.Error())
		os.Exit(4)
	}
	LexFile(src)
}

func LexFile(src *locerr.Source) {
	lex := syntax.NewLexer(src)

	// Start to lex the source in other goroutine
	go lex.Lex()

	// tokens will be sent from lex.Tokens channel
	for {
		select {
		case tok := <-lex.Tokens:
			switch tok.Kind {
			case token.Illegal:
				fmt.Printf("Lexing invalid token at %v\n", tok.Start)
				return
			case token.Eof:
				fmt.Println("End of input")
				return
			default:
				fmt.Printf("Token: %s", tok.String())
			}
		}
	}
}
