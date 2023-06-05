package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func SDerivO(r, char, out *ast.SExpr) comicro.Goal {
	return comini.DisjPlus(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.CallFresh(func(c *ast.SExpr) comicro.Goal {
			return comicro.Conj(
				comicro.EqualO(r, CharFromSExpr(c)),
				comini.IfThenElseO(
					comicro.EqualO(char, c),
					comicro.EqualO(out, EmptyStr()),
					comicro.EqualO(out, EmptySet()),
				),
			)
		}),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(da *ast.SExpr) comicro.Goal {
				return comicro.CallFresh(func(b *ast.SExpr) comicro.Goal {
					return comicro.CallFresh(func(db *ast.SExpr) comicro.Goal {
						return comini.ConjPlus(
							comicro.EqualO(r, Or(a, b)),
							SDerivO(a, char, da),
							SDerivO(b, char, db),
							SimpleOrO(da, db, out),
						)
					})
				})
			})
		}),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(da *ast.SExpr) comicro.Goal {
				return comicro.CallFresh(func(na *ast.SExpr) comicro.Goal {
					return comicro.CallFresh(func(b *ast.SExpr) comicro.Goal {
						return comicro.CallFresh(func(db *ast.SExpr) comicro.Goal {
							return comicro.CallFresh(func(ca *ast.SExpr) comicro.Goal {
								return comicro.CallFresh(func(cb *ast.SExpr) comicro.Goal {
									return comini.ConjPlus(
										comicro.EqualO(r, Concat(a, b)),
										SDerivO(a, char, da),
										SDerivO(b, char, db),
										NullO(a, na),
										SimpleConcatO(da, b, ca),
										SimpleConcatO(na, db, cb),
										SimpleOrO(ca, cb, out),
									)
								})
							})
						})
					})
				})
			})
		}),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(da *ast.SExpr) comicro.Goal {
				return comini.ConjPlus(
					comicro.EqualO(r, Star(a)),
					SDerivO(a, char, da),
					SimpleConcatO(da, r, out),
				)
			})
		}),
	)
}
