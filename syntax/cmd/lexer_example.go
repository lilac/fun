package main

import (
	"flag"
	"fmt"
	"github.com/lilac/fun-lang/syntax"
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

	var src *syntax.Source
	var err error

	if flag.NArg() == 0 {
		src, err = syntax.NewSourceFromFile("")
	} else {
		file := filepath.FromSlash(flag.Arg(0))
		src, err = syntax.NewSourceFromFile(file)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on opening source: %s\n", err.Error())
		os.Exit(4)
	}
	LexFile(src)
}

func LexFile(src *syntax.Source) {
	lex := syntax.NewLexer(src)

	for _, tok := range lex.LexAll() {
		switch tok.Kind {
		case syntax.Illegal:
			fmt.Printf("Lexing invalid token at %v\n", tok.Start())
			return
		default:
			fmt.Printf("%s ", tok.Value)
		}
	}

}
