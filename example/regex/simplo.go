package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
)

func SimplO(r, s *Regex) comicro.Goal {
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

func simplO(r, s *Regex) comicro.Goal {
	return comicro.CallFresh(func(r1 *Regex) comicro.Goal {
		return comicro.CallFresh(func(r2 *Regex) comicro.Goal {
			return comini.Conjs(
				comini.Disjs(
					comini.Conjs(
						comicro.EqualO(r, Or(r1, r2)),
						SimpleOrO(r1, r2, s),
					),
					comini.Conjs(
						comicro.EqualO(r, Concat(r1, r2)),
						SimpleConcatO(r1, r2, s),
					),
				),
			)
		})
	})
}

func SimpleOrO(r1, r2, res *Regex) comicro.Goal {
	r := comicro.EqualO(res, Or(r1, r2))
	return comini.Disjs(
		comicro.Conj(
			IsEmptySet(r1),
			comicro.EqualO(res, r2),
		),
		comicro.Conj(
			IsEmptySet(r2),
			comicro.EqualO(res, r1),
		),
		comini.Conjs(IsEmptyStr(r1), IsNotEmptySet(r2), r),
		comini.Conjs(IsEmptyStr(r2), IsNotEmptySet(r1), r),
		comini.Conjs(IsChar(r1), IsNotEmptySet(r2), r),
		comini.Conjs(IsChar(r2), IsNotEmptySet(r1), r),
		comini.Conjs(IsOr(r1), IsNotEmptySet(r2), r),
		comini.Conjs(IsOr(r2), IsNotEmptySet(r1), r),
		comini.Conjs(IsConcat(r1), IsNotEmptySet(r2), r),
		comini.Conjs(IsConcat(r2), IsNotEmptySet(r1), r),
		comini.Conjs(IsStar(r1), IsNotEmptySet(r2), r),
		comini.Conjs(IsStar(r2), IsNotEmptySet(r1), r),
	)
}

func SimpleConcatO(r1, r2, res *Regex) comicro.Goal {
	r := comicro.EqualO(res, Concat(r1, r2))
	return comini.Disjs(
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
		comini.Conjs(IsChar(r1), IsNotEmpty(r2), r),
		comini.Conjs(IsChar(r2), IsNotEmpty(r1), r),
		comini.Conjs(IsOr(r1), IsNotEmpty(r2), r),
		comini.Conjs(IsOr(r2), IsNotEmpty(r1), r),
		comini.Conjs(IsConcat(r1), IsNotEmpty(r2), r),
		comini.Conjs(IsConcat(r2), IsNotEmpty(r1), r),
		comini.Conjs(IsStar(r1), IsNotEmpty(r2), r),
		comini.Conjs(IsStar(r2), IsNotEmpty(r1), r),
	)
}

func IsNotEmpty(r *Regex) comicro.Goal {
	return comini.Disjs(
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsNotEmptySet(r *Regex) comicro.Goal {
	return comini.Disjs(
		IsEmptyStr(r),
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsEmptySet(r *Regex) comicro.Goal {
	return comicro.EqualO(r, EmptySet())
}

func IsEmptyStr(r *Regex) comicro.Goal {
	return comicro.EqualO(r, EmptyStr())
}

func IsChar(r *Regex) comicro.Goal {
	return comicro.CallFresh(func(c *rune) comicro.Goal {
		return comicro.EqualO(r, CharPtr(c))
	})
}

func IsOr(r *Regex) comicro.Goal {
	return comicro.CallFresh(func(r1 *Regex) comicro.Goal {
		return comicro.CallFresh(func(r2 *Regex) comicro.Goal {
			return comicro.EqualO(r, Or(r1, r2))
		})
	})
}

func IsConcat(r *Regex) comicro.Goal {
	return comicro.CallFresh(func(r1 *Regex) comicro.Goal {
		return comicro.CallFresh(func(r2 *Regex) comicro.Goal {
			return comicro.EqualO(r, Concat(r1, r2))
		})
	})
}

func IsStar(r *Regex) comicro.Goal {
	return comicro.CallFresh(func(r1 *Regex) comicro.Goal {
		return comicro.EqualO(r, Star(r1))
	})
}
