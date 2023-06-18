package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func SDerivO(r, char, out *ast.SExpr) comicro.Goal {
	return comini.Disjs(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptySet()),
		),
		DeriveCharO(r, char, out),
		comicro.CallFresh(&ast.SExpr{}, func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(&ast.SExpr{}, func(da *ast.SExpr) comicro.Goal {
				return comicro.CallFresh(&ast.SExpr{}, func(b *ast.SExpr) comicro.Goal {
					return comicro.CallFresh(&ast.SExpr{}, func(db *ast.SExpr) comicro.Goal {
						return comini.Conjs(
							comicro.EqualO(r, Or(a, b)),
							SDerivO(a, char, da),
							SDerivO(b, char, db),
							SimpleOrO(da, db, out),
						)
					})
				})
			})
		}),
		comicro.CallFresh(&ast.SExpr{}, func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(&ast.SExpr{}, func(da *ast.SExpr) comicro.Goal {
				return comicro.CallFresh(&ast.SExpr{}, func(na *ast.SExpr) comicro.Goal {
					return comicro.CallFresh(&ast.SExpr{}, func(b *ast.SExpr) comicro.Goal {
						return comicro.CallFresh(&ast.SExpr{}, func(db *ast.SExpr) comicro.Goal {
							return comicro.CallFresh(&ast.SExpr{}, func(ca *ast.SExpr) comicro.Goal {
								return comicro.CallFresh(&ast.SExpr{}, func(cb *ast.SExpr) comicro.Goal {
									return comini.Conjs(
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
		comicro.CallFresh(&ast.SExpr{}, func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(&ast.SExpr{}, func(da *ast.SExpr) comicro.Goal {
				return comini.Conjs(
					comicro.EqualO(r, Star(a)),
					SDerivO(a, char, da),
					SimpleConcatO(da, r, out),
				)
			})
		}),
	)
}
