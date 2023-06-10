package comini

import (
	"github.com/awalterschulze/gominikanren/comicro"
)

// Conde is a disjunction of conjunctions
func Conde(gs ...[]comicro.Goal) comicro.Goal {
	conj := make([]comicro.Goal, len(gs))
	for i, v := range gs {
		conj[i] = Conjs(v...)
	}
	return Disjs(conj...)
}
