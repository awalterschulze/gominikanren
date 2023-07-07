package gomini

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// AppendO is a goal that appends two lists into the third list.
func AppendO(l, t, out *ast.SExpr) Goal {
	return DisjO(
		ConjO(
			NullO(l),
			EqualO(t, out),
		),
		ExistO(func(a *ast.SExpr) Goal {
			return ExistO(func(d *ast.SExpr) Goal {
				return ExistO(func(res *ast.SExpr) Goal {
					return ConjO(
						ConsO(a, d, l),
						ConjO(
							ConsO(a, res, out),
							AppendO(d, t, res),
						),
					)
				})
			})
		}),
	)
}

// NullO is a goal that checks if a list is null.
func NullO(x *ast.SExpr) Goal {
	return EqualO(x, nil)
}

// ConsO is a goal that conses the first two expressions into the third.
func ConsO(a, d, p *ast.SExpr) Goal {
	return EqualO(ast.Cons(a, d), p)
}

// CarO is a goal where the second parameter is the head of the list in the first parameter.
func CarO(p, a *ast.SExpr) Goal {
	return ExistO(func(d *ast.SExpr) Goal {
		return EqualO(ast.Cons(a, d), p)
	})
}

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
	s := NewState(ast.CreateVar)
	sanys := Run(ctx, -1, s, func(q *ast.SExpr) Goal {
		return ExistO(func(x *ast.SExpr) Goal {
			return ExistO(func(y *ast.SExpr) Goal {
				return ConjO(
					EqualO(
						ast.Cons(x, ast.Cons(y, nil)),
						q,
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
	s := NewState(ast.CreateVar)
	sanys := Run(ctx, -1, s, func(q *ast.SExpr) Goal {
		return AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.Cons(ast.NewSymbol("b"), nil),
			q,
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
	s := NewState(ast.CreateVar)
	sanys := Run(ctx, -1, s, func(q *ast.SExpr) Goal {
		return AppendO(
			ast.Cons(ast.NewSymbol("a"), nil),
			ast.NewSymbol("b"),
			q,
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
	state := NewState()
	var y *ast.SExpr
	state, y = newVarWithName(state, "y", &ast.SExpr{})
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
		EqualO(ast.NewSymbol("#t"), y),
		EqualO(ast.NewSymbol("#f"), y),
	)
	ss := NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want := "(({,y: #t}, {,v1: (c o r n)} . 2))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
