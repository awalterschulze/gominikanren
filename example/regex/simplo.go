package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
)

func SimpleOrO(r1, r2, res *Regex) comicro.Goal {
	r := comicro.EqualO(res, Or(r1, r2))
	return comicro.Disjs(
		comicro.Conj(
			IsEmptySet(r1),
			comicro.EqualO(res, r2),
		),
		comicro.Conj(
			IsEmptySet(r2),
			comicro.EqualO(res, r1),
		),
		comicro.Conjs(IsEmptyStr(r1), IsNotEmptySet(r2), r),
		comicro.Conjs(IsEmptyStr(r2), IsNotEmptySet(r1), r),
		comicro.Conjs(IsChar(r1), IsNotEmptySet(r2), r),
		comicro.Conjs(IsChar(r2), IsNotEmptySet(r1), r),
		comicro.Conjs(IsOr(r1), IsNotEmptySet(r2), r),
		comicro.Conjs(IsOr(r2), IsNotEmptySet(r1), r),
		comicro.Conjs(IsConcat(r1), IsNotEmptySet(r2), r),
		comicro.Conjs(IsConcat(r2), IsNotEmptySet(r1), r),
		comicro.Conjs(IsStar(r1), IsNotEmptySet(r2), r),
		comicro.Conjs(IsStar(r2), IsNotEmptySet(r1), r),
	)
}

func SimpleConcatO(r1, r2, res *Regex) comicro.Goal {
	r := comicro.EqualO(res, Concat(r1, r2))
	return comicro.Disjs(
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
		comicro.Conjs(IsChar(r1), IsNotEmpty(r2), r),
		comicro.Conjs(IsChar(r2), IsNotEmpty(r1), r),
		comicro.Conjs(IsOr(r1), IsNotEmpty(r2), r),
		comicro.Conjs(IsOr(r2), IsNotEmpty(r1), r),
		comicro.Conjs(IsConcat(r1), IsNotEmpty(r2), r),
		comicro.Conjs(IsConcat(r2), IsNotEmpty(r1), r),
		comicro.Conjs(IsStar(r1), IsNotEmpty(r2), r),
		comicro.Conjs(IsStar(r2), IsNotEmpty(r1), r),
	)
}

func IsNotEmpty(r *Regex) comicro.Goal {
	return comicro.Disjs(
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsNotEmptySet(r *Regex) comicro.Goal {
	return comicro.Disjs(
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
	return comicro.Exists(func(c *rune) comicro.Goal {
		return comicro.EqualO(r, CharPtr(c))
	})
}

func IsOr(r *Regex) comicro.Goal {
	return comicro.Exists(func(r1 *Regex) comicro.Goal {
		return comicro.Exists(func(r2 *Regex) comicro.Goal {
			return comicro.EqualO(r, Or(r1, r2))
		})
	})
}

func IsConcat(r *Regex) comicro.Goal {
	return comicro.Exists(func(r1 *Regex) comicro.Goal {
		return comicro.Exists(func(r2 *Regex) comicro.Goal {
			return comicro.EqualO(r, Concat(r1, r2))
		})
	})
}

func IsStar(r *Regex) comicro.Goal {
	return comicro.Exists(func(r1 *Regex) comicro.Goal {
		return comicro.EqualO(r, Star(r1))
	})
}
