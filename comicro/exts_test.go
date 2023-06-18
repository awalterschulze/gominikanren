package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOccurs(t *testing.T) {
	s := NewEmptyState()
	var x, y Var
	s, x = s.NewVarWithName("x")
	s, y = s.NewVarWithName("y")
	sxy := s.Copy()
	sxy = sxy.AddKeyValue(y, x)
	tests := []func() (Var, *ast.SExpr, *State, bool){
		tuple4(x, x.SExpr(), s, true),
		tuple4(x, y.SExpr(), s, false),
		tuple4(x, ast.NewList(y.SExpr()), sxy, true),
	}
	for _, test := range tests {
		v, w, s, want := test()
		t.Run("(occurs "+v.String()+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := occurs(v, w, s)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}

func TestExtsXA(t *testing.T) {
	got := NewEmptyState()
	var xvar Var
	got, xvar = got.NewVarWithName("x")
	want := got.Copy()
	var gotok bool
	got, gotok = exts(xvar, ast.NewSymbol("a"), got)
	if !gotok {
		t.Fatalf("expected ok")
	}
	want = want.AddKeyValue(xvar, ast.NewSymbol("a"))
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestExtsXX(t *testing.T) {
	got := NewEmptyState()
	var xvar Var
	got, xvar = got.NewVarWithName("x")
	var gotok bool
	_, gotok = exts(xvar, xvar, got)
	if gotok {
		t.Fatalf("expected !ok")
	}
}

func TestExtsXY(t *testing.T) {
	got := NewEmptyState()
	var xvar, yvar Var
	got, xvar = got.NewVarWithName("x")
	got, yvar = got.NewVarWithName("y")
	want := got.Copy()
	var gotok bool
	got, gotok = exts(xvar, yvar, got)
	if !gotok {
		t.Fatalf("expected ok")
	}
	want = want.AddKeyValue(xvar, yvar)
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestExtsXYZ(t *testing.T) {
	got := NewEmptyState()
	var xvar, yvar, zvar Var
	got, xvar = got.NewVarWithName("x")
	got, yvar = got.NewVarWithName("y")
	got, zvar = got.NewVarWithName("z")
	got = got.AddKeyValue(zvar, xvar)
	got = got.AddKeyValue(yvar, zvar)
	want := got.Copy()
	var gotok bool
	got, gotok = exts(xvar, ast.NewSymbol("e"), got)
	if !gotok {
		t.Fatalf("expected ok")
	}
	want = want.AddKeyValue(xvar, ast.NewSymbol("e"))
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
