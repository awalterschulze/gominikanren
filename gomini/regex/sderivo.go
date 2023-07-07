package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func SDerivO(r *Regex, char *rune, out *Regex) Goal {
	return DisjO(
		ConjO(
			EqualO(r, EmptySet()),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(r, EmptyStr()),
			EqualO(out, EmptySet()),
		),
		DeriveCharO(r, char, out),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(da *Regex) Goal {
				return ExistO(func(b *Regex) Goal {
					return ExistO(func(db *Regex) Goal {
						return ConjO(
							EqualO(r, Or(a, b)),
							SDerivO(a, char, da),
							SDerivO(b, char, db),
							SimpleOrO(da, db, out),
						)
					})
				})
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(da *Regex) Goal {
				return ExistO(func(na *Regex) Goal {
					return ExistO(func(b *Regex) Goal {
						return ExistO(func(db *Regex) Goal {
							return ExistO(func(ca *Regex) Goal {
								return ExistO(func(cb *Regex) Goal {
									return ConjO(
										EqualO(r, Concat(a, b)),
										SDerivO(a, char, da),
										SDerivO(b, char, db),
										NullO(a, na),
										SimpleConcatO(da, b, ca),
										SimpleConcatO(na, db, cb),
										SimpleOrO(ca, cb, out),
									)
								})
							})
						})
					})
				})
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(da *Regex) Goal {
				return ConjO(
					EqualO(r, Star(a)),
					SDerivO(a, char, da),
					SimpleConcatO(da, r, out),
				)
			})
		}),
	)
}
