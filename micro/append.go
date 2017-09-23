package micro

import "github.com/awalterschulze/gominikanren/sexpr/ast"

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
func AppendO(l, t, out *ast.SExpr) Goal {
	return func(s Substitution) StreamOfSubstitutions {
		return Suspension(
			func() StreamOfSubstitutions {
				return DisjointO(
					ConjunctionO(
						NullO(l),
						EqualO(t, out),
					),
					CallFresh("a", func(a *ast.SExpr) Goal {
						return CallFresh("d", func(d *ast.SExpr) Goal {
							return CallFresh("res", func(res *ast.SExpr) Goal {
								return ConjunctionO(
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

func NullO(x *ast.SExpr) Goal {
	return EqualO(x, ast.NewList(false))
}

func ConsO(a, d, p *ast.SExpr) Goal {
	return func(s Substitution) StreamOfSubstitutions {
		l := cons(a, d)
		return EqualO(l, p)(s)
	}
}

func cons(car, cdr *ast.SExpr) *ast.SExpr {
	if cdr.List != nil {
		return ast.Prepend(cdr.List.IsQuoted(), car, cdr.List)
	}
	return ast.NewList(false, car, cdr)
}

func CarO(p, a *ast.SExpr) Goal {
	return CallFresh("d", func(d *ast.SExpr) Goal {
		return EqualO(cons(a, d), p)
	})
}
