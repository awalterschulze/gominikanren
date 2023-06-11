package comicro

import (
	"fmt"
	"sort"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// State is a product of a list of substitutions and a variable counter.
type State struct {
	Substitutions
	Counter uint64
}

// String returns a string representation of State.
func (s *State) String() string {
	if s.Substitutions == nil {
		return fmt.Sprintf("(() . %d)", s.Counter)
	}
	return fmt.Sprintf("(%s . %d)", s.Substitutions.String(), s.Counter)
}

func (s *State) AddCounter() *State {
	return &State{
		Substitutions: s.Substitutions,
		Counter:       s.Counter + 1,
	}
}

func (s *State) Copy() *State {
	return &State{
		Substitutions: s.Substitutions.Copy(),
		Counter:       s.Counter,
	}
}

// EmptyState returns an empty state.
func EmptyState() *State {
	return &State{}
}

// Substitutions is a list of substitutions represented by a sexprs pair.
type Substitutions map[Var]*ast.SExpr

func (s Substitutions) Copy() Substitutions {
	m := make(Substitutions, len(s))
	for k, v := range s {
		m[k] = v
	}
	return m
}

func (s Substitutions) String() string {
	ks := keys(s)
	sort.Slice(ks, func(i, j int) bool { return ks[i] < ks[j] })
	sexprs := make([]*ast.SExpr, len(s))
	for i, k := range ks {
		v := s[k]
		vv := Var(k)
		vvv := ast.Cons(vv.SExpr(), v)
		sexprs[len(s)-1-i] = vvv
	}
	l := ast.NewList(sexprs...).String()
	return l[1 : len(l)-1]
}

func (s Substitutions) AddPair(key Var, value *ast.SExpr) Substitutions {
	var ss Substitutions
	if s == nil {
		ss = map[Var]*ast.SExpr{}
	} else {
		ss = s.Copy()
	}
	ss[key] = value
	return ss
}

func (s Substitutions) Get(key Var) (*ast.SExpr, bool) {
	if s == nil {
		return nil, false
	}
	v, ok := s[key]
	return v, ok
}
