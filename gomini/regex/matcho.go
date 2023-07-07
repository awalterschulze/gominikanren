package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func MatchO(r *Regex, s *string, res *Regex) Goal {
	if s == nil {
		return NullO(r, res)
	}
	return ExistO(func(dr *Regex) Goal {
		return ConjO(
			SDerivOs(r, *s, dr),
			NullO(dr, res),
		)
	})
}

func SDerivOs(r *Regex, s string, res *Regex) Goal {
	ss := []rune(s)
	if len(ss) == 0 {
		return EqualO(res, r)
	}
	if len(ss) == 1 {
		c := ss[0]
		return SDerivO(r, &c, res)
	}
	head := ss[0]
	tail := string(ss[1:])
	return ExistO(func(dr *Regex) Goal {
		return ConjO(
			SDerivO(r, &head, dr),
			SDerivOs(dr, tail, res),
		)
	})
}
