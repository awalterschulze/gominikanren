package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type SExpr struct {
	List *List
	Atom *Atom
}

func (s *SExpr) IsAssosiation() bool {
	if s.List == nil {
		//fmt.Printf("list == nil: %v\n", s)
		return false
	}
	if len(s.List.Items) < 1 {
		//fmt.Printf("list2 == nil\n")
		return false
	}
	if !s.List.Items[0].IsVariable() {
		//fmt.Printf("not var %v\n", s.List.Items[0])
		return false
	}
	if len(s.List.Items) == 2 {
		return true
	}
	if len(s.List.Items) == 3 {
		//fmt.Printf("3\n")
		if s.List.Items[1].String() == "." {
			return true
		}
		//fmt.Printf("not .\n")
	}
	return false
}

// IsPair returns whether we can take the car and the cdr of a List without problems.
func (s *SExpr) IsPair() bool {
	if s.List != nil {
		return len(s.List.Items) >= 1
	}
	return false
}

func (s *SExpr) Car() *SExpr {
	if s.List != nil {
		return s.List.Head()
	}
	return s
}

func (s *SExpr) RemoveIDs() {
	if s.Atom != nil && s.Atom.Var != nil {
		s.Atom.Var.ID = 0
	}
	if s.List != nil {
		for i := range s.List.Items {
			s.List.Items[i].RemoveIDs()
		}
	}
}

func (s *SExpr) RemoveQuotes() {
	if s.List != nil {
		s.List.Quoted = ""
		for i := range s.List.Items {
			s.List.Items[i].RemoveQuotes()
		}
	}
}

func (l *List) Head() *SExpr {
	return pushQuote(l.IsQuoted(), l.Items[0])
}

func Cons(car, cdr *SExpr) *SExpr {
	items := []*SExpr{car}
	quoted := false
	if car.List != nil {
		quoted = car.List.IsQuoted()
	}
	if cdr.List != nil {
		items = append(items, cdr.List.Items...)
		quoted = quoted || cdr.List.IsQuoted()
	} else {
		items = append(items, cdr)
	}
	return NewList(quoted, items...)
}

func (s *SExpr) Cdr() *SExpr {
	if s.List != nil {
		tail := s.List.Tail()
		if len(tail.List.Items) == 1 {
			return pushQuote(tail.List.IsQuoted(), tail.List.Items[0])
		}
		return tail
	}
	return NewList(false)
}

func (l *List) Tail() *SExpr {
	if len(l.Items) > 2 && l.Items[1].String() == "." {
		return NewList(l.IsQuoted(), l.Items[2:]...)
	}
	return NewList(l.IsQuoted(), l.Items[1:]...)
}

func (s *SExpr) Equal(ss *SExpr) bool {
	return deriveEqual(s, ss)
}

func (s *SExpr) GoString() string {
	return deriveGoString(s)
}

func (s *SExpr) IsVariable() bool {
	if s.Atom == nil {
		return false
	}
	return s.Atom.Var != nil
}

func (s *SExpr) String() string {
	if s.List != nil {
		return s.List.String()
	}
	return s.Atom.String()
}

type List struct {
	Quoted string
	Items  []*SExpr
}

func NewList(quoted bool, sexprs ...*SExpr) *SExpr {
	q := ""
	if quoted {
		q = "`"
	}
	return &SExpr{
		List: &List{
			Quoted: q,
			Items:  sexprs,
		},
	}
}

func Prepend(quoted bool, s *SExpr, l *List) *SExpr {
	es := make([]*SExpr, len(l.Items)+1)
	copy(es[1:], l.Items)
	es[0] = s
	return NewList(quoted, es...)
}

func Append(l *List, s *SExpr) *List {
	l.Items = append(l.Items, s)
	return l
}

func (l *List) IsNil() bool {
	return len(l.Items) == 0
}

func (l *List) Equal(ll *List) bool {
	if l == nil && ll == nil {
		return true
	}
	if l == nil || ll == nil {
		return false
	}
	if len(l.Items) == 0 && len(ll.Items) == 0 {
		return true
	}
	return deriveEqualItems(l.Items, ll.Items)
}

func (l *List) GoString() string {
	return deriveGoStringList(l)
}

func (l *List) IsQuoted() bool {
	return len(l.Quoted) > 0
}

func pushQuote(quoted bool, s *SExpr) *SExpr {
	if !quoted {
		return s
	}
	if s.List == nil {
		return s
	}
	return NewList(quoted, s.List.Items...)
}

func (l *List) String() string {
	ss := make([]string, len(l.Items))
	for i := range l.Items {
		ss[i] = l.Items[i].String()
	}
	return l.Quoted + "(" + strings.Join(ss, " ") + ")"
}

func NewAssosiation(v, s *SExpr) *SExpr {
	return &SExpr{
		List: &List{
			Quoted: "",
			Items: []*SExpr{
				v,
				NewSymbol("."),
				s,
			},
		},
	}
}

type Atom struct {
	Str    *string
	Symbol *string
	Float  *float64
	Int    *int64
	Var    *Variable
}

func (a *Atom) GoString() string {
	return deriveGoStringAtom(a)
}

func (a *Atom) String() string {
	if a.Str != nil {
		return strconv.Quote(*a.Str)
	}
	if a.Symbol != nil {
		return *a.Symbol
	}
	if a.Float != nil {
		return strconv.FormatFloat(*a.Float, 'f', -1, 64)
	}
	if a.Int != nil {
		return strconv.FormatInt(*a.Int, 10)
	}
	return a.Var.String()
}

func ParseString(s string) (*SExpr, error) {
	ss, err := strconv.Unquote(s)
	if err != nil {
		return nil, err
	}
	return NewString(ss), nil
}

func NewString(s string) *SExpr {
	return &SExpr{
		Atom: &Atom{
			Str: &s,
		},
	}
}

func NewSymbol(s string) *SExpr {
	return &SExpr{
		Atom: &Atom{
			Symbol: &s,
		},
	}
}

func ParseFloat(s string) (*SExpr, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return NewFloat(f), nil
}

func NewFloat(f float64) *SExpr {
	return &SExpr{
		Atom: &Atom{
			Float: &f,
		},
	}
}

func ParseInt(s string) (*SExpr, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewInt(i), nil
}

func NewInt(i int64) *SExpr {
	return &SExpr{
		Atom: &Atom{
			Int: &i,
		},
	}
}

func ParseVariable(s string) (*SExpr, error) {
	if s[0] != ',' {
		return nil, fmt.Errorf("not a variable")
	}
	return NewVariable(s[1:]), nil
}

func NewVariable(s string) *SExpr {
	return &SExpr{
		Atom: &Atom{
			Var: &Variable{
				Name: s,
			},
		},
	}
}

type Variable struct {
	Name string
	ID   int64
}

func (v *Variable) String() string {
	if v.ID == 0 {
		return "," + v.Name
	}
	return "," + v.Name + "_" + strconv.Itoa(int(v.ID))
}

func (v *Variable) GoString() string {
	return deriveGoStringVar(v)
}

func (v *Variable) Equal(vv *Variable) bool {
	return deriveEqualVar(v, vv)
}
