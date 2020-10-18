package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
AppendO is a goal that appends two lists into the third list.

scheme code:

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
    return micro.Zzz(
				micro.DisjointO(
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
				),
		)
}

// NullO is a goal that checks if a list is null.
func NullO(x *ast.SExpr) micro.Goal {
	return micro.EqualO(x, nil)
}

// ConsO is a goal that conses the first two expressions into the third.
func ConsO(a, d, p *ast.SExpr) micro.Goal {
    return func() micro.GoalFn {
	return func(s *micro.State) micro.StreamOfStates {
		l := ast.Cons(a, d)
		return micro.EqualO(l, p)()(s)
	}
    }
}

// CarO is a goal where the second parameter is the head of the list in the first parameter.
func CarO(p, a *ast.SExpr) micro.Goal {
	return micro.CallFresh(func(d *ast.SExpr) micro.Goal {
		return micro.EqualO(ast.Cons(a, d), p)
	})
}
