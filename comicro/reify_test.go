package comicro

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
func TestReify(t *testing.T) {
	s := NewEmptyState().WithReifyNames(ast.ReifyName)
	var x, u, v, w, y, z *ast.SExpr
	s, u = NewVar(s, &ast.SExpr{})
	s, v = NewVar(s, &ast.SExpr{})
	s, w = NewVar(s, &ast.SExpr{})
	s, x = NewVar(s, &ast.SExpr{})
	s, y = NewVar(s, &ast.SExpr{})
	s, z = NewVar(s, &ast.SExpr{})
	xvar, _ := s.GetVar(x)
	yvar, _ := s.GetVar(y)
	wvar, _ := s.GetVar(w)
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

func TestNoReify(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	initial := NewEmptyState().WithReifyNames(ast.ReifyName)
	var x *ast.SExpr
	initial, x = NewVarWithName(initial, "x", &ast.SExpr{})
	xvar, _ := initial.GetVar(x)
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
