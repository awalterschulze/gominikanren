package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func SimplO(r, s *ast.SExpr) comicro.Goal {
	return ShortCircuitDisj(
		simplO(r, s),
		comicro.EqualO(r, s),
	)
}

func ShortCircuitDisj(a, b comicro.Goal) comicro.Goal {
	return comini.IfThenElseO(
		a,
		a,
		b,
	)
}

func simplO(r, s *ast.SExpr) comicro.Goal {
	return comicro.Fresh(2, func(vars ...*ast.SExpr) comicro.Goal {
		r1, r2 := vars[0], vars[1]
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
}

func SimpleOrO(r1, r2, res *ast.SExpr) comicro.Goal {
	r := comicro.EqualO(res, Or(r1, r2))
	return comini.DisjPlus(
		comicro.Conj(
			IsEmptySet(r1),
			comicro.EqualO(res, r2),
		),
		comicro.Conj(
			IsEmptySet(r2),
			comicro.EqualO(res, r1),
		),
		comini.ConjPlus(IsEmptyStr(r1), IsNotEmptySet(r2), r),
		comini.ConjPlus(IsEmptyStr(r2), IsNotEmptySet(r1), r),
		comini.ConjPlus(IsChar(r1), IsNotEmptySet(r2), r),
		comini.ConjPlus(IsChar(r2), IsNotEmptySet(r1), r),
		comini.ConjPlus(IsOr(r1), IsNotEmptySet(r2), r),
		comini.ConjPlus(IsOr(r2), IsNotEmptySet(r1), r),
		comini.ConjPlus(IsConcat(r1), IsNotEmptySet(r2), r),
		comini.ConjPlus(IsConcat(r2), IsNotEmptySet(r1), r),
		comini.ConjPlus(IsStar(r1), IsNotEmptySet(r2), r),
		comini.ConjPlus(IsStar(r2), IsNotEmptySet(r1), r),
	)
}

func SimpleConcatO(r1, r2, res *ast.SExpr) comicro.Goal {
	r := comicro.EqualO(res, Concat(r1, r2))
	return comini.DisjPlus(
		comicro.Conj(
			IsEmptySet(r1),
			comicro.EqualO(res, EmptySet()),
		),
		comicro.Conj(
			IsEmptySet(r2),
			comicro.EqualO(res, EmptySet()),
		),
		comicro.Conj(
			IsEmptyStr(r1),
			comicro.EqualO(res, r2),
		),
		comicro.Conj(
			IsEmptyStr(r2),
			comicro.EqualO(res, r1),
		),
		comini.ConjPlus(IsChar(r1), IsNotEmpty(r2), r),
		comini.ConjPlus(IsChar(r2), IsNotEmpty(r1), r),
		comini.ConjPlus(IsOr(r1), IsNotEmpty(r2), r),
		comini.ConjPlus(IsOr(r2), IsNotEmpty(r1), r),
		comini.ConjPlus(IsConcat(r1), IsNotEmpty(r2), r),
		comini.ConjPlus(IsConcat(r2), IsNotEmpty(r1), r),
		comini.ConjPlus(IsStar(r1), IsNotEmpty(r2), r),
		comini.ConjPlus(IsStar(r2), IsNotEmpty(r1), r),
	)
}

func IsNotEmpty(r *ast.SExpr) comicro.Goal {
	return comini.DisjPlus(
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsNotEmptySet(r *ast.SExpr) comicro.Goal {
	return comini.DisjPlus(
		IsEmptyStr(r),
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsEmptySet(r *ast.SExpr) comicro.Goal {
	return comicro.EqualO(r, EmptySet())
}

func IsEmptyStr(r *ast.SExpr) comicro.Goal {
	return comicro.EqualO(r, EmptyStr())
}

func IsChar(r *ast.SExpr) comicro.Goal {
	return comicro.CallFresh(func(c *ast.SExpr) comicro.Goal {
		return comicro.EqualO(r, CharFromSExpr(c))
	})
}

func IsOr(r *ast.SExpr) comicro.Goal {
	return comicro.Fresh(2, func(vars ...*ast.SExpr) comicro.Goal {
		r1, r2 := vars[0], vars[1]
		return comicro.EqualO(r, Or(r1, r2))
	})
}

func IsConcat(r *ast.SExpr) comicro.Goal {
	return comicro.Fresh(2, func(vars ...*ast.SExpr) comicro.Goal {
		r1, r2 := vars[0], vars[1]
		return comicro.EqualO(r, Concat(r1, r2))
	})
}

func IsStar(r *ast.SExpr) comicro.Goal {
	return comicro.CallFresh(func(r1 *ast.SExpr) comicro.Goal {
		return comicro.EqualO(r, Star(r1))
	})
}
