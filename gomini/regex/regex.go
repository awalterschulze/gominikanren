package regex

type RegexType int32

const (
	EmptySetType RegexType = 0
	EmptyStrType RegexType = 1
	CharType     RegexType = 2
	OrType       RegexType = 3
	ConcatType   RegexType = 4
	StarType     RegexType = 5
)

type Regex struct {
	Type *RegexType
	Char *rune
	R1   *Regex
	R2   *Regex
}

func EmptySet() *Regex {
	typ := EmptySetType
	return &Regex{Type: &typ}
}

func EmptyStr() *Regex {
	typ := EmptyStrType
	return &Regex{Type: &typ}
}

func Char(c rune) *Regex {
	typ := CharType
	return &Regex{Type: &typ, Char: &c}
}

func CharPtr(c *rune) *Regex {
	typ := CharType
	return &Regex{Type: &typ, Char: c}
}

func Or(a, b *Regex) *Regex {
	typ := OrType
	return &Regex{Type: &typ, R1: a, R2: b}
}

func Concat(a, b *Regex) *Regex {
	typ := ConcatType
	return &Regex{Type: &typ, R1: a, R2: b}
}

func Star(a *Regex) *Regex {
	typ := StarType
	return &Regex{Type: &typ, R1: a}
}

func (r *Regex) String() string {
	if r == nil {
		return "nil"
	}
	if r.Type == nil {
		return "nilType"
	}
	switch *r.Type {
	case EmptySetType:
		return "∅"
	case EmptyStrType:
		return "ε"
	case CharType:
		return string(*r.Char)
	case OrType:
		return "(" + r.R1.String() + "|" + r.R2.String() + ")"
	case ConcatType:
		return r.R1.String() + r.R2.String()
	case StarType:
		return "(" + r.R1.String() + ")*"
	}
	panic("unreachable")
}
