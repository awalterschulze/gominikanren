package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	defaultState := NewEmptyState()
	var a, v, w, x, y, z Var
	defaultState, v = defaultState.NewVarWithName("v")
	defaultState, w = defaultState.NewVarWithName("w")
	defaultState, x = defaultState.NewVarWithName("x")
	defaultState, y = defaultState.NewVarWithName("y")
	defaultState, z = defaultState.NewVarWithName("z")
	zaxwyz := defaultState.Copy()
	zaxwyz, a = zaxwyz.NewVarWithName("a")
	zaxwyz = zaxwyz.AddKeyValue(z, a)
	zaxwyz = zaxwyz.AddKeyValue(x, w)
	zaxwyz = zaxwyz.AddKeyValue(y, z)
	xyvxwx := defaultState.Copy()
	xyvxwx = xyvxwx.AddKeyValue(x, y)
	xyvxwx = xyvxwx.AddKeyValue(v, x)
	xyvxwx = xyvxwx.AddKeyValue(w, x)
	xbxywxe := defaultState.Copy()
	xbxywxe = xbxywxe.AddKeyValue(x, ast.NewSymbol("b"))
	xbxywxe = xbxywxe.AddKeyValue(z, y)
	xbxywxe = xbxywxe.AddKeyValue(w, ast.NewList(x.SExpr(), ast.NewSymbol("e"), z.SExpr()))
	xezxyz := defaultState.Copy()
	xezxyz = xezxyz.AddKeyValue(x, ast.NewSymbol("e"))
	xezxyz = xezxyz.AddKeyValue(z, x)
	xezxyz = xezxyz.AddKeyValue(y, z)
	tests := []func() (Var, *State, string){
		tuple3(z, zaxwyz, "a"),
		tuple3(y, zaxwyz, "a"),
		tuple3(x, zaxwyz, "w"),
		tuple3(x, xyvxwx, "y"),
		tuple3(v, xyvxwx, "y"),
		tuple3(w, xyvxwx, "y"),
		tuple3(w, xbxywxe, "("+x.String()+" e "+z.String()+")"),
		tuple3(y, xezxyz, "e"),
	}
	for _, test := range tests {
		v, s, want := test()
		t.Run("(walk "+v.String()+" "+s.String()+")", func(t *testing.T) {
			gotany := Lookup(v, s)
			got := Lookup(v, s).(interface{ String() string }).String()
			if vgot, ok := gotany.(Var); ok {
				got = s.GetName(vgot)
			}
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestReplaceAll(t *testing.T) {
	s := NewEmptyState()
	var a, v, w, x, y, z Var
	s, v = s.NewVarWithName("v")
	s, w = s.NewVarWithName("w")
	s, x = s.NewVarWithName("x")
	s, y = s.NewVarWithName("y")
	s, z = s.NewVarWithName("z")
	zaxwyz := s.Copy()
	zaxwyz, a = zaxwyz.NewVarWithName("a")
	zaxwyz = zaxwyz.AddKeyValue(z, a)
	zaxwyz = zaxwyz.AddKeyValue(x, w)
	zaxwyz = zaxwyz.AddKeyValue(y, z)
	xyvxwx := s.Copy()
	xyvxwx = xyvxwx.AddKeyValue(x, y)
	xyvxwx = xyvxwx.AddKeyValue(v, x)
	xyvxwx = xyvxwx.AddKeyValue(w, x)
	xbxywxe := s.Copy()
	xbxywxe = xbxywxe.AddKeyValue(x, ast.NewSymbol("b"))
	xbxywxe = xbxywxe.AddKeyValue(z, y)
	xbxywxe = xbxywxe.AddKeyValue(w, ast.NewList(x.SExpr(), ast.NewSymbol("e"), z.SExpr()))
	xezxyz := s.Copy()
	xezxyz = xezxyz.AddKeyValue(x, ast.NewSymbol("e"))
	xezxyz = xezxyz.AddKeyValue(z, x)
	xezxyz = xezxyz.AddKeyValue(y, z)
	tests := []func() (Var, *State, string){
		tuple3(z, zaxwyz, "a"),
		tuple3(y, zaxwyz, "a"),
		tuple3(x, zaxwyz, "w"),
		tuple3(x, xyvxwx, "y"),
		tuple3(v, xyvxwx, "y"),
		tuple3(w, xyvxwx, "y"),
		tuple3(w, xbxywxe, "(b e "+y.String()+")"),
		tuple3(y, xezxyz, "e"),
	}
	for _, test := range tests {
		start, state, want := test()
		t.Run("(walk "+start.String()+" "+state.String()+")", func(t *testing.T) {
			got := ReplaceAll(start, state)
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
