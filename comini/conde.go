package comini

import (
	"github.com/awalterschulze/gominikanren/comicro"
)

/*
Conde is a disjunction of conjunctions

(define- syntax conde
(syntax- rules ()

	((_ (g0 g . . . ) . . . ) (disj+ (conj+ g0 g . . . ) . . . ))))
*/
func Conde(gs ...[]comicro.Goal) comicro.Goal {
	conj := make([]comicro.Goal, len(gs))
	for i, v := range gs {
		conj[i] = ConjPlus(v...)
	}
	return DisjPlus(conj...)
}
