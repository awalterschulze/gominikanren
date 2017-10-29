package mini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(let
	(
		(q (var 'q))
	)
	(map
		(reify q)
		(run-goal #f
			(call/fresh 'x
				(lambda (x)
					(call/fresh 'y
						(lambda (y)
							(conj
								(== `(,x ,y) q)
								(appendo x y `(cake & ice d t))
							)
						)
					)
				)
			)
		)
	)
)
*/
// results in all the combinations of two lists that when appended will result in (cake & ice d t)
func TestAppendOAllCombinations(t *testing.T) {
	states := micro.RunGoal(
		-1,
		micro.CallFresh(func(x *ast.SExpr) micro.Goal {
			return micro.CallFresh(func(y *ast.SExpr) micro.Goal {
				return micro.ConjunctionO(
					micro.EqualO(
						ast.Cons(x, ast.Cons(y, nil)),
						ast.NewVariable("q"),
					),
					AppendO(
						x,
						y,
						ast.NewList(
							ast.NewSymbol("cake"),
							ast.NewSymbol("&"),
							ast.NewSymbol("ice"),
							ast.NewSymbol("d"),
							ast.NewSymbol("t"),
						),
					),
				)
			})
		}),
	)
	sexprs := micro.Reify(ast.NewVariable("q"), states)
	got := ast.NewList(sexprs...).String()
	want := "((() (cake & ice d t)) ((cake) (& ice d t)) ((cake &) (ice d t)) ((cake & ice) (d t)) ((cake & ice d) (t)) ((cake & ice d t) ()))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestAppendOSingleList(t *testing.T) {
	subs := micro.RunGoal(
		-1,
		AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.Cons(ast.NewSymbol("b"), nil),
			ast.NewVariable("q"),
		),
	)
	ss := micro.Reify(ast.NewVariable("q"), subs)
	got := ast.NewList(ss...).String()
	want := "((a b))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestAppendOSingleAtom(t *testing.T) {
	subs := micro.RunGoal(
		-1,
		AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.NewSymbol("b"),
			ast.NewVariable("q"),
		),
	)
	ss := micro.Reify(ast.NewVariable("q"), subs)
	got := ast.NewList(ss...).String()
	want := "((a . b))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestCarO(t *testing.T) {
	ifte := IfThenElseO(
		CarO(
			ast.NewList(
				ast.NewSymbol("a"),
				ast.NewSymbol("c"),
				ast.NewSymbol("o"),
				ast.NewSymbol("r"),
				ast.NewSymbol("n"),
			),
			ast.NewSymbol("a"),
		),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
	)
	ss := ifte(micro.EmptyState())
	got := ss.String()
	want := "(((,y . #t) (,v0 c o r n) . 1))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
