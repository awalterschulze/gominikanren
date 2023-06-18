package comini

import (
	"context"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// AppendO is a goal that appends two lists into the third list.
func AppendO(l, t, out *ast.SExpr) comicro.Goal {
	return comicro.Disj(
		comicro.Conj(
			NullO(l),
			comicro.EqualO(t, out),
		),
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(d comicro.Var) comicro.Goal {
				return comicro.CallFresh(func(res comicro.Var) comicro.Goal {
					return comicro.Conj(
						ConsO(a.SExpr(), d.SExpr(), l),
						comicro.Conj(
							ConsO(a.SExpr(), res.SExpr(), out),
							AppendO(d.SExpr(), t, res.SExpr()),
						),
					)
				})
			})
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
	return comicro.CallFresh(func(d comicro.Var) comicro.Goal {
		return comicro.EqualO(ast.Cons(a, d.SExpr()), p)
	})
}
