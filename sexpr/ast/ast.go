package ast

import (
	"fmt"
	"strconv"
)

func Cons(car *SExpr, cdr *SExpr) *SExpr {
	return &SExpr{Pair: &Pair{Car: car, Cdr: cdr}}
}

type SExpr struct {
	Pair *Pair
	Atom *Atom
}

func (s *SExpr) Cdr() *SExpr {
	if s.Pair != nil {
		return s.Pair.Cdr
	}
	panic("not implemented for atom")
}

func (s *SExpr) Car() *SExpr {
	if s.Pair != nil {
		return s.Pair.Car
	}
	panic("not implemented for atom")
}

type Pair struct {
	Car *SExpr
	Cdr *SExpr
}

func (s *SExpr) String() string {
	if s == nil {
		return "()"
	}
	if s.Atom != nil {
		return s.Atom.String()
	}
	return "(" + s.Pair.String() + ")"
}

func (p *Pair) String() string {
	if p == nil {
		return ""
	}
	if p.Cdr == nil {
		return p.Car.String()
	}
	if p.Cdr.Pair != nil {
		return p.Car.String() + " " + p.Cdr.Pair.String()
	}
	return p.Car.String() + " . " + p.Cdr.String()
}

func (s *SExpr) Equal(ss *SExpr) bool {
	return deriveEqual(s, ss)
}

func (s *SExpr) GoString() string {
	return deriveGoString(s)
}

func (s *SExpr) IsVariable() bool {
	return s != nil && s.Atom != nil && s.Atom.Var != nil
}

func (s *SExpr) IsPair() bool {
	return s != nil && s.Pair != nil
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
}

func (v *Variable) String() string {
	return "," + v.Name
}

func (v *Variable) GoString() string {
	return deriveGoStringVar(v)
}

func (v *Variable) Equal(vv *Variable) bool {
	return deriveEqualVar(v, vv)
}

func NewList(ss ...*SExpr) *SExpr {
	if len(ss) == 0 {
		return nil
	}
	if len(ss) == 1 {
		return Cons(ss[0], nil)
	}
	return Cons(ss[0], NewList(ss[1:]...))
}
