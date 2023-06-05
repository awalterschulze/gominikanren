package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func SimplO(r, s *ast.SExpr) comicro.Goal {
	return comini.IfThenElseO(
		simplO(r, s),
		simplO(r, s),
		comicro.EqualO(r, s),
	)
}

func simplO(r, s *ast.SExpr) comicro.Goal {
	return comicro.CallFresh(func(r1 *ast.SExpr) comicro.Goal {
		return comicro.CallFresh(func(r2 *ast.SExpr) comicro.Goal {
			return comicro.CallFresh(func(sr1 *ast.SExpr) comicro.Goal {
				return comicro.CallFresh(func(sr2 *ast.SExpr) comicro.Goal {
					return comini.ConjPlus(
						comini.DisjPlus(
							comini.ConjPlus(
								comicro.EqualO(r, Or(r1, r2)),
								SimpleOrO(r1, r2, s),
							),
							comini.ConjPlus(
								comicro.EqualO(r, Concat(r1, r2)),
								SimpleConcatO(r1, r2, s),
							),
						),
					)
				})
			})
		})
	})
}

func ShortCircuitDisj(a, b comicro.Goal) comicro.Goal {
	return comini.IfThenElseO(
		a,
		a,
		b,
	)
}

func SimpleOrO(r1, r2, res *ast.SExpr) comicro.Goal {
	return ShortCircuitDisj(
		comini.DisjPlus(
			comini.ConjPlus(
				comicro.EqualO(r1, EmptySet()),
				comicro.EqualO(res, r2),
			),
			comini.ConjPlus(
				comicro.EqualO(r2, EmptySet()),
				comicro.EqualO(res, r1),
			),
		),
		comicro.EqualO(res, Or(r1, r2)),
	)
}

func SimpleConcatO(r1, r2, res *ast.SExpr) comicro.Goal {
	return ShortCircuitDisj(
		comini.DisjPlus(
			comini.ConjPlus(
				comicro.EqualO(r1, EmptySet()),
				comicro.EqualO(res, EmptySet()),
			),
			comini.ConjPlus(
				comicro.EqualO(r2, EmptySet()),
				comicro.EqualO(res, EmptySet()),
			),
			comini.ConjPlus(
				comicro.EqualO(r1, EmptyStr()),
				comicro.EqualO(res, r2),
			),
			comini.ConjPlus(
				comicro.EqualO(r2, EmptyStr()),
				comicro.EqualO(res, r1),
			),
		),
		comicro.EqualO(res, Concat(r1, r2)),
	)
}
