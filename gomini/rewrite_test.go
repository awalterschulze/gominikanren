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
	s := NewEmptyState().WithVarCreators(ast.CreateVar)
	var x, u, v, w, y, z *ast.SExpr
	s, u = NewVar(s, &ast.SExpr{})
	s, v = NewVar(s, &ast.SExpr{})
	s, w = NewVar(s, &ast.SExpr{})
	s, x = NewVar(s, &ast.SExpr{})
	s, y = NewVar(s, &ast.SExpr{})
	s, z = NewVar(s, &ast.SExpr{})
	xvar, _ := s.castVar(x)
	yvar, _ := s.castVar(y)
	wvar, _ := s.castVar(w)
	a1 := ast.NewList(x, u, w, y, z, ast.NewList(ast.NewList(ast.NewSymbol("ice")), z))
	a2 := ast.Cons(y, ast.NewSymbol("corn"))
	a3 := ast.NewList(w, v, u)
	e := ast.NewList(a1, a2, a3)
	s = s.AddKeyValue(xvar, e.Car().Cdr())
	s = s.AddKeyValue(yvar, e.Cdr().Car().Cdr())
	s = s.AddKeyValue(wvar, e.Cdr().Cdr().Car().Cdr())
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
	initial := NewEmptyState().WithVarCreators(ast.CreateVar)
	var x *ast.SExpr
	initial, x = NewVarWithName(initial, "x", &ast.SExpr{})
	xvar, _ := initial.castVar(x)
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		x,
	)
	g := Disj(e1, e2)
	stream := NewStreamForGoal(ctx, g, initial)
	states := Take(ctx, 5, stream)
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
	s := NewEmptyState()
	var a, v, w, x, y, z *ast.SExpr
	s, v = NewVarWithName(s, "v", &ast.SExpr{})
	s, w = NewVarWithName(s, "w", &ast.SExpr{})
	s, x = NewVarWithName(s, "x", &ast.SExpr{})
	s, y = NewVarWithName(s, "y", &ast.SExpr{})
	s, z = NewVarWithName(s, "z", &ast.SExpr{})
	xvar, _ := s.castVar(x)
	yvar, _ := s.castVar(y)
	zvar, _ := s.castVar(z)
	wvar, _ := s.castVar(w)
	vvar, _ := s.castVar(v)
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
		startName := state.getName(start)
		t.Run("(walk "+startName+" "+state.String()+")", func(t *testing.T) {
			got := rewrite(start, state)
			if vgot, ok := state.castVar(got); ok {
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
