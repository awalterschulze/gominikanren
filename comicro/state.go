package comicro

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// State is a product of a list of substitutions and a variable counter.
type State struct {
	Substitutions map[Var]any
	Names         map[Var]string
	Counter       uint64
}

// NewEmptyState returns an empty state.
func NewEmptyState() *State {
	return &State{}
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

func (v Var) SExpr() *ast.SExpr {
	return ast.NewVar(fmt.Sprintf("v%d", v), uint64(v))
}

func (v Var) String() string {
	return v.SExpr().String()
}

func isvar(a any) bool {
	if IsNil(a) {
		return false
	}
	v := reflect.ValueOf(a)
	if v.Type() != varType {
		return false
	}
	return true
}

var varType = reflect.TypeOf(Var(0))

var sexprType = reflect.TypeOf(ast.SExpr{})
var atomType = reflect.TypeOf(ast.Atom{})

func isvarSExpr(s any) bool {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return false
	}
	if v.IsNil() {
		return false
	}
	v = v.Elem()
	k := v.Kind()
	switch k {
	case reflect.Struct:
		if v.Type() != sexprType {
			return false
		}
		atomValue := v.Field(1)
		if atomValue.IsNil() {
			return false
		}
		if atomValue.Elem().Type() != atomType {
			return false
		}
		varValue := atomValue.Elem().FieldByName("Var")
		if varValue.IsNil() {
			return false
		}
		return true
	}
	return false
}

func (s *State) NewVar() (*State, Var) {
	return s.NewVarWithName("v" + strconv.Itoa(int(s.Counter)))
}

func (s *State) NewVarWithName(name string) (*State, Var) {
	if s == nil {
		s = NewEmptyState()
	}
	v := Var(s.Counter)
	names := copyMap(s.Names)
	names[v] = name
	return &State{
		Substitutions: s.Substitutions,
		Counter:       s.Counter + 1,
		Names:         names,
	}, v
}

func (s *State) GetVar(a any) (Var, bool) {
	if isvar(a) {
		return a.(Var), true
	}
	if isvarSExpr(a) {
		return Var(a.(*ast.SExpr).Atom.Var.Index), true
	}
	return 0, false
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
	return &State{
		Substitutions: substitutions,
		Counter:       s.Counter,
		Names:         names,
	}
}

func copyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
