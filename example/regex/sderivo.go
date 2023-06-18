package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
)

func SDerivO(r *Regex, char *rune, out *Regex) comicro.Goal {
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
							SDerivO(a, char, da),
							SDerivO(b, char, db),
							SimpleOrO(da, db, out),
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
							return comicro.CallFresh(&Regex{}, func(ca *Regex) comicro.Goal {
								return comicro.CallFresh(&Regex{}, func(cb *Regex) comicro.Goal {
									return comini.Conjs(
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
		comicro.CallFresh(&Regex{}, func(a *Regex) comicro.Goal {
			return comicro.CallFresh(&Regex{}, func(da *Regex) comicro.Goal {
				return comini.Conjs(
					comicro.EqualO(r, Star(a)),
					SDerivO(a, char, da),
					SimpleConcatO(da, r, out),
				)
			})
		}),
	)
}
