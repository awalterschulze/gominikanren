package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	s := NewEmptyState()
	var a, v, w, x, y, z *ast.SExpr
	s, v = NewVarWithName(s, "v", &ast.SExpr{})
	s, w = NewVarWithName(s, "w", &ast.SExpr{})
	s, x = NewVarWithName(s, "x", &ast.SExpr{})
	s, y = NewVarWithName(s, "y", &ast.SExpr{})
	s, z = NewVarWithName(s, "z", &ast.SExpr{})
	xvar, _ := s.GetVar(x)
	yvar, _ := s.GetVar(y)
	zvar, _ := s.GetVar(z)
	wvar, _ := s.GetVar(w)
	vvar, _ := s.GetVar(v)
	zaxwyz := s.Copy()
	zaxwyz, a = NewVarWithName(zaxwyz, "a", &ast.SExpr{})
	zaxwyz = zaxwyz.AddKeyValue(zvar, a)
	zaxwyz = zaxwyz.AddKeyValue(xvar, w)
	zaxwyz = zaxwyz.AddKeyValue(yvar, z)
	xyvxwx := s.Copy()
	xyvxwx = xyvxwx.AddKeyValue(xvar, y)
	xyvxwx = xyvxwx.AddKeyValue(vvar, x)
	xyvxwx = xyvxwx.AddKeyValue(wvar, x)
	xbxywxe := s.Copy()
	xbxywxe = xbxywxe.AddKeyValue(xvar, ast.NewSymbol("b"))
	xbxywxe = xbxywxe.AddKeyValue(zvar, y)
	xbxywxe = xbxywxe.AddKeyValue(wvar, ast.NewList(x, ast.NewSymbol("e"), z))
	xezxyz := s.Copy()
	xezxyz = xezxyz.AddKeyValue(xvar, ast.NewSymbol("e"))
	xezxyz = xezxyz.AddKeyValue(zvar, x)
	xezxyz = xezxyz.AddKeyValue(yvar, z)
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
		startName := state.GetName(start)
		t.Run("(walk "+startName+" "+s.String()+")", func(t *testing.T) {
			got := Lookup(start, state)
			if vgot, ok := state.GetVar(got); ok {
				got = state.GetName(vgot)
			} else {
				got = got.(interface{ String() string }).String()
			}
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestReplaceAll(t *testing.T) {
	s := NewEmptyState()
	var a, v, w, x, y, z *ast.SExpr
	s, v = NewVarWithName(s, "v", &ast.SExpr{})
	s, w = NewVarWithName(s, "w", &ast.SExpr{})
	s, x = NewVarWithName(s, "x", &ast.SExpr{})
	s, y = NewVarWithName(s, "y", &ast.SExpr{})
	s, z = NewVarWithName(s, "z", &ast.SExpr{})
	xvar, _ := s.GetVar(x)
	yvar, _ := s.GetVar(y)
	zvar, _ := s.GetVar(z)
	wvar, _ := s.GetVar(w)
	vvar, _ := s.GetVar(v)
	zaxwyz := s.Copy()
	zaxwyz, a = NewVarWithName(zaxwyz, "a", &ast.SExpr{})
	zaxwyz = zaxwyz.AddKeyValue(zvar, a)
	zaxwyz = zaxwyz.AddKeyValue(xvar, w)
	zaxwyz = zaxwyz.AddKeyValue(yvar, z)
	xyvxwx := s.Copy()
	xyvxwx = xyvxwx.AddKeyValue(xvar, y)
	xyvxwx = xyvxwx.AddKeyValue(vvar, x)
	xyvxwx = xyvxwx.AddKeyValue(wvar, x)
	xbxywxe := s.Copy()
	xbxywxe = xbxywxe.AddKeyValue(xvar, ast.NewSymbol("b"))
	xbxywxe = xbxywxe.AddKeyValue(zvar, y)
	xbxywxe = xbxywxe.AddKeyValue(wvar, ast.NewList(x, ast.NewSymbol("e"), z))
	xezxyz := s.Copy()
	xezxyz = xezxyz.AddKeyValue(xvar, ast.NewSymbol("e"))
	xezxyz = xezxyz.AddKeyValue(zvar, x)
	xezxyz = xezxyz.AddKeyValue(yvar, z)
	tests := []func() (Var, *State, string){
		tuple3(zvar, zaxwyz, "a"),
		tuple3(yvar, zaxwyz, "a"),
		tuple3(xvar, zaxwyz, "w"),
		tuple3(xvar, xyvxwx, "y"),
		tuple3(vvar, xyvxwx, "y"),
		tuple3(wvar, xyvxwx, "y"),
		tuple3(wvar, xbxywxe, "(b e "+y.String()+")"),
		tuple3(yvar, xezxyz, "e"),
	}
	for _, test := range tests {
		start, state, want := test()
		startName := state.GetName(start)
		t.Run("(walk "+startName+" "+state.String()+")", func(t *testing.T) {
			got := rewrite(start, state)
			if vgot, ok := state.GetVar(got); ok {
				got = state.GetName(vgot)
			} else {
				got = got.(interface{ String() string }).String()
			}
			if want != got {
				t.Fatalf("got %T:%s want %T:%s", got, got, want, want)
			}
		})
	}
}
