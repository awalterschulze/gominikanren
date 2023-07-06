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
	return &State{varCreators: varCreators, substitutions: make(map[Var]any), placeholders: make(map[Var]any), names: make(map[Var]string)}
}

func (s *State) GetQueryVar() *Var {
	return s.queryVar
}

func (s *State) Set(key Var, value any) *State {
	ss := s.copy()
	ss.substitutions[key] = value
	return ss
}

func (s *State) Get(key Var) (any, bool) {
	a, ok := s.substitutions[key]
	return a, ok
}

func (s *State) GetPlaceHolder(key Var) any {
	placeholder, ok := s.placeholders[key]
	if !ok {
		panic(fmt.Sprintf("Var %v not found", key))
	}
	return placeholder
}

func (s *State) CastVar(x any) (Var, bool) {
	if avar, ok := x.(Var); ok {
		return avar, true
	}
	if !isPointerValue(x) {
		return 0, false
	}
	key := Var(reflect.ValueOf(x).Pointer())
	_, ok := s.placeholders[key]
	return key, ok
}

type Var uintptr

func NewVar[A any](s *State, typ A) (*State, A) {
	return newVarWithName(s, "v"+strconv.Itoa(len(s.placeholders)), typ)
}

func newVarWithName[A any](s *State, name string, typ A) (*State, A) {
	res := &State{
		substitutions: s.substitutions,
		placeholders:  copyMap(s.placeholders),
		names:         copyMap(s.names),
		queryVar:      s.queryVar,
		varCreators:   s.varCreators,
	}
	vvalue := s.newVarValue(typ, name)
	vvar := Var(reflect.ValueOf(vvalue).Pointer())
	res.names[vvar] = name
	res.placeholders[vvar] = vvalue
	if s.queryVar == nil {
		res.queryVar = &vvar
	}
	return res, vvalue.(A)
}

func (s *State) newVarValue(varType any, name string) any {
	for _, create := range s.varCreators {
		if val, ok := create(varType, name); ok {
			if !isPointerValue(val) {
				panic(fmt.Sprintf("cannot make a variable that is not a pointer, slice or map: %#v", val))
			}
			return val
		}
	}
	return reflect.New(reflect.TypeOf(varType).Elem()).Interface()
}

func (s *State) copy() *State {
	return &State{
		substitutions: copyMap(s.substitutions),
		placeholders:  copyMap(s.placeholders),
		queryVar:      s.queryVar,
		names:         copyMap(s.names),
		varCreators:   s.varCreators,
	}
}

func isPointerValue(a any) bool {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		return true
	default:
		return false
	}
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
	name, ok := s.names[v]
	if ok {
		return name
	}
	return "v?"
}

// String returns a string representation of State.
func (s *State) String() string {
	if len(s.substitutions) == 0 {
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
