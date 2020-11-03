package micro

import (
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
	u := ast.NewVar("u", 0)
	v := ast.NewVar("v", 1)
	w := ast.NewVar("w", 2)
	x := ast.NewVar("x", 3)
	y := ast.NewVar("y", 4)
	z := ast.NewVar("z", 5)
	a1 := ast.NewList(x, u, w, y, z, ast.NewList(ast.NewList(ast.NewSymbol("ice")), z))
	a2 := ast.Cons(y, ast.NewSymbol("corn"))
	a3 := ast.NewList(w, v, u)
	e := ast.NewList(a1, a2, a3)
	ss := Substitutions{
		indexOf(x): e.Car().Cdr(),
		indexOf(y): e.Cdr().Car().Cdr(),
		indexOf(w): e.Cdr().Cdr().Car().Cdr(),
	}
	gote := ReifyIntVarFromState(indexOf(x))(&State{Substitutions: ss})
	got := gote.String()
	want := "(_0 (_1 _0) corn _2 ((ice) _2))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNoReify(t *testing.T) {
	x := ast.NewVariable("x")
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		x,
	)
	g := Disj(e1, e2)
	states := RunGoal(5, g)
	ss := make([]*ast.SExpr, len(states))
	strs := make([]string, len(states))
	r := ReifyIntVarFromState(indexOf(x))
	for i, s := range states {
		ss[i] = r(s)
		strs[i] = ss[i].String()
	}
	got := "(" + strings.Join(strs, " ") + ")"
	want := "(olive oil)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
