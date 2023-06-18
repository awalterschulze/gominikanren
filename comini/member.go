package comini

import (
	"github.com/awalterschulze/gominikanren/comicro"
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
func MemberO(x, y *ast.SExpr) comicro.Goal {
	return comicro.CallFresh(&ast.SExpr{}, func(a *ast.SExpr) comicro.Goal {
		return comicro.CallFresh(&ast.SExpr{}, func(d *ast.SExpr) comicro.Goal {
			return comicro.Conj(
				comicro.EqualO(y, ast.Cons(a, d)),
				comicro.Disj(
					comicro.EqualO(x, a),
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
func MapO(f func(*ast.SExpr, *ast.SExpr) comicro.Goal, x, y *ast.SExpr) comicro.Goal {
	return Conde(
		[]comicro.Goal{
			comicro.EqualO(x, nil),
			comicro.EqualO(y, nil),
		},
		[]comicro.Goal{
			comicro.CallFresh(&ast.SExpr{}, func(xa *ast.SExpr) comicro.Goal {
				return comicro.CallFresh(&ast.SExpr{}, func(xd *ast.SExpr) comicro.Goal {
					return comicro.CallFresh(&ast.SExpr{}, func(ya *ast.SExpr) comicro.Goal {
						return comicro.CallFresh(&ast.SExpr{}, func(yd *ast.SExpr) comicro.Goal {
							return Conjs(
								comicro.EqualO(x, ast.Cons(xa, xd)),
								comicro.EqualO(y, ast.Cons(ya, yd)),
								f(xa, ya),
								MapO(f, xd, yd),
							)
						})
					})
				})
			}),
		},
	)
}
