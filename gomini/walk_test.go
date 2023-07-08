package gomini

import (
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	s := NewState()
	var v, w, x, y, z *ast.SExpr
	s, v = newVarWithName(s, "v", &ast.SExpr{})
	s, w = newVarWithName(s, "w", &ast.SExpr{})
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	s, y = newVarWithName(s, "y", &ast.SExpr{})
	s, z = newVarWithName(s, "z", &ast.SExpr{})
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)
	zvar, _ := s.CastVar(z)
	wvar, _ := s.CastVar(w)
	vvar, _ := s.CastVar(v)

	zaxwyz, a := newVarWithName(s, "a", &ast.SExpr{})
	zaxwyz = zaxwyz.Set(zvar, a)
	zaxwyz = zaxwyz.Set(xvar, w)
	zaxwyz = zaxwyz.Set(yvar, z)

	xyvxwx := s.Set(xvar, y)
	xyvxwx = xyvxwx.Set(vvar, x)
	xyvxwx = xyvxwx.Set(wvar, x)

	xbxywxe := s.Set(xvar, ast.NewSymbol("b"))
	xbxywxe = xbxywxe.Set(zvar, y)
	xbxywxe = xbxywxe.Set(wvar, ast.NewList(x, ast.NewSymbol("e"), z))

	xezxyz := s.Set(xvar, ast.NewSymbol("e"))
	xezxyz = xezxyz.Set(zvar, x)
	xezxyz = xezxyz.Set(yvar, z)
	tests := []func() (Var, *State, string){
		tuple3(zvar, zaxwyz, "a"),
		tuple3(yvar, zaxwyz, "a"),
		tuple3(xvar, zaxwyz, "w"),
		tuple3(xvar, xyvxwx, "y"),
		tuple3(vvar, xyvxwx, "y"),
		tuple3(wvar, xyvxwx, "y"),
		tuple3(wvar, xbxywxe, "("+x.String()+" e "+z.String()+")"),
		tuple3(yvar, xezxyz, "e"),
	}
	for _, test := range tests {
		start, state, want := test()
		startName := state.getName(start)
		t.Run("(walk "+startName+" "+s.String()+")", func(t *testing.T) {
			got := walk(start, state)
			if vgot, ok := state.CastVar(got); ok {
				got = state.getName(vgot)
			} else {
				got = got.(interface{ String() string }).String()
			}
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func newState(substs map[string]any) *State {
	s := NewState(ast.CreateVar)
	values := make(map[string]any)
	vars := make(map[string]Var)
	for k, _ := range substs {
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

// Good Example
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
