package gomini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestLookup(t *testing.T) {
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
	zaxwyz = zaxwyz.AddKeyValue(zvar, a)
	zaxwyz = zaxwyz.AddKeyValue(xvar, w)
	zaxwyz = zaxwyz.AddKeyValue(yvar, z)

	xyvxwx := s.AddKeyValue(xvar, y)
	xyvxwx = xyvxwx.AddKeyValue(vvar, x)
	xyvxwx = xyvxwx.AddKeyValue(wvar, x)

	xbxywxe := s.AddKeyValue(xvar, ast.NewSymbol("b"))
	xbxywxe = xbxywxe.AddKeyValue(zvar, y)
	xbxywxe = xbxywxe.AddKeyValue(wvar, ast.NewList(x, ast.NewSymbol("e"), z))

	xezxyz := s.AddKeyValue(xvar, ast.NewSymbol("e"))
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
		startName := state.getName(start)
		t.Run("(walk "+startName+" "+s.String()+")", func(t *testing.T) {
			got := Lookup(start, state)
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
