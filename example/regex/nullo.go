package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
)

func NullO(r, out *Regex) comicro.Goal {
	return comini.Disjs(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptyStr()),
		),
		comicro.Conj(
			comicro.Exists(func(char *rune) comicro.Goal {
				return comicro.EqualO(r, CharPtr(char))
			}),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Exists(func(a *Regex) comicro.Goal {
			return comicro.Exists(func(b *Regex) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Or(a, b)),
					comini.Disjs(
						comini.Conjs(
							NullO(a, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.Conjs(
							NullO(b, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
						comini.Conjs(
							NullO(a, EmptySet()),
							NullO(b, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
					),
				)
			})
		}),
		comicro.Exists(func(a *Regex) comicro.Goal {
			return comicro.Exists(func(b *Regex) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(r, Concat(a, b)),
					comini.Disjs(
						comini.Conjs(
							NullO(a, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.Conjs(
							NullO(b, EmptySet()),
							comicro.EqualO(out, EmptySet()),
						),
						comini.Conjs(
							NullO(a, EmptyStr()),
							NullO(b, EmptyStr()),
							comicro.EqualO(out, EmptyStr()),
						),
					),
				)
			})
		}),
		comicro.Exists(func(a *Regex) comicro.Goal {
			return comicro.Conj(
				comicro.EqualO(r, Star(a)),
				comicro.EqualO(out, EmptyStr()),
			)
		}),
	)
}
