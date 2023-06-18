package comini

import (
	"context"
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
	sanys := comicro.Run(ctx, -1, func(q comicro.Var) comicro.Goal {
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
	sexprs := fmap(sanys, func(a any) *ast.SExpr {
		return a.(*ast.SExpr)
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
	sanys := comicro.Run(ctx, -1, func(q comicro.Var) comicro.Goal {
		return AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.Cons(ast.NewSymbol("b"), nil),
			q.SExpr(),
		)
	})
	sexprs := fmap(sanys, func(a any) *ast.SExpr {
		return a.(*ast.SExpr)
	})
	got := ast.NewList(sexprs...).String()
	want := "((a b))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestAppendOSingleAtom(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sanys := comicro.Run(ctx, -1, func(q comicro.Var) comicro.Goal {
		return AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.NewSymbol("b"),
			q.SExpr(),
		)
	})
	sexprs := fmap(sanys, func(a any) *ast.SExpr {
		return a.(*ast.SExpr)
	})
	got := ast.NewList(sexprs...).String()
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
	state := comicro.NewEmptyState()
	var y comicro.Var
	state, y = state.NewVarWithName("y")
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
		comicro.EqualO(ast.NewSymbol("#t"), y.SExpr()),
		comicro.EqualO(ast.NewSymbol("#f"), y.SExpr()),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want := "(({,y: #t}, {,v1: (c o r n)} . 2))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
