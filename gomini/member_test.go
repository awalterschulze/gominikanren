package gomini

import (
	"context"
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// defines two candidate functions for unrolling
// the idea is to do partial application and 'unroll' recursive functions
// into a disjunction/conjunction of goals, so that they can be
// executed in parallel in order to speed them up
// unrolled versions are faster even without concurrency

/*
MemberO checks membership in a list
we make a function that does partial application
i.e. given y = (1 2 3), (membero x y) unrolls to
(conde [(== x 1)][(== x 2)][(== x 3)])
(and of course then we can use a parallel version of conde)
unrolling when x is bound but y is unbound does not make sense
but we can optimise if we detect y is bound

(define (membero x y)

	(fresh (a d)
	(== y `(,a . ,d))
	(conde
	  [(== x a)]
	  [(membero x d)]
	  )))
*/
func MemberO(x, y *ast.SExpr) Goal {
	return ExistO(func(a *ast.SExpr) Goal {
		return ExistO(func(d *ast.SExpr) Goal {
			return ConjO(
				EqualO(y, ast.Cons(a, d)),
				DisjO(
					EqualO(x, a),
					MemberO(x, d),
				),
			)
		})
	})
}

/*
MapO defines a relation between inputs and outputs of a function f
inputs and outputs are given as equal-length lists
we make a function that inspects and unrolls this
i.e. given x = (1 2 3), (mapo f x y) unrolls to the following
(fresh (y1 y2 y3) (f 1 y1) (f 2 y2) (f 3 y3) (== x `(,y1 ,y2 ,y3)))
of which each subgoal can be evaluated in parallel
here unrolling makes sense whether x or y is bound
we can detect either and unroll before applying search

(define (mapo f x y)

	(conde
	  [(== x '()) (== y '())]
	  [(fresh (xa xd ya yd)
	      (== `(,xa . ,xd) x)
	      (== `(,ya . ,yd) y)
	      (f xa ya)
	      (mapo f xd yd)
	    )]
	))
*/
func MapO(f func(*ast.SExpr, *ast.SExpr) Goal, x, y *ast.SExpr) Goal {
	return DisjO(
		ConjO(
			EqualO(x, nil),
			EqualO(y, nil),
		),
		ConjO(
			ExistO(func(xa *ast.SExpr) Goal {
				return ExistO(func(xd *ast.SExpr) Goal {
					return ExistO(func(ya *ast.SExpr) Goal {
						return ExistO(func(yd *ast.SExpr) Goal {
							return ConjO(
								EqualO(x, ast.Cons(xa, xd)),
								EqualO(y, ast.Cons(ya, yd)),
								f(xa, ya),
								MapO(f, xd, yd),
							)
						})
					})
				})
			}),
		),
	)
}

func TestMemberO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(ast.CreateVar)
	list := ast.Cons(ast.NewInt(0), ast.Cons(ast.NewInt(1), ast.Cons(ast.NewInt(2), nil)))
	sanys := RunTake(ctx, -1, s, func(q *ast.SExpr) Goal {
		return MemberO(
			q,
			list,
		)
	})
	sexprs := ast.Sort(fmap(sanys, func(a any) *ast.SExpr {
		return a.(*ast.SExpr)
	}))
	if len(sexprs) != 3 {
		t.Fatalf("expected len %d, but got len %d instead", 3, len(sexprs))
	}
	for i, sexpr := range sexprs {
		if *sexpr.Atom.Int != int64(i) {
			t.Fatalf("expected %d, but got %d instead", i, *sexpr.Atom.Int)
		}
	}
}

func TestMapO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(ast.CreateVar)
	list := ast.Cons(ast.NewInt(0), ast.Cons(ast.NewInt(1), ast.Cons(ast.NewInt(2), nil)))
	sexprs := RunTake(ctx, -1, s, func(q *ast.SExpr) Goal {
		return MapO(
			func(x, y *ast.SExpr) Goal { return EqualO(x, y) },
			q,
			list,
		)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	if !reflect.DeepEqual(sexprs[0], list) {
		t.Fatalf("expected output matches input list, but got %#v", sexprs[0])
	}
}
