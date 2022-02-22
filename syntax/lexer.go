package syntax

import (
	"bytes"
	"fmt"
	"github.com/lilac/funlang/token"
	"github.com/rhysd/locerr"
	"io"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*Lexer) stateFn

//const eof = -1

// Lexer instance which contains lexing states.
type Lexer struct {
	state   stateFn
	start   locerr.Pos
	current locerr.Pos
	src     *locerr.Source
	input   *bytes.Reader
	Tokens  chan token.Token
	top     rune
	eof     bool
	// Function called when error occurs.
	// By default it outputs an error to stderr.
	Error func(msg string, pos locerr.Pos)
}

// NewLexer creates new Lexer instance.
func NewLexer(src *locerr.Source) *Lexer {
	start := locerr.Pos{
		Offset: 0,
		Line:   1,
		Column: 1,
		File:   src,
	}
	return &Lexer{
		state:   lex,
		start:   start,
		current: start,
		input:   bytes.NewReader(src.Code),
		src:     src,
		Tokens:  make(chan token.Token),
		Error:   nil,
	}
}

// Lex starts lexing. Lexed tokens will be queued into channel in lexer.
func (l *Lexer) Lex() {
	// Set top to peek current rune
	l.forward()
	for l.state != nil {
		l.state = l.state(l)
	}
}

func (l *Lexer) emit(kind token.Kind) {
	l.Tokens <- token.Token{
		kind,
		l.start,
		l.current,
		l.src,
	}
	l.start = l.current
}

func (l *Lexer) emitIdent(ident string) {
	if len(ident) == 1 {
		// Shortcut because no keyword is one character. It must be identifier
		l.emit(token.Ident)
		return
	}

	switch ident {
	case "true", "false":
		l.emit(token.Bool)
	case "if":
		l.emit(token.If)
	case "then":
		l.emit(token.Then)
	case "else":
		l.emit(token.Else)
	case "let":
		l.emit(token.Let)
	case "in":
		l.emit(token.In)
	case "end":
		l.emit(token.End)
	case "val":
		l.emit(token.Val)
	case "rec":
		l.emit(token.Rec)
	case "not":
		l.emit(token.Not)
	case "match":
		l.emit(token.Match)
	case "with":
		l.emit(token.With)
	case "fn":
		l.emit(token.Fn)
	case "fun":
		l.emit(token.Fun)
	case "type":
		l.emit(token.Type)

	default:
		l.emit(token.Ident)
	}
}

func (l *Lexer) emitIllegal(reason string) {
	l.errmsg(reason)
	t := token.Token{
		token.Illegal,
		l.start,
		l.current,
		l.src,
	}
	l.Tokens <- t
	l.start = l.current
}

func (l *Lexer) expected(s string, actual rune) {
	l.emitIllegal(fmt.Sprintf("Expected %s but got '%c'(%d)", s, actual, actual))
}

func (l *Lexer) unclosedComment(expected string) {
	l.emitIllegal(fmt.Sprintf("Expected '%s' for closing comment but got Eof", expected))
}

func (l *Lexer) forward() {
	r, _, err := l.input.ReadRune()
	if err == io.EOF {
		l.top = 0
		l.eof = true
		return
	}

	if err != nil {
		panic(err)
	}

	if !utf8.ValidRune(r) {
		l.emitIllegal(fmt.Sprintf("Invalid UTF-8 character '%c' (%d)", r, r))
		l.eof = true
		return
	}

	l.top = r
	l.eof = false
}

func (l *Lexer) eat() {
	size := utf8.RuneLen(l.top)
	l.current.Offset += size

	// TODO: Consider \n\r
	if l.top == '\n' {
		l.current.Line++
		l.current.Column = 1
	} else {
		l.current.Column += size
	}

	l.forward()
}

func (l *Lexer) consume() {
	if l.eof {
		return
	}
	l.eat()
	l.start = l.current
}

func (l *Lexer) errmsg(msg string) {
	if l.Error == nil {
		return
	}
	l.Error(msg, l.current)
}

func (l *Lexer) eatIdent() bool {
	if !isLetter(l.top) {
		l.expected("letter for head character of identifer", l.top)
		return false
	}
	l.eat()

	for isLetter(l.top) || isDigit(l.top) || l.top == '\'' {
		l.eat()
	}
	return true
}

func lexComment(l *Lexer) stateFn {
	for {
		if l.eof {
			l.unclosedComment("*")
			return nil
		}
		if l.top == '*' {
			l.eat()
			if l.eof {
				l.unclosedComment(")")
				return nil
			}
			if l.top == ')' {
				l.eat()
				l.emit(token.Comment)
				return lex
			}
		}
		l.eat()
	}
}

func lexLeftParen(l *Lexer) stateFn {
	l.eat()
	if l.top == '*' {
		l.eat()
		return lexComment
	}
	l.emit(token.Lparen)
	return lex
}

func lexAdditiveOp(l *Lexer) stateFn {
	op := token.Plus
	if l.top == '-' {
		op = token.Minus
	}
	l.eat()

	switch l.top {
	case '>':
		if op == token.Minus {
			// Lexing '->'
			l.eat()
			l.emit(token.MinusGreater)
		} else {
			l.emit(op)
		}
	default:
		l.emit(op)
	}

	return lex
}

func lexMultOp(l *Lexer) stateFn {
	op := token.Star
	if l.top == '/' {
		op = token.Slash
	}
	l.eat()
	l.emit(op)

	return lex
}

func lexBar(l *Lexer) stateFn {
	l.eat() // Eat first '|'

	switch l.top {
	case '|':
		l.eat()
		l.emit(token.BarBar)
	default:
		l.emit(token.Bar)
	}

	return lex
}

func lexLogicalAnd(l *Lexer) stateFn {
	prev := l.top
	l.eat()

	if prev != l.top {
		l.expected("logical operator &&", l.top)
		return nil
	}
	l.eat()
	l.emit(token.AndAnd)

	return lex
}

func lexLess(l *Lexer) stateFn {
	l.eat()
	switch l.top {
	case '>':
		l.eat()
		l.emit(token.LessGreater)
	case '=':
		l.eat()
		l.emit(token.LessEqual)
	case '-':
		l.eat()
		l.emit(token.LessMinus)
	default:
		l.emit(token.Less)
	}
	return lex
}

func lexGreater(l *Lexer) stateFn {
	l.eat()
	switch l.top {
	case '=':
		l.eat()
		l.emit(token.GreaterEqual)
	default:
		l.emit(token.Greater)
	}
	return lex
}

// e.g. 123.45e10
func lexNumber(l *Lexer) stateFn {
	tok := token.Int

	// Eat first digit. It's known as digit in lex()
	l.eat()
	for isDigit(l.top) {
		l.eat()
	}

	// Note: Allow 1. as 1.0
	if l.top == '.' {
		tok = token.Float
		l.eat()
		for isDigit(l.top) {
			l.eat()
		}
	}

	if l.top == 'e' || l.top == 'E' {
		tok = token.Float
		l.eat()
		if l.top == '+' || l.top == '-' {
			l.eat()
		}
		if !isDigit(l.top) {
			l.expected("number for exponential part of float literal", l.top)
			return nil
		}
		for isDigit(l.top) {
			l.eat()
		}
	}

	l.emit(tok)
	return lex
}

func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' ||
		'A' <= r && r <= 'Z' ||
		r == '_' ||
		r >= utf8.RuneSelf && unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func lexIdent(l *Lexer) stateFn {
	if !l.eatIdent() {
		return nil
	}
	i := string(l.src.Code[l.start.Offset:l.current.Offset])
	l.emitIdent(i)
	return lex
}

func lexStringLiteral(l *Lexer) stateFn {
	l.eat() // Eat first '"'
	for !l.eof {
		if l.top == '\\' {
			// Skip escape ('\' and next char)
			l.eat()
			l.eat()
		}
		if l.top == '"' {
			l.eat()
			l.emit(token.StringLiteral)
			return lex
		}
		l.eat()
	}
	l.emitIllegal("Unclosed string literal")
	return nil
}

func lexLbracket(l *Lexer) stateFn {
	l.eat() // Eat '['
	l.emit(token.Lbracket)
	return lex
}

func lex(l *Lexer) stateFn {
	for {
		if l.eof {
			l.emit(token.Eof)
			return nil
		}
		switch l.top {
		case '(':
			return lexLeftParen
		case ')':
			l.eat()
			l.emit(token.Rparen)
		case '+':
			return lexAdditiveOp
		case '-':
			return lexAdditiveOp
		case '*':
			return lexMultOp
		case '/':
			return lexMultOp
		case '%':
			l.eat()
			l.emit(token.Percent)
		case '=':
			l.eat()
			switch l.top {
			case '>':
				l.eat()
				l.emit(token.Arrow)
			}
			l.emit(token.Equal)
		case '<':
			return lexLess
		case '>':
			return lexGreater
		case ',':
			l.eat()
			l.emit(token.Comma)
		case '.':
			l.eat()
			l.emit(token.Dot)
		case ';':
			l.eat()
			l.emit(token.Semicolon)
		case '|':
			return lexBar
		case '&':
			return lexLogicalAnd
		case '"':
			return lexStringLiteral
		case ':':
			l.eat()
			l.emit(token.Colon)
		case '[':
			return lexLbracket
		case ']':
			l.eat()
			l.emit(token.Rbracket)
		default:
			switch {
			case unicode.IsSpace(l.top):
				l.consume()
			case isDigit(l.top):
				return lexNumber
			default:
				return lexIdent
			}
		}
	}
}
