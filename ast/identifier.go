package ast

/*
 Identifiers
 id	::=  	letter⟨letter | digit | ' | _⟩*	alphanumeric
 		⟨! | % | & | $ | # | + | - | / | : | < | = | > | ? | @ | \ | ~ | ` | ^ | | | *⟩+   	symbolic (not allowed for type variables or module language identifiers)
 var	::=	'⟨letter | digit | ' | _⟩*	unconstrained
 		''⟨letter | digit | ' | _⟩*	equality
 longid  ::=	id1.···.idn	qualified (n ≥ 1)
 lab	::=	id	identifier
 		num	number (may not start with 0)
*/

type Identifier struct {
	Value string // final distinct name of the id
	Name  string // the original name written in code
}

type Var struct {
	HasToken
	Id Identifier
}

func (i Identifier) String() string {
	if len(i.Value) > 0 {
		return i.Value
	}
	return i.Name
}

func (v Var) String() string {
	return v.Id.String()
}
