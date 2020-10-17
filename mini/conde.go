package mini

import (
    "github.com/awalterschulze/gominikanren/micro"
)

/*
(define- syntax disj+
(syntax-rules ()
    ((_ g) (Zzz g))
    ((_ g0 g . . . ) (disj (Zzz g0) (disj+ g . . . )))))
*/
func DisjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO()
	}
	if len(gs) == 1 {
		return Zzz(gs[0])
	}
	g1 := Zzz(gs[0])
	g2 := DisjPlus(gs[1:]...)
	return func(s *micro.State) micro.StreamOfStates {
		g1s := g1(s)
		g2s := g2(s)
		return micro.Mplus(g1s, g2s)
	}
}

// TODO: should be list of lists of goals
func Conde(gs ...micro.Goal) micro.Goal {
	return DisjPlus(gs...)
}

func Zzz(g micro.Goal) micro.Goal {
	return func(s *micro.State) micro.StreamOfStates {
		return micro.Suspension(func() micro.StreamOfStates {
            return g(s)
        })
    }
}

func Zzz2(g func()micro.Goal) micro.Goal {
	return func(s *micro.State) micro.StreamOfStates {
		return micro.Suspension(func() micro.StreamOfStates {
            return g()(s)
        })
    }
}
