package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func DerivO(r, char, out *ast.SExpr) comicro.Goal {
	return comini.DisjPlus(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptySet()),
		),
		DeriveCharO(r, char, out),
		comicro.Fresh(4, func(vars ...*ast.SExpr) comicro.Goal {
			a, da, b, db := vars[0], vars[1], vars[2], vars[3]
			return comini.ConjPlus(
				comicro.EqualO(r, Or(a, b)),
				DerivO(a, char, da),
				DerivO(b, char, db),
				comicro.EqualO(out, Or(da, db)),
			)
		}),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(da *ast.SExpr) comicro.Goal {
				return comicro.CallFresh(func(na *ast.SExpr) comicro.Goal {
					return comicro.CallFresh(func(b *ast.SExpr) comicro.Goal {
						return comicro.CallFresh(func(db *ast.SExpr) comicro.Goal {
							return comini.ConjPlus(
								comicro.EqualO(r, Concat(a, b)),
								DerivO(a, char, da),
								DerivO(b, char, db),
								NullO(a, na),
								comicro.EqualO(out, Or(Concat(da, b), Concat(na, db))),
							)
						})
					})
				})
			})
		}),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(da *ast.SExpr) comicro.Goal {
				return comini.ConjPlus(
					comicro.EqualO(r, Star(a)),
					DerivO(a, char, da),
					comicro.EqualO(out, Concat(da, Star(a))),
				)
			})
		}),
	)
}

func DeriveCharO(r, char, out *ast.SExpr) comicro.Goal {
	return comini.DisjPlus(
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('a')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('a'))),
			comicro.EqualO(out, EmptyStr()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('a')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('b'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('a')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('c'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('b')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('a'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('b')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('b'))),
			comicro.EqualO(out, EmptyStr()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('b')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('c'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('c')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('a'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('c')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('b'))),
			comicro.EqualO(out, EmptySet()),
		),
		comini.ConjPlus(
			comicro.EqualO(char, CharSymbol('c')),
			comicro.EqualO(r, CharFromSExpr(CharSymbol('c'))),
			comicro.EqualO(out, EmptyStr()),
		),
	)
}
