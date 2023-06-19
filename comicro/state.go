package comicro

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// State is a product of a list of substitutions and a variable counter.
type State struct {
	Substitutions map[Var]any
	Names         map[Var]string
	Pointers      map[Var]reflect.Value
	ReifyNames    []ReifyName
	FirstVar      *Var
	Counter       uint64
}

type ReifyName func(varType any, name string) (any, bool)

// NewEmptyState returns an empty state.
func NewEmptyState() *State {
	return &State{}
}

func (s *State) WithReifyNames(reifyNames ...ReifyName) *State {
	res := s.Copy()
	res.ReifyNames = append(s.ReifyNames, reifyNames...)
	return res
}

// String returns a string representation of State.
func (s *State) String() string {
	if s.Substitutions == nil {
		return fmt.Sprintf("(() . %d)", s.Counter)
	}
	ks := keys(s.Substitutions)
	sort.Slice(ks, func(i, j int) bool { return ks[i] < ks[j] })
	ss := make([]string, len(s.Substitutions))
	for i, k := range ks {
		v := s.Substitutions[k]
		vstr := fmt.Sprintf("%v", v)
		kstr := s.GetName(k)
		if vvar, ok := v.(Var); ok {
			vstr = s.GetName(vvar)
		}
		kstr = "," + kstr
		ss[i] = fmt.Sprintf("{%s: %s}", kstr, vstr)
	}
	return fmt.Sprintf("(%s . %d)", strings.Join(ss, ", "), s.Counter)
}

func (s *State) Equal(other *State) bool {
	return s.String() == other.String()
}

type Var uintptr

func NewVar[A any](s *State, typ *A) (*State, *A) {
	return NewVarWithName(s, "v"+strconv.Itoa(int(s.Counter)), typ)
}

func NewVarWithName[A any](s *State, name string, typ *A) (*State, *A) {
	if s == nil {
		s = NewEmptyState()
	}
	v := reflect.New(reflect.TypeOf(typ).Elem())
	key := Var(v.Pointer())
	names := copyMap(s.Names)
	names[key] = name
	pointers := copyMap(s.Pointers)
	pointers[key] = v
	res := &State{
		Substitutions: s.Substitutions,
		Counter:       s.Counter + 1,
		Pointers:      pointers,
		Names:         names,
		FirstVar:      s.FirstVar,
		ReifyNames:    s.ReifyNames,
	}
	if s.FirstVar == nil {
		res.FirstVar = &key
	}
	return res, v.Interface().(*A)
}

func (s *State) GetFirstVar() *Var {
	return s.FirstVar
}

func (s *State) GetVar(a any) (Var, bool) {
	if s == nil {
		return 0, false
	}
	if avar, ok := a.(Var); ok {
		return avar, true
	}
	v := reflect.ValueOf(a)
	key := Var(v.Pointer())
	_, ok := s.Pointers[key]
	return key, ok
}

func (s *State) LookupValue(key Var) any {
	placeholder, ok := s.Pointers[key]
	if !ok {
		panic(fmt.Sprintf("Var %v not found", key))
	}
	return placeholder.Interface()
}

func (s *State) SameVar(a, b Var) bool {
	return s.Pointers[a].Pointer() == s.Pointers[b].Pointer()
}

func (s *State) GetName(v Var) string {
	if s == nil {
		return "v0"
	}
	if s != nil && s.Names != nil {
		name, ok := s.Names[v]
		if ok {
			return name
		}
	}
	return "v?"
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

func (s *State) AddKeyValue(key Var, value any) *State {
	var ss *State
	if s == nil {
		ss = NewEmptyState()
		ss.Substitutions = make(map[Var]any)
	} else {
		ss = s.Copy()
	}
	ss.Substitutions[key] = value
	return ss
}

func (s *State) Copy() *State {
	if s == nil {
		return nil
	}
	names := copyMap(s.Names)
	substitutions := copyMap(s.Substitutions)
	pointers := copyMap(s.Pointers)
	return &State{
		Substitutions: substitutions,
		Counter:       s.Counter,
		Pointers:      pointers,
		FirstVar:      s.FirstVar,
		Names:         names,
		ReifyNames:    s.ReifyNames,
	}
}

func (s *State) CopyWithoutSubstitutions() *State {
	res := s.Copy()
	res.Substitutions = make(map[Var]any)
	return res
}

func copyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (s *State) GetReifyName(v Var) any {
	varType := s.LookupValue(v)
	name := "," + s.GetName(v)
	for _, r := range s.ReifyNames {
		if val, ok := r(varType, name); ok {
			return val
		}
	}
	panic(fmt.Sprintf("unable to reify var (%s) with type: %T", name, varType))
}
