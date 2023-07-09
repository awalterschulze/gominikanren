package gomini

import (
	"context"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
scheme code:

	(let
		(
			(a1 `(,x . (,u ,w ,y ,z ((ice) ,z))))
			(a2 `(,y . corn))
			(a3 `(,w . (,v ,u)))
		)
		(let
			(
				(s `(,a1 ,a2 ,a3))
			)
			(
				(reify x) s
			)
		)
	)
*/
func TestRewrite1(t *testing.T) {
	s := NewState(ast.CreateVar)
	var x, u, v, w, y, z *ast.SExpr
	s, u = NewVar[*ast.SExpr](s)
	s, v = NewVar[*ast.SExpr](s)
	s, w = NewVar[*ast.SExpr](s)
	s, x = NewVar[*ast.SExpr](s)
	s, y = NewVar[*ast.SExpr](s)
	s, z = NewVar[*ast.SExpr](s)
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)
	wvar, _ := s.CastVar(w)
	a1 := ast.NewList(x, u, w, y, z, ast.NewList(ast.NewList(ast.NewSymbol("ice")), z))
	a2 := ast.Cons(y, ast.NewSymbol("corn"))
	a3 := ast.NewList(w, v, u)
	e := ast.NewList(a1, a2, a3)
	s = s.Set(xvar, e.Car().Cdr())
	s = s.Set(yvar, e.Cdr().Car().Cdr())
	s = s.Set(wvar, e.Cdr().Cdr().Car().Cdr())
	gote := rewrite(xvar, s).(*ast.SExpr)
	got := gote.String()
	want := "(v0 (v1 v0) corn v5 ((ice) v5))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNoRewrites(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	initial := NewState(ast.CreateVar)
	var x *ast.SExpr
	initial, x = newVarWithName[*ast.SExpr](initial, "x")
	xvar, _ := initial.CastVar(x)
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		x,
	)
	g := DisjO(e1, e2)
	stream := NewStreamForGoal(ctx, g, initial)
	states := take(ctx, 5, stream)
	ss := make([]*ast.SExpr, len(states))
	strs := make([]string, len(states))
	for i, s := range states {
		ss[i] = rewrite(xvar, s).(*ast.SExpr)
		strs[i] = ss[i].String()
	}
	got := "(" + strings.Join(strs, " ") + ")"
	want1 := "(olive oil)"
	want2 := "(oil olive)"
	if got != want1 && got != want2 {
		t.Fatalf("got %s != want %s or %s", got, want1, want2)
	}
}

func TestRewriteInternal(t *testing.T) {
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
		tuple3(wvar, xbxywxe, "(b e "+y.String()+")"),
		tuple3(yvar, xezxyz, "e"),
	}
	for _, test := range tests {
		start, state, want := test()
		startName := state.getName(start)
		t.Run("(rewrite "+startName+" "+state.String()+")", func(t *testing.T) {
			got := rewrite(start, state)
			if vgot, ok := state.CastVar(got); ok {
				got = state.getName(vgot)
			} else {
				got = got.(interface{ String() string }).String()
			}
			if want != got {
				t.Fatalf("got %T:%s want %T:%s", got, got, want, want)
			}
		})
	}
}
