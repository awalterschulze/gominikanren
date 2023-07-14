package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func SDerivO(r *Regex, char *rune, dr *Regex) Goal {
	return DisjO(
		ConjO(
			EqualO(r, EmptySet()),
			EqualO(dr, EmptySet()),
		),
		ConjO(
			EqualO(r, EmptyStr()),
			EqualO(dr, EmptySet()),
		),
		DeriveCharO(r, char, dr),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(da *Regex) Goal {
				return ExistO(func(b *Regex) Goal {
					return ExistO(func(db *Regex) Goal {
						return ConjO(
							EqualO(r, Or(a, b)),
							SDerivO(a, char, da),
							SDerivO(b, char, db),
							SimpleOrO(da, db, dr),
						)
					})
				})
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(da *Regex) Goal {
				return ExistO(func(b *Regex) Goal {
					return ExistO(func(ca *Regex) Goal {
						return ConjO(
							EqualO(r, Concat(a, b)),
							SDerivO(a, char, da),
							SimpleConcatO(da, b, ca),
							DisjO(
								ExistO(func(db *Regex) Goal {
									return ConjO(
										IsNullO(a),
										SDerivO(b, char, db),
										EqualO(dr, Or(ca, db)),
									)
								}),
								EqualO(dr, ca),
							),
						)
					})
				})
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(da *Regex) Goal {
				return ConjO(
					EqualO(r, Star(a)),
					SDerivO(a, char, da),
					SimpleConcatO(da, r, dr),
				)
			})
		}),
	)
}
