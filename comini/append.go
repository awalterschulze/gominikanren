package comini

import (
	"context"

	"github.com/awalterschulze/gominikanren/comicro"
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
func AppendO(l, t, out *ast.SExpr) comicro.Goal {
	return comicro.Disj(
		comicro.Conj(
			NullO(l),
			comicro.EqualO(t, out),
		),
		comicro.Fresh(3, func(vars ...*ast.SExpr) comicro.Goal {
			a, d, res := vars[0], vars[1], vars[2]
			return comicro.Conj(
				ConsO(a, d, l),
				comicro.Conj(
					ConsO(a, res, out),
					AppendO(d, t, res),
				),
			)
		}),
	)
}

// NullO is a goal that checks if a list is null.
func NullO(x *ast.SExpr) comicro.Goal {
	return comicro.EqualO(x, nil)
}

// ConsO is a goal that conses the first two expressions into the third.
func ConsO(a, d, p *ast.SExpr) comicro.Goal {
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		l := ast.Cons(a, d)
		comicro.EqualO(l, p)(ctx, s, ss)
	}
}

// CarO is a goal where the second parameter is the head of the list in the first parameter.
func CarO(p, a *ast.SExpr) comicro.Goal {
	return comicro.CallFresh(func(d *ast.SExpr) comicro.Goal {
		return comicro.EqualO(ast.Cons(a, d), p)
	})
}
