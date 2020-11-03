package peano

import (
	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// peano numbers and extralogical convenience functions
var Zero = ast.NewInt(0)
var One = Makenat(1)

/*
   (define (succ x y)
       (== y `(,x)))
*/
func Succ(prev, next *ast.SExpr) micro.Goal {
	return micro.EqualO(next, ast.Cons(prev, nil))
}

func Makenat(n int) *ast.SExpr {
	if n == 0 {
		return Zero
	}
	return ast.Cons(Makenat(n-1), nil)
}

func Parsenat(x *ast.SExpr) int {
	if x == Zero {
		return 0
	}
	return Parsenat(x.Car()) + 1
}

func ToList(list []int) *ast.SExpr {
	var out *ast.SExpr = nil
	for i := 0; i < len(list); i++ {
		out = ast.Cons(Makenat(list[len(list)-i-1]), out)
	}
	return out
}

/*
   ;; x + y = z
   (define (nat-plus x y z)
   (conde
       [(== x 0) (== y z)]
       [(fresh (a b) (succ a x) (succ b z) (nat-plus a y b))]
   ))
*/
func Natplus(x, y, z *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, Zero), micro.EqualO(y, z)},
		[]micro.Goal{
			micro.CallFresh(func(a *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(b *ast.SExpr) micro.Goal {
					return mini.ConjPlus(
						Succ(a, x),
						Succ(b, z),
						Natplus(a, y, b),
					)
				})
			}),
		},
	)
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
func Leq(x, y *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, Zero)},
		[]micro.Goal{
			micro.CallFresh(func(n *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(m *ast.SExpr) micro.Goal {
					return mini.ConjPlus(
						Succ(n, x),
						Succ(m, y),
						Leq(n, m),
					)
				})
			}),
		},
	)
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
func Half(x, y *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, Zero), micro.EqualO(y, Zero)},
		[]micro.Goal{micro.EqualO(x, One), micro.EqualO(y, Zero)},
		[]micro.Goal{
			micro.CallFresh(func(n *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(m *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(z *ast.SExpr) micro.Goal {
						return mini.ConjPlus(
							Succ(m, y),
							Succ(z, x),
							Succ(n, z),
							Half(n, m),
						)
					})
				})
			}),
		},
	)
}
