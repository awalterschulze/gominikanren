package gomini

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func getVars(s *State) map[string]Var {
	vars := make(map[string]Var)
	for k, v := range s.names {
		vars[v] = k
	}
	return vars
}

func addVars(s *state, v any) (*state, any) {
	switch v := v.(type) {
	case string:
		if _, ok := s.values[v]; !ok {
			fmt.Printf("adding %s\n", v)
			var vvar *ast.SExpr
			s.State, vvar = newVarWithName(s.State, v, &ast.SExpr{})
			s.values[v] = vvar
			return s, vvar
		} else {
			fmt.Printf("already added %s\n", v)
			return s, s.values[v]
		}
	case []string:
		xs := make([]*ast.SExpr, len(v))
		for i := 0; i < len(v); i++ {
			if _, ok := s.values[v[i]]; !ok {
				fmt.Printf("adding %s\n", v[i])
				var vvar *ast.SExpr
				s.State, vvar = newVarWithName(s.State, v[i], &ast.SExpr{})
				s.values[v[i]] = vvar
				xs[i] = vvar
			} else {
				fmt.Printf("already added %s\n", v[i])
				xs[i] = s.values[v[i]].(*ast.SExpr)
			}
		}
		return s, ast.NewList(xs...)
	}
	panic("unreachable")
}

type state struct {
	*State
	values map[string]any
}

func newState(substs map[string]any) *state {
	s := &state{NewState(ast.CreateVar), make(map[string]any)}
	for k := range substs {
		s, _ = addVars(s, k)
	}
	vars := getVars(s.State)
	m := make(map[Var]any)
	for k, v := range substs {
		var value any
		s, value = addVars(s, v)
		m[vars[k]] = value
	}
	for k, v := range m {
		s.State = s.State.Set(k, v)
	}
	return s
}

func toString(s *State, v any) string {
	if vvar, ok := s.CastVar(v); ok {
		return s.getName(vvar)
	}
	switch v := v.(type) {
	case *ast.SExpr:
		if v.IsPair() {
			ss := []string{}
			for v.IsPair() {
				car := v.Car()
				cdr := v.Cdr()
				v = cdr
				ss = append(ss, car.String())
			}
			return "[" + strings.Join(ss, " ") + "]"
		}
	}
	return v.(interface{ String() string }).String()
}

func walks(start string, s *state) string {
	vars := getVars(s.State)
	startvar := vars[start]
	return toString(s.State, walk(startvar, s.State))
}

func rewrites(start string, s *state) string {
	vars := getVars(s.State)
	startvar := vars[start]
	return toString(s.State, rewrite(startvar, s.State))
}

func unifys(x, y any, s *state) map[string]string {
	var xvalue, yvalue any
	s, xvalue = addVars(s, x)
	s, yvalue = addVars(s, y)
	s.State = unify(xvalue, yvalue, s.State)
	if s.State == nil {
		return nil
	}
	ms := make(map[string]string)
	for k, v := range s.substitutions {
		ms[toString(s.State, k)] = toString(s.State, v)
	}
	return ms
}
func TestUnifyDeep(t *testing.T) {
	substs := map[string]any{
		"a": "b",
	}
	s := newState(substs)
	got := unifys([]string{"a", "c"}, []string{"a", "d"}, s)
	want := map[string]string{
		"a": "b",
		"c": "d",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestUnifyNoCycle(t *testing.T) {
	substs := map[string]any{
		"a": "b",
	}
	s := newState(substs)
	got := unifys("b", "c", s)
	want := map[string]string{
		"a": "b",
		"b": "c",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestUnifyHasCycle1(t *testing.T) {
	substs := map[string]any{
		"x": []string{"y"},
	}
	s := newState(substs)
	got := unifys("y", "x", s)
	var want map[string]string = nil
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestUnifyHasCycle2(t *testing.T) {
	substs := map[string]any{
		"x": []string{"y"},
	}
	s := newState(substs)
	got := unifys("x", "y", s)
	var want map[string]string = nil
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestWalkZA(t *testing.T) {
	substs := map[string]any{
		"z": "a",
		"x": "w",
		"y": "z",
	}
	s := newState(substs)
	got := walks("z", s)
	want := "a"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteZA(t *testing.T) {
	substs := map[string]any{
		"z": "a",
		"x": "w",
		"y": "z",
	}
	s := newState(substs)
	got := rewrites("z", s)
	want := "a"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkYA(t *testing.T) {
	substs := map[string]any{
		"z": "a",
		"x": "w",
		"y": "z",
	}
	s := newState(substs)
	got := walks("y", s)
	want := "a"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteYA(t *testing.T) {
	substs := map[string]any{
		"z": "a",
		"x": "w",
		"y": "z",
	}
	s := newState(substs)
	got := rewrites("y", s)
	want := "a"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkXW(t *testing.T) {
	substs := map[string]any{
		"z": "a",
		"x": "w",
		"y": "z",
	}
	s := newState(substs)
	got := walks("x", s)
	want := "w"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteXW(t *testing.T) {
	substs := map[string]any{
		"z": "a",
		"x": "w",
		"y": "z",
	}
	s := newState(substs)
	got := rewrites("x", s)
	want := "w"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkXY(t *testing.T) {
	substs := map[string]any{
		"x": "y",
		"v": "x",
		"w": "x",
	}
	s := newState(substs)
	got := walks("x", s)
	want := "y"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteXY(t *testing.T) {
	substs := map[string]any{
		"x": "y",
		"v": "x",
		"w": "x",
	}
	s := newState(substs)
	got := rewrites("x", s)
	want := "y"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkVY(t *testing.T) {
	substs := map[string]any{
		"x": "y",
		"v": "x",
		"w": "x",
	}
	s := newState(substs)
	got := walks("v", s)
	want := "y"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteVY(t *testing.T) {
	substs := map[string]any{
		"x": "y",
		"v": "x",
		"w": "x",
	}
	s := newState(substs)
	got := rewrites("v", s)
	want := "y"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkWY(t *testing.T) {
	substs := map[string]any{
		"x": "y",
		"v": "x",
		"w": "x",
	}
	s := newState(substs)
	got := walks("w", s)
	want := "y"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteWY(t *testing.T) {
	substs := map[string]any{
		"x": "y",
		"v": "x",
		"w": "x",
	}
	s := newState(substs)
	got := rewrites("w", s)
	want := "y"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkYE(t *testing.T) {
	substs := map[string]any{
		"x": "e",
		"z": "x",
		"y": "z",
	}
	s := newState(substs)
	got := walks("y", s)
	want := "e"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkABCDE(t *testing.T) {
	substs := map[string]any{
		"a": "b",
		"b": "c",
		"c": "d",
		"e": []string{"a", "b", "c"},
	}
	got := walks("e", newState(substs))
	want := "[a b c]"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteABCDE(t *testing.T) {
	substs := map[string]any{
		"a": "b",
		"b": "c",
		"c": "d",
		"e": []string{"a", "b", "c"},
	}
	s := newState(substs)
	got := rewrites("e", s)
	want := "[d d d]"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteYE(t *testing.T) {
	substs := map[string]any{
		"x": "e",
		"z": "x",
		"y": "z",
	}
	s := newState(substs)
	got := rewrites("y", s)
	want := "e"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestWalkList(t *testing.T) {
	substs := map[string]any{
		"x": "b",
		"z": "y",
		"w": []string{"x", "e", "z"},
	}
	s := newState(substs)
	got := walks("w", s)
	want := "[x e z]"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestRewriteList(t *testing.T) {
	substs := map[string]any{
		"x": "b",
		"z": "y",
		"w": []string{"x", "e", "z"},
	}
	s := newState(substs)
	got := rewrites("w", s)
	want := "[b e y]"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}
