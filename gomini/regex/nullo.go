package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func NullO(r, out *Regex) Goal {
	return Disjs(
		Conj(
			EqualO(r, EmptySet()),
			EqualO(out, EmptySet()),
		),
		Conj(
			EqualO(r, EmptyStr()),
			EqualO(out, EmptyStr()),
		),
		Conj(
			Exists(func(char *rune) Goal {
				return EqualO(r, CharPtr(char))
			}),
			EqualO(out, EmptySet()),
		),
		Exists(func(a *Regex) Goal {
			return Exists(func(b *Regex) Goal {
				return Conj(
					EqualO(r, Or(a, b)),
					Disjs(
						Conjs(
							NullO(a, EmptyStr()),
							EqualO(out, EmptyStr()),
						),
						Conjs(
							NullO(b, EmptyStr()),
							EqualO(out, EmptyStr()),
						),
						Conjs(
							NullO(a, EmptySet()),
							NullO(b, EmptySet()),
							EqualO(out, EmptySet()),
						),
					),
				)
			})
		}),
		Exists(func(a *Regex) Goal {
			return Exists(func(b *Regex) Goal {
				return Conj(
					EqualO(r, Concat(a, b)),
					Disjs(
						Conjs(
							NullO(a, EmptySet()),
							EqualO(out, EmptySet()),
						),
						Conjs(
							NullO(b, EmptySet()),
							EqualO(out, EmptySet()),
						),
						Conjs(
							NullO(a, EmptyStr()),
							NullO(b, EmptyStr()),
							EqualO(out, EmptyStr()),
						),
					),
				)
			})
		}),
		Exists(func(a *Regex) Goal {
			return Conj(
				EqualO(r, Star(a)),
				EqualO(out, EmptyStr()),
			)
		}),
	)
}
