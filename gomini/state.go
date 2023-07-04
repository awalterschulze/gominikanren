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
	placeholders  map[Var]any
	queryVar      *Var

	names       map[Var]string
	varCreators []VarCreator
}

// VarCreator is a function that creates a variable of a given type and name for debugging purposes.
type VarCreator func(varType any, name string) (any, bool)

// NewState returns an empty state.
// Provide optional VarCreator, which is used to give variables printable names for debugging purposes.
func NewState(varCreators ...VarCreator) *State {
	return &State{varCreators: varCreators}
}

type Var uintptr

func NewVar[A any](s *State, typ A) (*State, A) {
	return newVarWithName(s, "v"+strconv.Itoa(len(s.placeholders)), typ)
}

func newVarWithName[A any](s *State, name string, typ A) (*State, A) {
	if s == nil {
		s = NewState()
	}
	res := &State{
		substitutions: s.substitutions,
		placeholders:  copyMap(s.placeholders),
		names:         copyMap(s.names),
		queryVar:      s.queryVar,
		varCreators:   s.varCreators,
	}
	vvalue := s.newVarValue(typ, name)
	key := Var(reflect.ValueOf(vvalue).Pointer())
	res.names[key] = name
	res.placeholders[key] = vvalue
	if s.queryVar == nil {
		res.queryVar = &key
	}
	return res, vvalue.(A)
}

func (s *State) copy() *State {
	if s == nil {
		return nil
	}
	return &State{
		substitutions: copyMap(s.substitutions),
		placeholders:  copyMap(s.placeholders),
		queryVar:      s.queryVar,
		names:         copyMap(s.names),
		varCreators:   s.varCreators,
	}
}

func (s *State) newVarValue(varType any, name string) any {
	for _, create := range s.varCreators {
		if val, ok := create(varType, name); ok {
			v := reflect.ValueOf(val)
			switch v.Kind() {
			case reflect.Ptr, reflect.Slice, reflect.Map:
				// call to Pointer only works for these types and otherwise panics
			default:
				panic("cannot make a variable that is not a pointer, slice or map " + v.Type().String())
			}
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
	_, ok := s.placeholders[key]
	return key, ok
}

func (s *State) lookupPlaceholderValue(key Var) any {
	placeholder, ok := s.placeholders[key]
	if !ok {
		panic(fmt.Sprintf("Var %v not found", key))
	}
	return placeholder
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
		ss = NewState()
		ss.substitutions = make(map[Var]any)
	} else {
		ss = s.copy()
	}
	ss.substitutions[key] = value
	return ss
}

func copyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
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
		return fmt.Sprintf("(() . %d)", len(s.placeholders))
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
	return fmt.Sprintf("(%s . %d)", strings.Join(ss, ", "), len(s.placeholders))
}
