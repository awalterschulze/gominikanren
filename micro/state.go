package micro

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

// EmptyState returns an empty state.
func EmptyState() *State {
	return &State{}
}

// SubPair is a substitution pair
type SubPair struct {
	Key   uint64
	Value *ast.SExpr
}

// Substitutions is a list of substitutions represented by a sexprs pair.
type Substitutions []SubPair

func (s Substitutions) String() string {
	sexprs := make([]*ast.SExpr, len(s))
	sort.Slice(s, func(i, j int) bool { return s[i].Key < s[j].Key })
	for i, pair := range s {
		k, v := pair.Key, pair.Value
		vv := Var(k)
		vvv := ast.Cons(vv, v)
		sexprs[len(s)-1-i] = vvv
	}
	l := ast.NewList(sexprs...).String()
	return l[1 : len(l)-1]
}
