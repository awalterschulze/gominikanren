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
	s := NewEmptyState()
	var x, u, v, w, y, z Var
	s, u = s.NewVarWithName("u")
	s, v = s.NewVarWithName("v")
	s, w = s.NewVarWithName("w")
	s, x = s.NewVarWithName("x")
	s, y = s.NewVarWithName("y")
	s, z = s.NewVarWithName("z")
	a1 := ast.NewList(x.SExpr(), u.SExpr(), w.SExpr(), y.SExpr(), z.SExpr(), ast.NewList(ast.NewList(ast.NewSymbol("ice")), z.SExpr()))
	a2 := ast.Cons(y.SExpr(), ast.NewSymbol("corn"))
	a3 := ast.NewList(w.SExpr(), v.SExpr(), u.SExpr())
	e := ast.NewList(a1, a2, a3)
	s = s.AddKeyValue(x, e.Car().Cdr())
	s = s.AddKeyValue(y, e.Cdr().Car().Cdr())
	s = s.AddKeyValue(w, e.Cdr().Cdr().Car().Cdr())
	gote := reifyFromState(x, s).(*ast.SExpr)
	got := gote.String()
	want := "(_0 (_1 _0) corn _2 ((ice) _2))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNoReify(t *testing.T) {
	initial := NewEmptyState()
	var x Var
	initial, x = initial.NewVarWithName("x")
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		x,
	)
	g := Disj(e1, e2)
	states := RunGoal(context.Background(), 5, g)
	ss := make([]*ast.SExpr, len(states))
	strs := make([]string, len(states))
	for i, s := range states {
		ss[i] = reifyFromState(x, s).(*ast.SExpr)
		strs[i] = ss[i].String()
	}
	got := "(" + strings.Join(strs, " ") + ")"
	want1 := "(olive oil)"
	want2 := "(oil olive)"
	if got != want1 && got != want2 {
		t.Fatalf("got %s != want %s or %s", got, want1, want2)
	}
}
