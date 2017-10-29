package micro

import (
	"strings"
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
	a1 := "(,x . (,u ,w ,y ,z ((ice) ,z)))"
	a2 := "(,y . corn)"
	a3 := "(,w . (,v ,u))"
	s := "(" + a1 + " " + a2 + " " + a3 + ")"
	e, err := sexpr.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	if !e.IsPair() {
		t.Fatalf("expected list")
	}
	gote := reifyStateVar(ast.NewVariable("x"), &State{Substitution: e.Pair})
	got := gote.String()
	want := "(_0 (_1 _0) corn _2 ((ice) _2))"
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
	states := RunGoal(5, g)
	ss := make([]*ast.SExpr, len(states))
	strs := make([]string, len(states))
	for i, s := range states {
		ss[i] = reifyStateVar(ast.NewVariable("x"), s)
		strs[i] = ss[i].String()
	}
	got := "(" + strings.Join(strs, " ") + ")"
	want := "(olive oil)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
