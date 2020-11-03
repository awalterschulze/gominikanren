package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
)

/*
Conde is a disjunction of conjunctions

(define- syntax conde
(syntax- rules ()
    ((_ (g0 g . . . ) . . . ) (disj+ (conj+ g0 g . . . ) . . . ))))
*/
func Conde(gs ...[]micro.Goal) micro.Goal {
	conj := make([]micro.Goal, len(gs))
	for i, v := range gs {
		conj[i] = ConjPlus(v...)
	}
	return DisjPlus(conj...)
}
