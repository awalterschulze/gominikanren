package gomini

import (
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func newState(substs map[string]any) *State {
	s := NewState(ast.CreateVar)
	values := make(map[string]any)
	vars := make(map[string]Var)
	for k := range substs {
		var kvalue any
		s, kvalue = newVarWithName(s, k, &ast.SExpr{})
		values[k] = kvalue
		kvar, _ := s.CastVar(kvalue)
		vars[k] = kvar
	}
	for k, v := range substs {
		switch v := v.(type) {
		case string:
			if vvar, ok := values[v]; ok {
				s = s.Set(vars[k], vvar)
			} else {
				s = s.Set(vars[k], ast.NewSymbol(v))
			}
		case []string:
			sexprs := make([]*ast.SExpr, len(v))
			for i := 0; i < len(v); i++ {
				velem := v[i]
				if vvar, ok := values[velem]; ok {
					sexprs[i] = vvar.(*ast.SExpr)
				} else {
					sexprs[i] = ast.NewSymbol(velem)
				}
			}
			s = s.Set(vars[k], ast.NewList(sexprs...))
		}
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

func walks(start string, s *State) string {
	vars := make(map[string]Var)
	for k, v := range s.names {
		vars[v] = k
	}
	startvar := vars[start]
	return toString(s, walk(startvar, s))
}

func rewrites(start string, s *State) string {
	vars := make(map[string]Var)
	for k, v := range s.names {
		vars[v] = k
	}
	startvar := vars[start]
	return toString(s, rewrite(startvar, s))
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
