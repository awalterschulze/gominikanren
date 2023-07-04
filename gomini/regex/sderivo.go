package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func SDerivO(r *Regex, char *rune, out *Regex) Goal {
	return Disjs(
		Conj(
			EqualO(r, EmptySet()),
			EqualO(out, EmptySet()),
		),
		Conj(
			EqualO(r, EmptyStr()),
			EqualO(out, EmptySet()),
		),
		DeriveCharO(r, char, out),
		Exists(func(a *Regex) Goal {
			return Exists(func(da *Regex) Goal {
				return Exists(func(b *Regex) Goal {
					return Exists(func(db *Regex) Goal {
						return Conjs(
							EqualO(r, Or(a, b)),
							SDerivO(a, char, da),
							SDerivO(b, char, db),
							SimpleOrO(da, db, out),
						)
					})
				})
			})
		}),
		Exists(func(a *Regex) Goal {
			return Exists(func(da *Regex) Goal {
				return Exists(func(na *Regex) Goal {
					return Exists(func(b *Regex) Goal {
						return Exists(func(db *Regex) Goal {
							return Exists(func(ca *Regex) Goal {
								return Exists(func(cb *Regex) Goal {
									return Conjs(
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
		Exists(func(a *Regex) Goal {
			return Exists(func(da *Regex) Goal {
				return Conjs(
					EqualO(r, Star(a)),
					SDerivO(a, char, da),
					SimpleConcatO(da, r, out),
				)
			})
		}),
	)
}
