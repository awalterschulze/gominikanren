package gomini

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// State is a product of a list of substitutions and a variable counter.
type State struct {
	substitutions map[Var]any
	queryVar      *Var
	counter       uint64
	pointers      map[Var]reflect.Value

	names       map[Var]string
	varCreators []VarCreator
}

type VarCreator func(varType any, name string) (any, bool)

// NewEmptyState returns an empty state.
func NewEmptyState() *State {
	return &State{}
}

type Var uintptr

func NewVar[A any](s *State, typ A) (*State, A) {
	return NewVarWithName(s, "v"+strconv.Itoa(int(s.counter)), typ)
}

func NewVarWithName[A any](s *State, name string, typ A) (*State, A) {
	if s == nil {
		s = NewEmptyState()
	}
	vvalue := s.newVarValue(typ, name)
	v := reflect.ValueOf(vvalue)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		// call to Pointer only works for these types and otherwise panics
	default:
		panic("cannot make a variable that is not a pointer, slice or map " + v.Type().String())
	}
	key := Var(v.Pointer())
	names := copyMap(s.names)
	names[key] = name
	pointers := copyMap(s.pointers)
	pointers[key] = v
	res := &State{
		substitutions: s.substitutions,
		counter:       s.counter + 1,
		pointers:      pointers,
		names:         names,
		queryVar:      s.queryVar,
		varCreators:   s.varCreators,
	}
	if s.queryVar == nil {
		res.queryVar = &key
	}
	return res, v.Interface().(A)
}

func (s *State) newVarValue(varType any, name string) any {
	for _, create := range s.varCreators {
		if val, ok := create(varType, name); ok {
			return val
		}
	}
	return reflect.New(reflect.TypeOf(varType).Elem()).Interface()
}

func (s *State) GetQueryVar() *Var {
	return s.queryVar
}

func (s *State) castVar(a any) (Var, bool) {
	if s == nil {
		return 0, false
	}
	if avar, ok := a.(Var); ok {
		return avar, true
	}
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		// call to Pointer only works for these types and otherwise panics
	default:
		return 0, false
	}
	key := Var(v.Pointer())
	_, ok := s.pointers[key]
	return key, ok
}

func (s *State) lookupPlaceholderValue(key Var) any {
	placeholder, ok := s.pointers[key]
	if !ok {
		panic(fmt.Sprintf("Var %v not found", key))
	}
	return placeholder.Interface()
}

func (s *State) isSameVar(a, b Var) bool {
	return s.pointers[a].Pointer() == s.pointers[b].Pointer()
}

func (s *State) findSubstitution(v Var) (any, bool) {
	if s == nil {
		return nil, false
	}
	if s.substitutions == nil {
		return nil, false
	}
	a, ok := s.substitutions[v]
	return a, ok
}

func (s *State) AddKeyValue(key Var, value any) *State {
	var ss *State
	if s == nil {
		ss = NewEmptyState()
		ss.substitutions = make(map[Var]any)
	} else {
		ss = s.Copy()
	}
	ss.substitutions[key] = value
	return ss
}

func (s *State) Copy() *State {
	if s == nil {
		return nil
	}
	names := copyMap(s.names)
	substitutions := copyMap(s.substitutions)
	pointers := copyMap(s.pointers)
	return &State{
		substitutions: substitutions,
		counter:       s.counter,
		pointers:      pointers,
		queryVar:      s.queryVar,
		names:         names,
		varCreators:   s.varCreators,
	}
}

func copyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// Provide optional VarCreator, which is used to give variables printable names for debugging purposes.
func (s *State) WithVarCreators(varCreators ...VarCreator) *State {
	res := s.Copy()
	res.varCreators = append(s.varCreators, varCreators...)
	return res
}

func (s *State) Equal(other *State) bool {
	return s.String() == other.String()
}

func (s *State) getName(v Var) string {
	if s == nil {
		return "v0"
	}
	if s != nil && s.names != nil {
		name, ok := s.names[v]
		if ok {
			return name
		}
	}
	return "v?"
}

// String returns a string representation of State.
func (s *State) String() string {
	if s.substitutions == nil {
		return fmt.Sprintf("(() . %d)", s.counter)
	}
	ks := keys(s.substitutions)
	sort.Slice(ks, func(i, j int) bool { return ks[i] < ks[j] })
	ss := make([]string, len(s.substitutions))
	for i, k := range ks {
		v := s.substitutions[k]
		vstr := fmt.Sprintf("%v", v)
		kstr := s.getName(k)
		if vvar, ok := v.(Var); ok {
			vstr = s.getName(vvar)
		}
		kstr = "," + kstr
		ss[i] = fmt.Sprintf("{%s: %s}", kstr, vstr)
	}
	return fmt.Sprintf("(%s . %d)", strings.Join(ss, ", "), s.counter)
}
