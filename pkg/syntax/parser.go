package syntax

import (
	"fmt"
	merror "github.com/hashicorp/go-multierror"
	"github.com/lilac/fun-lang/pkg/ast"
	"io"
	"os"
	"strings"
)

type Source struct {
	io.Reader
	Path string
}

func NewSourceFromFile(file string) (*Source, error) {
	var reader *os.File
	var err error
	if file == "" {
		reader = os.Stdin
	} else {
		reader, err = os.Open(file)
	}
	if err != nil {
		return nil, err
	} else {
		return &Source{reader, reader.Name()}, nil
	}
}

// NewDummySource make *Source with passed code. This is used for tests.
func NewDummySource(code string) *Source {
	return &Source{strings.NewReader(code), "<dummy>"}
}

/*func Parse(source *locerr.Source) (*ast.Module, error) {
	input := bytes.NewReader(source.Code)
	src := &Source{input, source.Path}
	return ParseReader(src)
}*/

func Parse(src *Source) (*ast.Module, error) {
	var err error
	lexer := NewLexer(src)
	lexer.OnError = func(s string) {
		e := fmt.Errorf("parse error at %v: %s", src.Path, s)
		err = merror.Append(err, e)
	}
	parser := funNewParser()
	status := parser.Parse(lexer)
	//fmt.Printf("Parse %s: status = %d\n", src.Path, status)
	module := parser.(*funParserImpl).lval.mod
	if err != nil {
		return module, err
	} else if status == 0 {
		return module, nil
	} else {
		pos := lexer.Current()
		err := fmt.Errorf("parse error at line %d column %d", pos.Line, pos.Column)
		return module, err
	}
}
