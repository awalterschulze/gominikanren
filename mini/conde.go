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
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			g1s := g1()(s)
			return micro.Bind(g1s, g2)
		}
	}
}

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
		return func(s *micro.State) micro.StreamOfStates {
			g1s := g1()(s)
			g2s := g2()(s)
			return micro.Mplus(g1s, g2s)
		}
	}
}

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
