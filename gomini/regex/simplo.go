package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func SimpleOrO(r1, r2, res *Regex) Goal {
	r := EqualO(res, Or(r1, r2))
	return Disjs(
		Conj(
			IsEmptySet(r1),
			EqualO(res, r2),
		),
		Conj(
			IsEmptySet(r2),
			EqualO(res, r1),
		),
		Conjs(IsEmptyStr(r1), IsNotEmptySet(r2), r),
		Conjs(IsEmptyStr(r2), IsNotEmptySet(r1), r),
		Conjs(IsChar(r1), IsNotEmptySet(r2), r),
		Conjs(IsChar(r2), IsNotEmptySet(r1), r),
		Conjs(IsOr(r1), IsNotEmptySet(r2), r),
		Conjs(IsOr(r2), IsNotEmptySet(r1), r),
		Conjs(IsConcat(r1), IsNotEmptySet(r2), r),
		Conjs(IsConcat(r2), IsNotEmptySet(r1), r),
		Conjs(IsStar(r1), IsNotEmptySet(r2), r),
		Conjs(IsStar(r2), IsNotEmptySet(r1), r),
	)
}

func SimpleConcatO(r1, r2, res *Regex) Goal {
	r := EqualO(res, Concat(r1, r2))
	return Disjs(
		Conj(
			IsEmptySet(r1),
			EqualO(res, EmptySet()),
		),
		Conj(
			IsEmptySet(r2),
			EqualO(res, EmptySet()),
		),
		Conj(
			IsEmptyStr(r1),
			EqualO(res, r2),
		),
		Conj(
			IsEmptyStr(r2),
			EqualO(res, r1),
		),
		Conjs(IsChar(r1), IsNotEmpty(r2), r),
		Conjs(IsChar(r2), IsNotEmpty(r1), r),
		Conjs(IsOr(r1), IsNotEmpty(r2), r),
		Conjs(IsOr(r2), IsNotEmpty(r1), r),
		Conjs(IsConcat(r1), IsNotEmpty(r2), r),
		Conjs(IsConcat(r2), IsNotEmpty(r1), r),
		Conjs(IsStar(r1), IsNotEmpty(r2), r),
		Conjs(IsStar(r2), IsNotEmpty(r1), r),
	)
}

func IsNotEmpty(r *Regex) Goal {
	return Disjs(
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsNotEmptySet(r *Regex) Goal {
	return Disjs(
		IsEmptyStr(r),
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsEmptySet(r *Regex) Goal {
	return EqualO(r, EmptySet())
}

func IsEmptyStr(r *Regex) Goal {
	return EqualO(r, EmptyStr())
}

func IsChar(r *Regex) Goal {
	return Exists(func(c *rune) Goal {
		return EqualO(r, CharPtr(c))
	})
}

func IsOr(r *Regex) Goal {
	return Exists(func(r1 *Regex) Goal {
		return Exists(func(r2 *Regex) Goal {
			return EqualO(r, Or(r1, r2))
		})
	})
}

func IsConcat(r *Regex) Goal {
	return Exists(func(r1 *Regex) Goal {
		return Exists(func(r2 *Regex) Goal {
			return EqualO(r, Concat(r1, r2))
		})
	})
}

func IsStar(r *Regex) Goal {
	return Exists(func(r1 *Regex) Goal {
		return EqualO(r, Star(r1))
	})
}
