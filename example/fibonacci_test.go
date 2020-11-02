package example

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

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
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, zero), micro.EqualO(y, zero)},
		[]micro.Goal{micro.EqualO(x, one), micro.EqualO(y, one)},
		[]micro.Goal{
			micro.CallFresh(func(n1 *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(n2 *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(f1 *ast.SExpr) micro.Goal {
						return micro.CallFresh(func(f2 *ast.SExpr) micro.Goal {
							return mini.ConjPlus(
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
