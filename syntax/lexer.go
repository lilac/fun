package syntax

import (
	"bufio"
	"fmt"
	"github.com/lilac/fun-lang/token"
	"io"
	"regexp"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*Lexer) stateFn

// Lexer instance which contains lexing states.
type Lexer struct {
	state   stateFn
	start   token.Position
	current token.Position
	src     string // the path to the source code
	input   io.RuneReader
	token   *token.Token
	buffer  []rune // partial runes of the current token being parsed
	top     rune
	eof     bool
	OnError ErrorFun // error listener
}

func (l *Lexer) Lex(lval *funSymType) int {
	if l.state != nil && l.Next() != nil {
		lval.token = l.token
		return lval.token.Kind
	}
	return Eof
}

// Next consumes the input until a token is parsed, and nil is returned when end of file reached.
func (l *Lexer) Next() *token.Token {
	l.token = nil // reset the current token
	for l.state != nil && l.token == nil {
		l.state = l.state(l)
	}
	return l.token
}

func (l Lexer) Error(s string) {
	fmt.Printf("Parsing error near %s(%v): %s\n", l.src, l.current, s)
	if l.OnError != nil {
		l.OnError(s)
	}
}

func (l Lexer) Current() token.Position {
	return l.current
}

func (l Lexer) Text() string {
	return string(l.buffer)
}

func (l Lexer) newToken(kind int) *token.Token {
	tok := token.NewToken(l.Text())
	tok.Location = token.Location{
		Start: l.start,
		End:   l.current,
		Path:  l.src,
	}
	tok.Kind = kind
	return tok
}

// NewLexer creates new Lexer instance.
func NewLexer(src *Source) *Lexer {
	start := token.Position{
		Line:   1,
		Column: 1,
	}
	l := &Lexer{
		state:   lex,
		start:   start,
		current: start,
		input:   bufio.NewReader(src.Reader),
		src:     src.Path,
		buffer:  nil,
	}
	// Look ahead to start parsing
	l.lookAhead()
	return l
}

// LexAll starts lexing.
func (l *Lexer) LexAll() []*token.Token {
	tokens := make([]*token.Token, 0, 10)
	for l.state != nil {
		if l.Next() != nil {
			tokens = append(tokens, l.token)
		} else {
			break
		}
	}
	return tokens
}

func (l *Lexer) emit(kind int) {
	l.token = l.newToken(kind)
	// reset the start position
	l.start = l.current
	l.buffer = nil // reset the buffer
}

func (l *Lexer) emitIdent(ident string) {
	opReg := regexp.MustCompile(`[+\-*/^=<>]+`)
	if opReg.MatchString(ident) {
		l.emit(Op)
	}
	if len(ident) == 1 {
		// Shortcut because no keyword is one character. It must be identifier
		l.emit(Ident)
		return
	}

	switch ident {
	case "true", "false":
		l.emit(Bool)
	case "if":
		l.emit(If)
	case "then":
		l.emit(Then)
	case "else":
		l.emit(Else)
	case "let":
		l.emit(Let)
	case "in":
		l.emit(In)
	case "end":
		l.emit(End)
	case "val":
		l.emit(Val)
	case "rec":
		l.emit(Rec)
	case "not":
		l.emit(Not)
	case "match":
		l.emit(Match)
	case "with":
		l.emit(With)
	case "fn":
		l.emit(Fn)
	case "fun":
		l.emit(Fun)
	case "type":
		l.emit(Type)

	default:
		l.emit(Ident)
	}
}

func (l *Lexer) emitIllegal(reason string) {
	l.reportError(reason)
	l.emit(Illegal)
}

func (l *Lexer) expected(s string, actual rune) {
	l.emitIllegal(fmt.Sprintf("Expected %s but got '%c'(%d)", s, actual, actual))
}

func (l *Lexer) unclosedComment(expected string) {
	l.emitIllegal(fmt.Sprintf("Expected '%s' for closing comment but got Eof", expected))
}

// look ahead by one char, and assign top and eof.
func (l *Lexer) lookAhead() {
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

// move ahead
func (l *Lexer) shift() {
	// TODO: Consider \n\r
	if l.top == '\n' {
		l.current.Line++
		l.current.Column = 1
	} else {
		l.current.Column += 1
	}
}

// update the current position, and push the top char to buffer, then look ahead.
func (l *Lexer) eat() {
	l.shift()
	l.buffer = append(l.buffer, l.top)
	l.lookAhead()
}

// skip the current char
func (l *Lexer) consume() {
	if l.eof {
		return
	}
	l.shift()
	l.lookAhead()
	l.start = l.current
}

func (l *Lexer) reportError(msg string) {
	l.Error(msg)
}

func (l *Lexer) eatIdent() bool {
	if !isLetter(l.top) {
		l.expected("letter for head character of identifier", l.top)
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
				l.emit(Comment)
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
	l.emit(LParen)
	return lex
}

func lexAdditiveOp(l *Lexer) stateFn {
	op := Plus
	if l.top == '-' {
		op = Minus
	}
	l.eat()

	switch l.top {
	case '>':
		if op == Minus {
			// Lexing '->'
			l.eat()
			l.emit(MinusGreater)
		} else {
			l.emit(op)
		}
	default:
		l.emit(op)
	}

	return lex
}

func lexMultOp(l *Lexer) stateFn {
	op := Star
	if l.top == '/' {
		op = Slash
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
		l.emit(BarBar)
	default:
		l.emit(Bar)
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
	l.emit(AndAnd)

	return lex
}

func lexLess(l *Lexer) stateFn {
	l.eat()
	switch l.top {
	case '>':
		l.eat()
		l.emit(LessGreater)
	case '=':
		l.eat()
		l.emit(LessEqual)
	case '-':
		l.eat()
		l.emit(LessMinus)
	default:
		l.emit(Less)
	}
	return lex
}

func lexGreater(l *Lexer) stateFn {
	l.eat()
	switch l.top {
	case '=':
		l.eat()
		l.emit(GreaterEqual)
	default:
		l.emit(Greater)
	}
	return lex
}

// e.g. 123.45e10
func lexNumber(l *Lexer) stateFn {
	tok := Int

	// Eat first digit. It's known as digit in lex()
	l.eat()
	for isDigit(l.top) {
		l.eat()
	}

	// Note: Allow 1. as 1.0
	if l.top == '.' {
		tok = Float
		l.eat()
		for isDigit(l.top) {
			l.eat()
		}
	}

	if l.top == 'e' || l.top == 'E' {
		tok = Float
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
	i := l.Text()
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
			l.emit(StringLiteral)
			return lex
		}
		l.eat()
	}
	l.emitIllegal("Unclosed string literal")
	return nil
}

func lexLBracket(l *Lexer) stateFn {
	l.eat() // Eat '['
	l.emit(LBracket)
	return lex
}

// lex is the initial state transformation function. It should eat/consume at least one char to move ahead,
// or stop when eof is true.
func lex(l *Lexer) stateFn {
	if l.eof {
		return nil
	}
	switch l.top {
	case '(':
		return lexLeftParen
	case ')':
		l.eat()
		l.emit(RParen)
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
		l.emit(Percent)
	case '=':
		l.eat()
		switch l.top {
		case '>':
			l.eat()
			l.emit(Arrow)
		default:
			l.emit(Equal)
		}
	case '<':
		return lexLess
	case '>':
		return lexGreater
	case ',':
		l.eat()
		l.emit(Comma)
	case '.':
		l.eat()
		l.emit(Dot)
	case ';':
		l.eat()
		l.emit(Semicolon)
	case '|':
		return lexBar
	case '&':
		return lexLogicalAnd
	case '"':
		return lexStringLiteral
	case ':':
		l.eat()
		l.emit(Colon)
	case '[':
		return lexLBracket
	case ']':
		l.eat()
		l.emit(RBracket)
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
	return lex
}
