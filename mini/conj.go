package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
)

/*
ConjPlus is a macro that extends conjunction to arbitrary arguments

(define- syntax conj+
(syntax- rules ()
    ((_ g) (Zzz g))
    ((_ g0 g . . . ) (conj (Zzz g0) (conj+ g . . . )))))
*/
func ConjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.SuccessO
	}
	if len(gs) == 1 {
		return micro.Zzz(gs[0])
	}
	g1 := micro.Zzz(gs[0])
	g2 := ConjPlus(gs[1:]...)
	return func(s *micro.State) *micro.StreamOfStates {
		g1s := g1(s)
		return micro.Bind(g1s, g2)
	}
}

// ConjPlusNoZzz is the conj+ macro without wrapping
// its goals in a suspension (zzz)
func ConjPlusNoZzz(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := ConjPlusNoZzz(gs[1:]...)
	return func(s *micro.State) *micro.StreamOfStates {
		g1s := g1(s)
		return micro.Bind(g1s, g2)
	}
}
