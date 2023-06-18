package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func NullO(r, out *ast.SExpr) comicro.Goal {
	return comini.Disjs(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptyStr()),
		),
		comicro.Conj(
			comicro.CallFresh(&ast.SExpr{}, func(char *ast.SExpr) comicro.Goal {
				return comicro.EqualO(r, CharFromSExpr(char))
			}),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.CallFresh(&ast.SExpr{}, func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(&ast.SExpr{}, func(b *ast.SExpr) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Or(a, b)),
					comini.Disjs(
						comini.Conjs(
							NullO(a, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.Conjs(
							NullO(b, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.Conjs(
							NullO(a, EmptySet()),
							NullO(b, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
					),
				)
			})
		}),
		comicro.CallFresh(&ast.SExpr{}, func(a *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(&ast.SExpr{}, func(b *ast.SExpr) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Concat(a, b)),
					comini.Disjs(
						comini.Conjs(
							NullO(a, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.Conjs(
							NullO(b, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.Conjs(
							NullO(a, EmptyStr()),
							NullO(b, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
					),
				)
			})
		}),
		comicro.CallFresh(&ast.SExpr{}, func(a *ast.SExpr) comicro.Goal {
			return comicro.Conj(
				comicro.EqualO(r, Star(a)),
				comicro.EqualO(out, EmptyStr()),
			)
		}),
	)
}
