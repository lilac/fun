package syntax

import (
	"bytes"
	"fmt"
	"github.com/lilac/fun-lang/ast"
	"github.com/rhysd/locerr"
	"io"
)

func Parse(source *locerr.Source) (*ast.Module, error) {
	input := bytes.NewReader(source.Code)
	return ParseReader(input)
}

func ParseReader(src io.Reader) (*ast.Module, error) {
	lexer := NewLexer(src)
	parser := funNewParser()
	status := parser.Parse(lexer)
	fmt.Printf("Parse result: %d\n", status)
	if status == 0 {
		return parser.(*funParserImpl).lval.mod, nil
	} else {
		return nil, fmt.Errorf("parse error at line %d column %d", lexer.Line(), lexer.Column())
	}
}
