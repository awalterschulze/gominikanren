package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
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
	a1 := "`(,x . (,u ,w ,y ,z ((ice) ,z)))"
	a2 := "`(,y . corn)"
	a3 := "`(,w . (,v ,u))"
	s := "(" + a1 + " " + a2 + " " + a3 + ")"
	e, err := sexpr.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	if e.List == nil {
		t.Fatalf("expected list")
	}
	x := ast.NewVariable("x")
	gote := reify(x)(e.List)
	got := gote.String()
	want := "`(_0 `(_1 _0) corn _2 `(`(ice) _2))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNoReify(t *testing.T) {
	e1 := EqualO(
		ast.NewSymbol("olive"),
		ast.NewVariable("x"),
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		ast.NewVariable("x"),
	)
	g := DisjointO(e1, e2)
	es := deriveFmapR(reify(ast.NewVariable("x")), RunGoal(5, g))
	l := ast.NewList(false, es...)
	got := l.String()
	want := "(olive oil)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
