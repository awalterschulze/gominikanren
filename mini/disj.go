package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
)

/*
DisjPlus is a macro that extends disjunction to arbitrary arguments

(define- syntax disj+
(syntax-rules ()
    ((_ g) (Zzz g))
    ((_ g0 g . . . ) (disj (Zzz g0) (disj+ g . . . )))))
*/
func DisjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return micro.Zzz(gs[0])
	}
	g1 := micro.Zzz(gs[0])
	g2 := DisjPlus(gs[1:]...)
	return func() micro.GoalFn {
		return func(s *micro.State) *micro.StreamOfStates {
			g1s := g1()(s)
			g2s := g2()(s)
			return micro.Mplus(g1s, g2s)
		}
	}
}

// DisjPlusNoZzz is the disj+ macro without wrapping
// its goals in a suspension (zzz)
func DisjPlusNoZzz(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := DisjPlusNoZzz(gs[1:]...)
	return func() micro.GoalFn {
		return func(s *micro.State) *micro.StreamOfStates {
			g1s := g1()(s)
			g2s := g2()(s)
			return micro.Mplus(g1s, g2s)
		}
	}
}
