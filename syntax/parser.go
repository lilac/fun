package syntax

import (
	"github.com/lilac/funlang/ast"
	"github.com/lilac/funlang/token"
	"github.com/rhysd/locerr"
)

// funLex implement the interface of funLexer
type funLex struct {
	lastToken *token.Token
	tokens    chan token.Token
	err       *locerr.Error
	result    *ast.Module
}

// Lex is called by the parser to get next token.
func (l *funLex) Lex(lval *funSymType) int {
	for {
		select {
		case t := <-l.tokens:
			lval.token = &t
			switch t.Kind {
			case token.Eof, token.Illegal:
				return eof
			case token.Comment:
				continue
			}
			l.lastToken = &t
			return int(t.Kind) + Illegal
		}
	}
}

// Error is called by the parser on parsing errors.
func (l *funLex) Error(msg string) {
	if l.err == nil {
		if l.lastToken != nil {
			l.err = locerr.ErrorAt(l.lastToken.Start, msg)
		} else {
			l.err = locerr.NewError(msg)
		}
	} else {
		if l.lastToken != nil {
			l.err = l.err.NoteAt(l.lastToken.Start, msg)
		} else {
			l.err = l.err.Note(msg)
		}
	}
}

func Parse(src *locerr.Source) (*ast.Module, error) {
	var lexErr *locerr.Error
	l := NewLexer(src)
	l.Error = func(msg string, pos locerr.Pos) {
		if lexErr == nil {
			lexErr = locerr.ErrorAt(pos, msg)
		} else {
			lexErr = lexErr.NoteAt(pos, msg)
		}
	}
	go l.Lex()
	parsed, err := ParseTokens(l.Tokens)
	if lexErr != nil {
		return nil, lexErr.Note("Lexing source into tokens failed")
	}
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

// ParseTokens parses given tokens and returns parsed module.
func ParseTokens(tokens chan token.Token) (*ast.Module, error) {
	funErrorVerbose = true

	l := &funLex{tokens: tokens}
	ret := funParse(l)

	if l.err != nil {
		l.Error("Error while parsing")
		return nil, l.err
	}

	root := l.result
	if ret != 0 || root == nil {
		panic("FATAL: Parse failed for unknown reason")
	}

	return root, nil
}
