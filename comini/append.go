package comini

import (
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
		comicro.Exists(func(a *ast.SExpr) comicro.Goal {
			return comicro.Exists(func(d *ast.SExpr) comicro.Goal {
				return comicro.Exists(func(res *ast.SExpr) comicro.Goal {
					return comicro.Conj(
						ConsO(a, d, l),
						comicro.Conj(
							ConsO(a, res, out),
							AppendO(d, t, res),
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
	return comicro.EqualO(ast.Cons(a, d), p)
}

// CarO is a goal where the second parameter is the head of the list in the first parameter.
func CarO(p, a *ast.SExpr) comicro.Goal {
	return comicro.Exists(func(d *ast.SExpr) comicro.Goal {
		return comicro.EqualO(ast.Cons(a, d), p)
	})
}
