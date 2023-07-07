package regex

import (
	. "github.com/awalterschulze/gominikanren/gomini"
)

func NullO(r, out *Regex) Goal {
	return DisjO(
		ConjO(
			EqualO(r, EmptySet()),
			EqualO(out, EmptySet()),
		),
		ConjO(
			EqualO(r, EmptyStr()),
			EqualO(out, EmptyStr()),
		),
		ConjO(
			ExistO(func(char *rune) Goal {
				return EqualO(r, CharPtr(char))
			}),
			EqualO(out, EmptySet()),
		),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(b *Regex) Goal {
				return ConjO(
					EqualO(r, Or(a, b)),
					DisjO(
						ConjO(
							NullO(a, EmptyStr()),
							EqualO(out, EmptyStr()),
						),
						ConjO(
							NullO(b, EmptyStr()),
							EqualO(out, EmptyStr()),
						),
						ConjO(
							NullO(a, EmptySet()),
							NullO(b, EmptySet()),
							EqualO(out, EmptySet()),
						),
					),
				)
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ExistO(func(b *Regex) Goal {
				return ConjO(
					EqualO(r, Concat(a, b)),
					DisjO(
						ConjO(
							NullO(a, EmptySet()),
							EqualO(out, EmptySet()),
						),
						ConjO(
							NullO(b, EmptySet()),
							EqualO(out, EmptySet()),
						),
						ConjO(
							NullO(a, EmptyStr()),
							NullO(b, EmptyStr()),
							EqualO(out, EmptyStr()),
						),
					),
				)
			})
		}),
		ExistO(func(a *Regex) Goal {
			return ConjO(
				EqualO(r, Star(a)),
				EqualO(out, EmptyStr()),
			)
		}),
	)
}
