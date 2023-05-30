package comini

import (
	"github.com/awalterschulze/gominikanren/micro"
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
func MemberO(x, y *ast.SExpr) micro.Goal {
	return micro.Zzz(micro.CallFresh(func(a *ast.SExpr) micro.Goal {
		return micro.CallFresh(func(d *ast.SExpr) micro.Goal {
			return micro.Conj(
				micro.EqualO(y, ast.Cons(a, d)),
				micro.Disj(
					micro.EqualO(x, a),
					MemberO(x, d),
				),
			)
		})
	}))
}

// MemberOUnrolled is the partial application version of MemberO
// taking the second argument (the list) first
func MemberOUnrolled(y *ast.SExpr) func(*ast.SExpr) micro.Goal {
	goals := []func(*ast.SExpr) micro.Goal{}
	for {
		if y == nil {
			break
		}
		car, cdr := y.Car(), y.Cdr()
		goals = append(goals, func(x *ast.SExpr) micro.Goal {
			return micro.EqualO(x, car)
		})
		y = cdr
	}
	n := len(goals)
	return func(x *ast.SExpr) micro.Goal {
		gs := make([]micro.Goal, n)
		for i, g := range goals {
			gs[i] = g(x)
		}
		return DisjPlus(gs...)
	}
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
func MapO(f func(*ast.SExpr, *ast.SExpr) micro.Goal, x, y *ast.SExpr) micro.Goal {
	return Conde(
		[]micro.Goal{
			micro.EqualO(x, nil),
			micro.EqualO(y, nil),
		},
		[]micro.Goal{
			// really missing Fresh (xa xd ya yd) here
			micro.CallFresh(func(xa *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(xd *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(ya *ast.SExpr) micro.Goal {
						return micro.CallFresh(func(yd *ast.SExpr) micro.Goal {
							return ConjPlus(
								micro.EqualO(x, ast.Cons(xa, xd)),
								micro.EqualO(y, ast.Cons(ya, yd)),
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

// MapOUnrolled is partial application of MapO with y already filled in
func MapOUnrolled(list *ast.SExpr) func(func(*ast.SExpr, *ast.SExpr) micro.Goal, *ast.SExpr) micro.Goal {
	goals := []func(func(*ast.SExpr, *ast.SExpr) micro.Goal, *ast.SExpr) micro.Goal{}
	for {
		if list == nil {
			break
		}
		car, cdr := list.Car(), list.Cdr()
		goals = append(goals, func(f func(*ast.SExpr, *ast.SExpr) micro.Goal, x *ast.SExpr) micro.Goal {
			return f(x, car)
		})
		list = cdr
	}
	n := len(goals)
	return func(f func(*ast.SExpr, *ast.SExpr) micro.Goal, x *ast.SExpr) micro.Goal {
		gs := make([]micro.Goal, n+1)
		vars := make([]*ast.SExpr, n)
		var cons *ast.SExpr = nil
		for i := 0; i < n; i++ {
			v := ast.NewVariable("")
			vars[i] = v
			gs[n-i] = goals[n-i-1](f, v)
			cons = ast.Cons(v, cons)
		}
		gs[0] = micro.EqualO(x, cons)
		return ConjPlus(gs...)
	}
}

// MapODoubleUnrolled is partial application of MapO with f and y already filled in
func MapODoubleUnrolled(f func(*ast.SExpr, *ast.SExpr) micro.Goal, list *ast.SExpr) func(*ast.SExpr) micro.Goal {
	goals := []func(*ast.SExpr) micro.Goal{}
	for {
		if list == nil {
			break
		}
		car, cdr := list.Car(), list.Cdr()
		goals = append(goals, func(x *ast.SExpr) micro.Goal {
			return f(x, car)
		})
		list = cdr
	}
	n := len(goals)
	return func(x *ast.SExpr) micro.Goal {
		gs := make([]micro.Goal, n+1)
		vars := make([]*ast.SExpr, n)
		var cons *ast.SExpr = nil
		for i := 0; i < n; i++ {
			v := ast.NewVariable("")
			vars[i] = v
			gs[n-i] = goals[n-i-1](v)
			cons = ast.Cons(v, cons)
		}
		gs[0] = micro.EqualO(x, cons)
		return ConjPlus(gs...)
	}
}
