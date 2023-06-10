package comini

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
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
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sexprs := comicro.Run(ctx, -1, func(q comicro.Var) comicro.Goal {
		return comicro.CallFresh(func(x comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(y comicro.Var) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(
						ast.Cons(x.SExpr(), ast.Cons(y.SExpr(), nil)),
						q.SExpr(),
					),
					AppendO(
						x.SExpr(),
						y.SExpr(),
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
		})
	})
	got := ast.NewList(ast.Sort(sexprs)...).String()
	want := "((() (cake & ice d t)) ((cake) (& ice d t)) ((cake &) (ice d t)) ((cake & ice) (d t)) ((cake & ice d) (t)) ((cake & ice d t) ()))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestAppendOSingleList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := comicro.Run(ctx, -1, func(q comicro.Var) comicro.Goal {
		return AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.Cons(ast.NewSymbol("b"), nil),
			q.SExpr(),
		)
	})
	got := ast.NewList(ss...).String()
	want := "((a b))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestAppendOSingleAtom(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := comicro.Run(ctx, -1, func(q comicro.Var) comicro.Goal {
		return AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.NewSymbol("b"),
			q.SExpr(),
		)
	})
	got := ast.NewList(ss...).String()
	want := "((a . b))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestCarO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	y := ast.NewVariable("y")
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
		comicro.EqualO(ast.NewSymbol("#t"), y),
		comicro.EqualO(ast.NewSymbol("#f"), y),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, comicro.EmptyState())
	got := ss.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", y.Atom.Var.Index), "y", -1)
	want := "(((,y . #t) (,v0 c o r n) . 1))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
