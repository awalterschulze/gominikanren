package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func MatchO(r, s, res *ast.SExpr) comicro.Goal {
	if s == nil {
		return NullO(r, res)
	}
	return comicro.CallFresh(func(dr *ast.SExpr) comicro.Goal {
		return comicro.Conj(
			SDerivOs(r, s, dr),
			NullO(dr, res),
		)
	})
}

func SDerivOs(r, s, res *ast.SExpr) comicro.Goal {
	if s == nil {
		return comicro.EqualO(res, r)
	}
	if !s.IsPair() {
		return SDerivO(r, s, res)
	}
	return comicro.CallFresh(func(dr *ast.SExpr) comicro.Goal {
		return comicro.Conj(
			SDerivO(r, s.Car(), dr),
			SDerivOs(dr, s.Cdr(), res),
		)
	})
}
