package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
)

func MatchO(r *Regex, s *string, res *Regex) comicro.Goal {
	if s == nil {
		return NullO(r, res)
	}
	return comicro.CallFresh(&Regex{}, func(dr *Regex) comicro.Goal {
		return comicro.Conj(
			SDerivOs(r, *s, dr),
			NullO(dr, res),
		)
	})
}

func SDerivOs(r *Regex, s string, res *Regex) comicro.Goal {
	ss := []rune(s)
	if len(ss) == 0 {
		return comicro.EqualO(res, r)
	}
	if len(ss) == 1 {
		c := ss[0]
		return SDerivO(r, &c, res)
	}
	head := ss[0]
	tail := string(ss[1:])
	return comicro.CallFresh(&Regex{}, func(dr *Regex) comicro.Goal {
		return comicro.Conj(
			SDerivO(r, &head, dr),
			SDerivOs(dr, tail, res),
		)
	})
}
