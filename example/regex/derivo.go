package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
)

func DerivO(r *Regex, char *rune, out *Regex) comicro.Goal {
	return comicro.Disjs(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptySet()),
		),
		DeriveCharO(r, char, out),
		comicro.Exists(func(a *Regex) comicro.Goal {
			return comicro.Exists(func(da *Regex) comicro.Goal {
				return comicro.Exists(func(b *Regex) comicro.Goal {
					return comicro.Exists(func(db *Regex) comicro.Goal {
						return comicro.Conjs(
							comicro.EqualO(r, Or(a, b)),
							DerivO(a, char, da),
							DerivO(b, char, db),
							comicro.EqualO(out, Or(da, db)),
						)
					})
				})
			})
		}),
		comicro.Exists(func(a *Regex) comicro.Goal {
			return comicro.Exists(func(da *Regex) comicro.Goal {
				return comicro.Exists(func(na *Regex) comicro.Goal {
					return comicro.Exists(func(b *Regex) comicro.Goal {
						return comicro.Exists(func(db *Regex) comicro.Goal {
							return comicro.Conjs(
								comicro.EqualO(r, Concat(a, b)),
								DerivO(a, char, da),
								DerivO(b, char, db),
								NullO(a, na),
								comicro.EqualO(out, Or(Concat(da, b), Concat(na, db))),
							)
						})
					})
				})
			})
		}),
		comicro.Exists(func(a *Regex) comicro.Goal {
			return comicro.Exists(func(da *Regex) comicro.Goal {
				return comicro.Conjs(
					comicro.EqualO(r, Star(a)),
					DerivO(a, char, da),
					comicro.EqualO(out, Concat(da, Star(a))),
				)
			})
		}),
	)
}

func DeriveCharO(r *Regex, char *rune, out *Regex) comicro.Goal {
	a := rune('a')
	b := rune('b')
	c := rune('c')
	return comicro.Disjs(
		comicro.Conjs(
			comicro.EqualO(char, &a),
			comicro.EqualO(r, Char('a')),
			comicro.EqualO(out, EmptyStr()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &a),
			comicro.EqualO(r, Char('b')),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &a),
			comicro.EqualO(r, Char('c')),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &b),
			comicro.EqualO(r, Char('a')),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &b),
			comicro.EqualO(r, Char('b')),
			comicro.EqualO(out, EmptyStr()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &b),
			comicro.EqualO(r, Char('c')),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &c),
			comicro.EqualO(r, Char('a')),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &c),
			comicro.EqualO(r, Char('b')),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conjs(
			comicro.EqualO(char, &c),
			comicro.EqualO(r, Char('c')),
			comicro.EqualO(out, EmptyStr()),
		),
	)
}
