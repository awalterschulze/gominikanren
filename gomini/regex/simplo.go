package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func SimpleOrO(r1, r2, res *Regex) Goal {
	r := EqualO(res, Or(r1, r2))
	return DisjO(
		ConjO(
			IsEmptySet(r1),
			EqualO(res, r2),
		),
		ConjO(
			IsEmptySet(r2),
			EqualO(res, r1),
		),
		ConjO(IsEmptyStr(r1), IsNotEmptySet(r2), r),
		ConjO(IsEmptyStr(r2), IsNotEmptySet(r1), r),
		ConjO(IsChar(r1), IsNotEmptySet(r2), r),
		ConjO(IsChar(r2), IsNotEmptySet(r1), r),
		ConjO(IsOr(r1), IsNotEmptySet(r2), r),
		ConjO(IsOr(r2), IsNotEmptySet(r1), r),
		ConjO(IsConcat(r1), IsNotEmptySet(r2), r),
		ConjO(IsConcat(r2), IsNotEmptySet(r1), r),
		ConjO(IsStar(r1), IsNotEmptySet(r2), r),
		ConjO(IsStar(r2), IsNotEmptySet(r1), r),
	)
}

func SimpleConcatO(r1, r2, res *Regex) Goal {
	r := EqualO(res, Concat(r1, r2))
	return DisjO(
		ConjO(
			IsEmptySet(r1),
			EqualO(res, EmptySet()),
		),
		ConjO(
			IsEmptySet(r2),
			EqualO(res, EmptySet()),
		),
		ConjO(
			IsEmptyStr(r1),
			EqualO(res, r2),
		),
		ConjO(
			IsEmptyStr(r2),
			EqualO(res, r1),
		),
		ConjO(IsChar(r1), IsNotEmpty(r2), r),
		ConjO(IsChar(r2), IsNotEmpty(r1), r),
		ConjO(IsOr(r1), IsNotEmpty(r2), r),
		ConjO(IsOr(r2), IsNotEmpty(r1), r),
		ConjO(IsConcat(r1), IsNotEmpty(r2), r),
		ConjO(IsConcat(r2), IsNotEmpty(r1), r),
		ConjO(IsStar(r1), IsNotEmpty(r2), r),
		ConjO(IsStar(r2), IsNotEmpty(r1), r),
	)
}

func IsNotEmpty(r *Regex) Goal {
	return DisjO(
		IsChar(r),
		IsOr(r),
		IsConcat(r),
		IsStar(r),
	)
}

func IsNotEmptySet(r *Regex) Goal {
	return DisjO(
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
	return ExistO(func(c *rune) Goal {
		return EqualO(r, CharPtr(c))
	})
}

func IsOr(r *Regex) Goal {
	return ExistO(func(r1 *Regex) Goal {
		return ExistO(func(r2 *Regex) Goal {
			return EqualO(r, Or(r1, r2))
		})
	})
}

func IsConcat(r *Regex) Goal {
	return ExistO(func(r1 *Regex) Goal {
		return ExistO(func(r2 *Regex) Goal {
			return EqualO(r, Concat(r1, r2))
		})
	})
}

func IsStar(r *Regex) Goal {
	return ExistO(func(r1 *Regex) Goal {
		return EqualO(r, Star(r1))
	})
}
