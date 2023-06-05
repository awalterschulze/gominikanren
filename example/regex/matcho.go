package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func MatchO(r, s, res *ast.SExpr) comicro.Goal {
	if s == nil {
		return NullO(r, res)
	}
	if !s.IsPair() {
		return comicro.CallFresh(func(dr *ast.SExpr) comicro.Goal {
			return comicro.Conj(
				SDerivO(r, s, dr),
				NullO(dr, res),
			)
		})
	}
	return comicro.CallFresh(func(dr *ast.SExpr) comicro.Goal {
		return comicro.Conj(
			SDerivO(r, s.Car(), dr),
			MatchO(dr, s.Cdr(), res),
		)
	})
}
