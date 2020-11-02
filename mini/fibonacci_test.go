package mini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// peano numbers and extralogical convenience functions
var zero = ast.NewInt(0)
var one = makenat(1)

/*
   (define (succ x y)
       (== y `(,x)))
*/
func succ(prev, next *ast.SExpr) micro.Goal {
	return micro.EqualO(next, ast.Cons(prev, nil))
}

func makenat(n int) *ast.SExpr {
	if n == 0 {
		return zero
	}
	return ast.Cons(makenat(n-1), nil)
}

func parsenat(x *ast.SExpr) int {
	if x == zero {
		return 0
	}
	return parsenat(x.Car()) + 1
}

/*
   ;; x + y = z
   (define (nat-plus x y z)
   (conde
       [(== x 0) (== y z)]
       [(fresh (a b) (succ a x) (succ b z) (nat-plus a y b))]
   ))
*/
func natplus(x, y, z *ast.SExpr) micro.Goal {
	return Conde(
		[]micro.Goal{micro.EqualO(x, zero), micro.EqualO(y, z)},
		[]micro.Goal{
			micro.CallFresh(func(a *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(b *ast.SExpr) micro.Goal {
					return ConjPlus(
						succ(a, x),
						succ(b, z),
						natplus(a, y, b),
					)
				})
			}),
		},
	)
}

func TestNatPlus(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return natplus(makenat(3), makenat(6), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := parsenat(sexprs[0]), 9
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
}

/*
(define (fib x y)
    (conde
    [(== x 0) (== y 0)]
    [(== x (make-nat 1)) (== y (make-nat 1))]
    [(fresh (n1 n2 f1 f2)
            (succ n1 x)
            (succ n2 n1)
            (fib n1 f1)
            (fib n2 f2)
            (nat-plus f1 f2 y))]
  ))
*/
func fib(x, y *ast.SExpr) micro.Goal {
	return Conde(
		[]micro.Goal{micro.EqualO(x, zero), micro.EqualO(y, zero)},
		[]micro.Goal{micro.EqualO(x, one), micro.EqualO(y, one)},
		[]micro.Goal{
			micro.CallFresh(func(n1 *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(n2 *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(f1 *ast.SExpr) micro.Goal {
						return micro.CallFresh(func(f2 *ast.SExpr) micro.Goal {
							return ConjPlus(
								succ(n1, x),
								succ(n2, n1),
								fib(n1, f1),
								fib(n2, f2),
								natplus(f1, f2, y),
							)
						})
					})
				})
			}),
		},
	)
}

func TestFibonacci(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return fib(makenat(15), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := parsenat(sexprs[0]), 610
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
}
