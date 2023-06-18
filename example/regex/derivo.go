package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func DerivO(r, char, out *ast.SExpr) comicro.Goal {
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
							DerivO(a.SExpr(), char, da.SExpr()),
							DerivO(b.SExpr(), char, db.SExpr()),
							comicro.EqualO(out, Or(da.SExpr(), db.SExpr())),
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
							return comini.Conjs(
								comicro.EqualO(r, Concat(a.SExpr(), b.SExpr())),
								DerivO(a.SExpr(), char, da.SExpr()),
								DerivO(b.SExpr(), char, db.SExpr()),
								NullO(a.SExpr(), na.SExpr()),
								comicro.EqualO(out, Or(Concat(da.SExpr(), b.SExpr()), Concat(na.SExpr(), db.SExpr()))),
							)
						})
					})
				})
			})
		}),
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(da comicro.Var) comicro.Goal {
				return comini.Conjs(
					comicro.EqualO(r, Star(a.SExpr())),
					DerivO(a.SExpr(), char, da.SExpr()),
					comicro.EqualO(out, Concat(da.SExpr(), Star(a.SExpr()))),
				)
			})
		}),
	)
}

func DeriveCharO(r, char, out *ast.SExpr) comicro.Goal {
	return comini.Disjs(
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('a')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('a'))),
			comicro.EqualO(out, EmptyStr()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('a')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('b'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('a')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('c'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('b')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('a'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('b')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('b'))),
			comicro.EqualO(out, EmptyStr()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('b')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('c'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('c')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('a'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('c')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('b'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, CharSymbol('c')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('c'))),
			comicro.EqualO(out, EmptyStr()),
		),
	)
}
