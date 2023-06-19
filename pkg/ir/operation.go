package ir

type Not struct {
	Child Exp
}

func (n Not) tag() expTag {
	return notTag
}

type Neg struct {
	Child Exp
}

func (n Neg) tag() expTag {
	return negTag
}

type Op int8

const (
	Add Op = iota
	Minus
	Mul
	Div
	Mod

	Eq
	NotEq
	Less
	LessEq
	Greater
	GreaterEq

	And
	Or
)

type BinaryOp struct {
	Left  Exp
	Op    Op
	Right Exp
}

func (i BinaryOp) tag() expTag {
	return binaryOpTag
}

func NewBinaryOp(op Op, left, right Exp) *BinaryOp {
	return &BinaryOp{
		Left:  left,
		Op:    op,
		Right: right,
	}
}
