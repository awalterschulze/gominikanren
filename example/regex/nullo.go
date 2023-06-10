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
			comicro.CallFresh(func(char comicro.Var) comicro.Goal {
				return comicro.EqualO(r, CharFromSExpr(char.SExpr()))
			}),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(b comicro.Var) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Or(a.SExpr(), b.SExpr())),
					comini.Disjs(
						comini.Conjs(
							NullO(a.SExpr(), EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.Conjs(
							NullO(b.SExpr(), EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.Conjs(
							NullO(a.SExpr(), EmptySet()),
							NullO(b.SExpr(), EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
					),
				)
			})
		}),
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.CallFresh(func(b comicro.Var) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Concat(a.SExpr(), b.SExpr())),
					comini.Disjs(
						comini.Conjs(
							NullO(a.SExpr(), EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.Conjs(
							NullO(b.SExpr(), EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.Conjs(
							NullO(a.SExpr(), EmptyStr()),
							NullO(b.SExpr(), EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
					),
				)
			})
		}),
		comicro.CallFresh(func(a comicro.Var) comicro.Goal {
			return comicro.Conj(
				comicro.EqualO(r, Star(a.SExpr())),
				comicro.EqualO(out, EmptyStr()),
			)
		}),
	)
}
