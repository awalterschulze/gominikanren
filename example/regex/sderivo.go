package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
)

func SDerivO(r *Regex, char *rune, out *Regex) comicro.Goal {
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
							SDerivO(a, char, da),
							SDerivO(b, char, db),
							SimpleOrO(da, db, out),
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
							return comicro.Exists(func(ca *Regex) comicro.Goal {
								return comicro.Exists(func(cb *Regex) comicro.Goal {
									return comicro.Conjs(
										comicro.EqualO(r, Concat(a, b)),
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
		comicro.Exists(func(a *Regex) comicro.Goal {
			return comicro.Exists(func(da *Regex) comicro.Goal {
				return comicro.Conjs(
					comicro.EqualO(r, Star(a)),
					SDerivO(a, char, da),
					SimpleConcatO(da, r, out),
				)
			})
		}),
	)
}
