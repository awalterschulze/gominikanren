package mini

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(define (appendo l t out)
	(lambda (s)
		(lambda ()
			(
				(disj
					(conj (nullo l) (== t out))
					(call/fresh 'a
						(lambda (a)
							(call/fresh 'd
								(lambda (d)
									(call/fresh 'res
										(lambda (res)
											(conj
												(conso a d l)
												(conj
													(conso a res out)
													(appendo d t res)
												)
											)
										)
									)
								)
							)
						)
					)
				)
			)
			s
		)
	)
)
*/
func AppendO(l, t, out *ast.SExpr) micro.Goal {
	return func(s *micro.State) micro.StreamOfSubstitutions {
		fmt.Printf("Appendo %v %v %v\n", l, t, out)
		return micro.Suspension(
			func() micro.StreamOfSubstitutions {
				return micro.DisjointO(
					micro.ConjunctionO(
						NullO(l),
						micro.EqualO(t, out),
					),
					micro.CallFresh(func(a *ast.SExpr) micro.Goal {
						return micro.CallFresh(func(d *ast.SExpr) micro.Goal {
							return micro.CallFresh(func(res *ast.SExpr) micro.Goal {
								return micro.ConjunctionO(
									ConsO(a, d, l),
									ConsO(a, res, out),
									AppendO(d, t, res),
								)
							})
						})
					}),
				)(s)
			},
		)
	}
}

func NullO(x *ast.SExpr) micro.Goal {
	return micro.EqualO(x, nil)
}

func ConsO(a, d, p *ast.SExpr) micro.Goal {
	return func(s *micro.State) micro.StreamOfSubstitutions {
		fmt.Printf("Conso %v %v %v\n", a, d, p)
		l := ast.Cons(a, d)
		return micro.EqualO(l, p)(s)
	}
}

func CarO(p, a *ast.SExpr) micro.Goal {
	return micro.CallFresh(func(d *ast.SExpr) micro.Goal {
		return micro.EqualO(ast.Cons(a, d), p)
	})
}
