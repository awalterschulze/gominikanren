package comicro

import (
	"fmt"
	"sort"
	"strings"
)

// State is a product of a list of substitutions and a variable counter.
type State struct {
	Substitutions Substitutions
	Counter       uint64
}

// String returns a string representation of State.
func (s *State) String() string {
	if s.Substitutions == nil {
		return fmt.Sprintf("(() . %d)", s.Counter)
	}
	return fmt.Sprintf("(%s . %d)", s.Substitutions.String(), s.Counter)
}

func (s *State) NewVar() (*State, Var) {
	v := NewVar(s.Counter)
	return &State{
		Substitutions: s.Substitutions,
		Counter:       s.Counter + 1,
	}, v
}

func (s *State) Get(v Var) (any, bool) {
	if s == nil {
		return nil, false
	}
	if s.Substitutions == nil {
		return nil, false
	}
	a, ok := s.Substitutions[v]
	return a, ok
}

func (s *State) LenSubstitutions() int {
	if s == nil {
		return 0
	}
	return len(s.Substitutions)
}

func (s *State) AddKeyValue(key Var, value any) *State {
	if s == nil {
		s = EmptyState()
	}
	if s.Substitutions == nil {
		s.Substitutions = map[Var]any{}
	}
	return &State{Substitutions: s.Substitutions.AddKeyValue(key, value), Counter: s.Counter}
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
type Substitutions map[Var]any

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
	ss := make([]string, len(s))
	for i, k := range ks {
		v := s[k]
		ss[i] = fmt.Sprintf("{%v: %v}", k, v)
	}
	return strings.Join(ss, ", ")
}

func (s Substitutions) AddKeyValue(key Var, value any) Substitutions {
	var ss Substitutions
	if s == nil {
		ss = map[Var]any{}
	} else {
		ss = s.Copy()
	}
	ss[key] = value
	return ss
}
