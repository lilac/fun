package ir

type Unit struct {
}

func (u Unit) tag() expTag {
	return unitTag
}

type Bool struct {
	Value bool
}

func (b Bool) tag() expTag {
	return boolTag
}

type Char struct {
	Value rune
}

func (c Char) tag() expTag {
	return charTag
}

type Int struct {
	Value int
}

func (i Int) tag() expTag {
	return intTag
}

type Float struct {
	Value float64
}

func (f Float) tag() expTag {
	return floatTag
}

type String struct {
	Value string
}

func (s String) tag() expTag {
	return stringTag
}
