package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func MatchO(r *Regex, s *String, res *Regex) Goal {
	if s == nil {
		return NullO(r, res)
	}
	return ExistO(func(dr *Regex) Goal {
		return ConjO(
			SDerivOs(r, s, dr),
			NullO(dr, res),
		)
	})
}

func IsMatchO(r *Regex, s *String) Goal {
	if s == nil {
		return IsNullO(r)
	}
	return ExistO(func(dr *Regex) Goal {
		return ConjO(
			SDerivOs(r, s, dr),
			IsNullO(dr),
		)
	})
}

func SDerivOs(r *Regex, s *String, res *Regex) Goal {
	if s.Value == nil {
		return EqualO(res, r)
	}
	if s.Next == nil {
		c := s.Value
		return SDerivO(r, c, res)
	}
	head := s.Value
	tail := s.Next
	return ExistO(func(dr *Regex) Goal {
		return ConjO(
			SDerivO(r, head, dr),
			SDerivOs(dr, tail, res),
		)
	})
}
