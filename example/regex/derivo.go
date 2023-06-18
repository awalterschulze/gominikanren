package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
)

func DerivO(r *Regex, char *rune, out *Regex) comicro.Goal {
	return comini.Disjs(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptySet()),
		),
		DeriveCharO(r, char, out),
		comicro.CallFresh(&Regex{}, func(a *Regex) comicro.Goal {
			return comicro.CallFresh(&Regex{}, func(da *Regex) comicro.Goal {
				return comicro.CallFresh(&Regex{}, func(b *Regex) comicro.Goal {
					return comicro.CallFresh(&Regex{}, func(db *Regex) comicro.Goal {
						return comini.Conjs(
							comicro.EqualO(r, Or(a, b)),
							DerivO(a, char, da),
							DerivO(b, char, db),
							comicro.EqualO(out, Or(da, db)),
						)
					})
				})
			})
		}),
		comicro.CallFresh(&Regex{}, func(a *Regex) comicro.Goal {
			return comicro.CallFresh(&Regex{}, func(da *Regex) comicro.Goal {
				return comicro.CallFresh(&Regex{}, func(na *Regex) comicro.Goal {
					return comicro.CallFresh(&Regex{}, func(b *Regex) comicro.Goal {
						return comicro.CallFresh(&Regex{}, func(db *Regex) comicro.Goal {
							return comini.Conjs(
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
		comicro.CallFresh(&Regex{}, func(a *Regex) comicro.Goal {
			return comicro.CallFresh(&Regex{}, func(da *Regex) comicro.Goal {
				return comini.Conjs(
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
	return comini.Disjs(
		comini.Conjs(
			comicro.EqualO(char, &a),
			comicro.EqualO(r, Char('a')),
			comicro.EqualO(out, EmptyStr()),
		),
		comini.Conjs(
			comicro.EqualO(char, &a),
			comicro.EqualO(r, Char('b')),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, &a),
			comicro.EqualO(r, Char('c')),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, &b),
			comicro.EqualO(r, Char('a')),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, &b),
			comicro.EqualO(r, Char('b')),
			comicro.EqualO(out, EmptyStr()),
		),
		comini.Conjs(
			comicro.EqualO(char, &b),
			comicro.EqualO(r, Char('c')),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, &c),
			comicro.EqualO(r, Char('a')),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, &c),
			comicro.EqualO(r, Char('b')),
			comicro.EqualO(out, EmptySet()),
		),
		comini.Conjs(
			comicro.EqualO(char, &c),
			comicro.EqualO(r, Char('c')),
			comicro.EqualO(out, EmptyStr()),
		),
	)
}
