package gomini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestHasCycle(t *testing.T) {
	s := NewState()
	var x, y *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	s, y = newVarWithName(s, "y", &ast.SExpr{})
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)

	sxy := s.Set(yvar, x)
	tests := []func() (Var, *ast.SExpr, *State, bool){
		tuple4(xvar, x, s, true),
		tuple4(xvar, y, s, false),
		tuple4(xvar, ast.NewList(y), sxy, true),
	}
	for _, test := range tests {
		v, w, s, want := test()
		vname := s.getName(v)
		t.Run("(hasCycle "+vname+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := hasCycle(v, w, s)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}

// exts either extends a substitution s with an association between the variable x and the value v ,
// or it returns nil if extending the substitution with the pair `(,x . ,v) would create a cycle.
func exts(x Var, v any, s *State) *State {
	if hasCycle(x, v, s) {
		return nil
	}
	return s.Set(x, v)
}

func TestExtsXA(t *testing.T) {
	got := NewState()
	var x *ast.SExpr
	got, x = newVarWithName(got, "x", &ast.SExpr{})
	xvar, _ := got.CastVar(x)
	want := got.copy()
	got = exts(xvar, ast.NewSymbol("a"), got)
	want = want.Set(xvar, ast.NewSymbol("a"))
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestExtsXX(t *testing.T) {
	got := NewState()
	var x *ast.SExpr
	got, x = newVarWithName(got, "x", &ast.SExpr{})
	xvar, _ := got.CastVar(x)
	res := exts(xvar, x, got)
	if res != nil {
		t.Fatalf("expected res == nil")
	}
}

func TestExtsXY(t *testing.T) {
	got := NewState()
	var x, y *ast.SExpr
	got, x = newVarWithName(got, "x", &ast.SExpr{})
	got, y = newVarWithName(got, "y", &ast.SExpr{})
	xvar, _ := got.CastVar(x)
	want := got.copy()
	got = exts(xvar, y, got)
	want = want.Set(xvar, y)
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestExtsXYZ(t *testing.T) {
	got := NewState()
	var x, y, z *ast.SExpr
	got, x = newVarWithName(got, "x", &ast.SExpr{})
	got, y = newVarWithName(got, "y", &ast.SExpr{})
	got, z = newVarWithName(got, "z", &ast.SExpr{})
	xvar, _ := got.CastVar(x)
	yvar, _ := got.CastVar(y)
	zvar, _ := got.CastVar(z)
	got = got.Set(zvar, x)
	got = got.Set(yvar, z)

	want := got.copy()
	got = exts(xvar, ast.NewSymbol("e"), got)
	want = want.Set(xvar, ast.NewSymbol("e"))
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
