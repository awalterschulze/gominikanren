package gomini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	s := NewState()
	var v, w, x, y, z *ast.SExpr
	s, v = newVarWithName[*ast.SExpr](s, "v")
	s, w = newVarWithName[*ast.SExpr](s, "w")
	s, x = newVarWithName[*ast.SExpr](s, "x")
	s, y = newVarWithName[*ast.SExpr](s, "y")
	s, z = newVarWithName[*ast.SExpr](s, "z")
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)
	zvar, _ := s.CastVar(z)
	wvar, _ := s.CastVar(w)
	vvar, _ := s.CastVar(v)

	zaxwyz, a := newVarWithName[*ast.SExpr](s, "a")
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
