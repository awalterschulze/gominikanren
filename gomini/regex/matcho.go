package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func MatchO(r *Regex, s *String, res *Regex) Goal {
	return ExistO(func(dr *Regex) Goal {
		return ConjO(
			SDerivOs(r, s, dr),
			NullO(dr, res),
		)
	})
}

func IsMatchO(r *Regex, s *String) Goal {
	return ExistO(func(dr *Regex) Goal {
		return ConjO(
			SDerivOs(r, s, dr),
			IsNullO(dr),
		)
	})
}

func SDerivOs(r *Regex, s *String, res *Regex) Goal {
	return DisjO(
		ConjO(
			EqualO(s, nil),
			EqualO(r, res),
		),
		ExistO(func(head *rune) Goal {
			return ExistO(func(tail *String) Goal {
				return ExistO(func(dr *Regex) Goal {
					return ConjO(
						EqualO(s, &String{head, tail}),
						SDerivO(r, head, dr),
						SDerivOs(dr, tail, res),
					)
				})
			})
		}),
	)
}
