package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func DerivO(r *Regex, char *rune, dr *Regex) Goal {
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
							DerivO(a, char, da),
							DerivO(b, char, db),
							EqualO(dr, Or(da, db)),
						)
					})
				})
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(b *Regex) Goal {
				return ExistO(func(da *Regex) Goal {
					return ExistO(func(db *Regex) Goal {
						return ConjO(
							EqualO(r, Concat(a, b)),
							DerivO(a, char, da),
							DerivO(b, char, db),
							DisjO(
								ConjO(IsNullO(a), EqualO(dr, Or(Concat(da, b), db))),
								EqualO(dr, Concat(da, b)),
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
					DerivO(a, char, da),
					EqualO(dr, Concat(da, Star(a))),
				)
			})
		}),
	)
}

func DeriveCharO(r *Regex, char *rune, dr *Regex) Goal {
	a := rune('a')
	b := rune('b')
	return DisjO(
		ConjO(
			EqualO(char, &a),
			EqualO(r, Char('a')),
			EqualO(dr, EmptyStr()),
		),
		ConjO(
			EqualO(char, &a),
			EqualO(r, Char('b')),
			EqualO(dr, EmptySet()),
		),
		ConjO(
			EqualO(char, &b),
			EqualO(r, Char('a')),
			EqualO(dr, EmptySet()),
		),
		ConjO(
			EqualO(char, &b),
			EqualO(r, Char('b')),
			EqualO(dr, EmptyStr()),
		),
	)
}
