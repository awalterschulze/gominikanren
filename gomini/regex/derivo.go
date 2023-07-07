package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func DerivO(r *Regex, char *rune, out *Regex) Goal {
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
							DerivO(a, char, da),
							DerivO(b, char, db),
							EqualO(out, Or(da, db)),
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
							return ConjO(
								EqualO(r, Concat(a, b)),
								DerivO(a, char, da),
								DerivO(b, char, db),
								NullO(a, na),
								EqualO(out, Or(Concat(da, b), Concat(na, db))),
							)
						})
					})
				})
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(da *Regex) Goal {
				return ConjO(
					EqualO(r, Star(a)),
					DerivO(a, char, da),
					EqualO(out, Concat(da, Star(a))),
				)
			})
		}),
	)
}

func DeriveCharO(r *Regex, char *rune, out *Regex) Goal {
	a := rune('a')
	b := rune('b')
	c := rune('c')
	return DisjO(
		ConjO(
			EqualO(char, &a),
			EqualO(r, Char('a')),
			EqualO(out, EmptyStr()),
		),
		ConjO(
			EqualO(char, &a),
			EqualO(r, Char('b')),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(char, &a),
			EqualO(r, Char('c')),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(char, &b),
			EqualO(r, Char('a')),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(char, &b),
			EqualO(r, Char('b')),
			EqualO(out, EmptyStr()),
		),
		ConjO(
			EqualO(char, &b),
			EqualO(r, Char('c')),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(char, &c),
			EqualO(r, Char('a')),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(char, &c),
			EqualO(r, Char('b')),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(char, &c),
			EqualO(r, Char('c')),
			EqualO(out, EmptyStr()),
		),
	)
}
