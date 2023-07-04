package gomini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOccurs(t *testing.T) {
	s := NewState()
	var x, y *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	s, y = newVarWithName(s, "y", &ast.SExpr{})
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)

	sxy := s.AddKeyValue(yvar, x)
	tests := []func() (Var, *ast.SExpr, *State, bool){
		tuple4(xvar, x, s, true),
		tuple4(xvar, y, s, false),
		tuple4(xvar, ast.NewList(y), sxy, true),
	}
	for _, test := range tests {
		v, w, s, want := test()
		vname := s.getName(v)
		t.Run("(occurs "+vname+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := occurs(v, w, s)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}

func TestExtsXA(t *testing.T) {
	got := NewState()
	var x *ast.SExpr
	got, x = newVarWithName(got, "x", &ast.SExpr{})
	xvar, _ := got.CastVar(x)
	want := got.copy()
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
	got := NewState()
	var x *ast.SExpr
	got, x = newVarWithName(got, "x", &ast.SExpr{})
	xvar, _ := got.CastVar(x)
	var gotok bool
	_, gotok = exts(xvar, x, got)
	if gotok {
		t.Fatalf("expected !ok")
	}
}

func TestExtsXY(t *testing.T) {
	got := NewState()
	var x, y *ast.SExpr
	got, x = newVarWithName(got, "x", &ast.SExpr{})
	got, y = newVarWithName(got, "y", &ast.SExpr{})
	xvar, _ := got.CastVar(x)
	want := got.copy()
	var gotok bool
	got, gotok = exts(xvar, y, got)
	if !gotok {
		t.Fatalf("expected ok")
	}
	want = want.AddKeyValue(xvar, y)
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
	got = got.AddKeyValue(zvar, x)
	got = got.AddKeyValue(yvar, z)

	want := got.copy()
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
