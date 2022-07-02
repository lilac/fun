package main

import (
	"flag"
	"fmt"
	"github.com/lilac/fun-lang/pkg/compiler"
	"github.com/lilac/fun-lang/pkg/syntax"
	"os"
)

var (
	help = flag.Bool("help", false, "Show this help")
)

const usageHeader = `Usage: fun [flags] [file]

  Compiler of the Fun language.
  When [file] is not given, it will read the source code from STDIN.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	var err error
	var file string = ""

	if flag.NArg() > 0 {
		file = flag.Arg(0)
	}
	src, err := syntax.NewSourceFromFile(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on opening file: %s\n", err.Error())
		os.Exit(1)
	}

	err = compiler.Compile(src)
	handleError(err)

	// code generation
	//fun := codegen.GenFunction(nil)
	//printer.Fprint(os.Stdout, token.NewFileSet(), fun)
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
