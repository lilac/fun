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
	pattern ast.Pattern
	patterns []ast.Pattern
	match []ast.Match
	funBind	[]ast.FunBind
	dec []ast.Dec
	mod *ast.Module
}

%token<token> Illegal
%token<token> Comment
%token<token> LParen
%token<token> RParen
%token<token> Ident
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

%right prec_if
%right prec_fn
%left Bar
%right Arrow
%left Comma Semicolon
%left BarBar
%left AndAnd
%left Equal LessGreater Less Greater LessEqual GreaterEqual
%left Plus Minus
%left Star Slash Percent
%right prec_unary_minus Not
%left prec_app
%left Dot

%type<mod> module
%type<dec> dec
%type<exp> exp con simple_exp
%type<pattern> pattern
%type<match> match
%type<patterns> patterns
%type<funBind> fun_bind

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
|	dec Fun fun_bind
	{
		dec := NewFunDec($3)
		$$ = append($1, dec)
	}

fun_bind:
	Ident patterns Equal exp
	{
		bind := NewFunBind($1, $2, nil, $4)
		$$ = []ast.FunBind{*bind}
	}
|	fun_bind Bar Ident patterns Equal exp
	{
		bind := NewFunBind($3, $4, nil, $6)
        	$$ = append($1, *bind)
	}

patterns:
	pattern
	{ $$ = []ast.Pattern{$1} }
|	patterns pattern
	{ $$ = append($1, $2) }

simple_exp:
	con
	{ $$ = $1 }
|	Ident
	{ $$ = NewVar($1) }
|	LParen exp RParen
	{ $$ = $2 }

exp:
	simple_exp
	{ $$ = $1 }
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

|	exp Comma exp
	{ $$ = NewTuple($1, $3) }
|	exp Semicolon exp
	{ $$ = NewSequence($1, $3) }
|	exp simple_exp
	%prec prec_app
	{ $$ = &ast.Apply{$1, $2} }

|	If exp Then exp Else exp
	%prec prec_if
	{ $$ = NewIfThen($1, $2, $4, $6) }
|	Let dec In exp End
	{ $$ = NewLet($1, $2, $4) }
|	Fn match
	%prec prec_fn
	{ $$ = NewFn($1, $2) }

match:
	pattern Arrow exp
	{
		m := NewMatch($1, $3)
		$$ = []ast.Match{*m}
	}
|	match Bar pattern Arrow exp
	{
		m := NewMatch($3, $5)
		$$ = append($1, *m)
	}

pattern:
	con
	{ $$ = NewConstPattern($1) }
|	Ident
	{ $$ = NewVarPattern($1) }

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
