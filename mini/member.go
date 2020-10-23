package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
Membero checks membership in a list

(define (membero x y)
  (fresh (a d)
  (== y `(,a . ,d))
  (conde
    [(== x a)]
    [(membero x d)]
    )))
*/
func Membero(x *ast.SExpr, y *ast.SExpr) micro.Goal {
	return micro.Zzz(micro.CallFresh(func(a *ast.SExpr) micro.Goal {
		return micro.CallFresh(func(d *ast.SExpr) micro.Goal {
			return micro.ConjunctionO(
				micro.EqualO(y, ast.Cons(a, d)),
				Conde(
					[]micro.Goal{micro.EqualO(x, a)},
					[]micro.Goal{Membero(x, d)},
				),
			)
		})
	}))
}
