package example

import (
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
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
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, zero), micro.EqualO(y, z)},
		[]micro.Goal{
			micro.CallFresh(func(a *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(b *ast.SExpr) micro.Goal {
					return mini.ConjPlus(
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
(define (leq x y)
  (conde
    [(== x 0)]
    [(fresh (n m)
        (succ n x)
        (succ m y)
        (leq n m)
    )]
  ))
*/
func leq(x, y *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, zero)},
		[]micro.Goal{
			micro.CallFresh(func(n *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(m *ast.SExpr) micro.Goal {
					return mini.ConjPlus(
						succ(n, x),
						succ(m, y),
						leq(n, m),
					)
				})
			}),
		},
	)
}

func TestLeq(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return leq(q, makenat(6))
	})
	if len(sexprs) != 7 {
		t.Fatalf("expected len %d, but got len %d instead", 7, len(sexprs))
	}
	// answer we're looking for is the list of 0..6 incl
	for i, sexpr := range sexprs {
		got, want := parsenat(sexpr), i
		if got != want {
			t.Fatalf("expected %d, but got %d instead", want, got)
		}
	}
	// we expect q to be any successor of 4
	sexprs = micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return leq(makenat(4), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	// representing succ as wrapping with a list, we expect to find a
	// var (reified to symbol _0) after taking car of answer 4 times
	got := sexprs[0]
	for i := 0; i < 4; i++ {
		got = got.Car()
	}
	want := ast.NewSymbol("_0")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, but got %v instead", want, got)
	}
}

/*
;; y = x // 2, integer division rounding down
(define (half x y)
    (conde
    [(== x 0) (== y 0)]
    [(== x (make-nat 1)) (== y 0)]
    [(fresh (n m z)
        (succ m y)
        (succ z x)
        (succ n z)
        (half n m)
    )]
  ))
*/
func half(x, y *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, zero), micro.EqualO(y, zero)},
		[]micro.Goal{micro.EqualO(x, one), micro.EqualO(y, zero)},
		[]micro.Goal{
			micro.CallFresh(func(n *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(m *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(z *ast.SExpr) micro.Goal {
						return mini.ConjPlus(
							succ(m, y),
							succ(z, x),
							succ(n, z),
							half(n, m),
						)
					})
				})
			}),
		},
	)
}

func TestHalf(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return half(makenat(7), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := parsenat(sexprs[0]), 3
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
	sexprs = micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return half(q, makenat(5))
	})
	if len(sexprs) != 2 {
		t.Fatalf("expected len %d, but got len %d instead", 2, len(sexprs))
	}
	got, want = parsenat(sexprs[0]), 10
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
	got, want = parsenat(sexprs[1]), 11
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
}

func toPeanoList(list []int) *ast.SExpr {
	var out *ast.SExpr = nil
	for i := 0; i < len(list); i++ {
		out = ast.Cons(makenat(list[len(list)-i-1]), out)
	}
	return out
}
