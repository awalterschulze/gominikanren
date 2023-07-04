package regex

import (
	. "github.com/awalterschulze/gominikanren/comicro"
)

func DerivO(r *Regex, char *rune, out *Regex) Goal {
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
							DerivO(a, char, da),
							DerivO(b, char, db),
							EqualO(out, Or(da, db)),
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
							return Conjs(
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
		Exists(func(a *Regex) Goal {
			return Exists(func(da *Regex) Goal {
				return Conjs(
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
	return Disjs(
		Conjs(
			EqualO(char, &a),
			EqualO(r, Char('a')),
			EqualO(out, EmptyStr()),
		),
		Conjs(
			EqualO(char, &a),
			EqualO(r, Char('b')),
			EqualO(out, EmptySet()),
		),
		Conjs(
			EqualO(char, &a),
			EqualO(r, Char('c')),
			EqualO(out, EmptySet()),
		),
		Conjs(
			EqualO(char, &b),
			EqualO(r, Char('a')),
			EqualO(out, EmptySet()),
		),
		Conjs(
			EqualO(char, &b),
			EqualO(r, Char('b')),
			EqualO(out, EmptyStr()),
		),
		Conjs(
			EqualO(char, &b),
			EqualO(r, Char('c')),
			EqualO(out, EmptySet()),
		),
		Conjs(
			EqualO(char, &c),
			EqualO(r, Char('a')),
			EqualO(out, EmptySet()),
		),
		Conjs(
			EqualO(char, &c),
			EqualO(r, Char('b')),
			EqualO(out, EmptySet()),
		),
		Conjs(
			EqualO(char, &c),
			EqualO(r, Char('c')),
			EqualO(out, EmptyStr()),
		),
	)
}
