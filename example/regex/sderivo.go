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
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(da comicro.Var) comicro.Goal {
				return comicro.CallFresh(func(b comicro.Var) comicro.Goal {
					return comicro.CallFresh(func(db comicro.Var) comicro.Goal {
						return comini.Conjs(
							comicro.EqualO(r, Or(a.SExpr(), b.SExpr())),
							SDerivO(a.SExpr(), char, da.SExpr()),
							SDerivO(b.SExpr(), char, db.SExpr()),
							SimpleOrO(da.SExpr(), db.SExpr(), out),
						)
					})
				})
			})
		}),
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(da comicro.Var) comicro.Goal {
				return comicro.CallFresh(func(na comicro.Var) comicro.Goal {
					return comicro.CallFresh(func(b comicro.Var) comicro.Goal {
						return comicro.CallFresh(func(db comicro.Var) comicro.Goal {
							return comicro.CallFresh(func(ca comicro.Var) comicro.Goal {
								return comicro.CallFresh(func(cb comicro.Var) comicro.Goal {
									return comini.Conjs(
										comicro.EqualO(r, Concat(a.SExpr(), b.SExpr())),
										SDerivO(a.SExpr(), char, da.SExpr()),
										SDerivO(b.SExpr(), char, db.SExpr()),
										NullO(a.SExpr(), na.SExpr()),
										SimpleConcatO(da.SExpr(), b.SExpr(), ca.SExpr()),
										SimpleConcatO(na.SExpr(), db.SExpr(), cb.SExpr()),
										SimpleOrO(ca.SExpr(), cb.SExpr(), out),
									)
								})
							})
						})
					})
				})
			})
		}),
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(da comicro.Var) comicro.Goal {
				return comini.Conjs(
					comicro.EqualO(r, Star(a.SExpr())),
					SDerivO(a.SExpr(), char, da.SExpr()),
					SimpleConcatO(da.SExpr(), r, out),
				)
			})
		}),
	)
}
