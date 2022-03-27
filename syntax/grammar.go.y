/*
 This is the grammar of the fun language.
 To build it:
 go get golang.org/x/tools/cmd/goyacc
 goyacc -p fun grammar.go.y (produces grammar.go)
*/

%{

package syntax

import (
	"github.com/lilac/fun-lang/token"
	"github.com/lilac/fun-lang/ast"
)
%}

%union {
	token *token.Token
	exp ast.Exp
	dec []ast.Dec
	mod *ast.Module
}

%token<token> Illegal
%token<token> Comment
%token<token> LParen
%token<token> RParen
%token<token> Ident
%token<token> Op
%token<token> Bool
%token<token> Not
%token<token> Int
%token<token> Float
%token<token> Minus
%token<token> Plus
%token<token> Equal
%token<token> LessGreater
%token<token> LessEqual
%token<token> Less
%token<token> Greater
%token<token> GreaterEqual
%token<token> If
%token<token> Then
%token<token> Else
%token<token> Let
%token<token> In
%token<token> End
%token<token> Val
%token<token> Rec
%token<token> Comma
%token<token> Dot
%token<token> LessMinus
%token<token> Semicolon
%token<token> Star
%token<token> Slash
%token<token> BarBar
%token<token> AndAnd
%token<token> StringLiteral
%token<token> Percent
%token<token> Match
%token<token> With
%token<token> Bar
%token<token> MinusGreater
%token<token> Arrow
%token<token> Fn
%token<token> Fun
%token<token> Colon
%token<token> Type
%token<token> LBracket
%token<token> RBracket

%left Comma
%left BarBar
%left AndAnd
%left Equal LessGreater Less Greater LessEqual GreaterEqual
%left Plus Minus
%left Star Slash Percent
%left prec_app
%right prec_unary_minus Not
%left Dot

%type<mod> module
%type<dec> dec
%type<exp> exp con

%start module

%%

module:
	dec
	{
	 	$$ = &ast.Module{$1}
	 	funrcvr.lval.mod = $$
	}

dec:
	/* empty */
	{ $$ = []ast.Dec{} }
|	dec Val Ident Equal exp
 	{
 		dec := NewValDec($3, $5)
 		$$ = append($1, dec)
 	}

exp:
	con
	{ $$ = $1 }
|	Ident
	{ $$ = NewVar($1) }
|	exp Plus exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp Minus exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp Star exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp Slash exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp Percent exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp Equal exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp LessGreater exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp Less exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp LessEqual exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp Greater exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp GreaterEqual exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp AndAnd exp
	{ $$ = NewInfixApp($1, $2, $3) }
|	exp BarBar exp
	{ $$ = NewInfixApp($1, $2, $3) }

|	Not exp
	{ $$ = NewNot($1, $2) }
|	Minus exp
	%prec prec_unary_minus
	{ $$ = NewNeg($1, $2) }
|	LParen exp RParen
	{ $$ = $2 }

con:
	LParen RParen
	{ $$ = NewUnit($1) }
|	Bool
	{ $$ = NewBool($1) }
|	Int
	{ $$ = NewInt($1, funlex.Error) }
|	Float
	{ $$ = NewFloat($1, funlex.Error) }
|	StringLiteral
	{ $$ = NewString($1, funlex.Error) }
%%

// The parser expects the lexer to return 0 on the end of file.
const Eof = 0
