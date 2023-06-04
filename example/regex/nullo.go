package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func NullO(r, out *ast.SExpr) comicro.Goal {
	return comini.DisjPlus(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptyStr()),
		),
		comicro.Conj(
			comicro.CallFresh(func(char *ast.SExpr) comicro.Goal {
				return comicro.EqualO(r, CharFromSExpr(char))
			}),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(b *ast.SExpr) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Or(a, b)),
					comini.DisjPlus(
						comini.ConjPlus(
							NullO(a, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.ConjPlus(
							NullO(b, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.ConjPlus(
							NullO(a, EmptySet()),
							NullO(b, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
					),
				)
			})
		}),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(b *ast.SExpr) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Concat(a, b)),
					comini.DisjPlus(
						comini.ConjPlus(
							NullO(a, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.ConjPlus(
							NullO(b, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.ConjPlus(
							NullO(a, EmptyStr()),
							NullO(b, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
					),
				)
			})
		}),
		comicro.CallFresh(func(a *ast.SExpr) comicro.Goal {
			return comicro.Conj(
				comicro.EqualO(r, Star(a)),
				comicro.EqualO(out, EmptyStr()),
			)
		}),
	)
}
