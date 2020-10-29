package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
)

type answer struct {
	i int
	s micro.StreamOfStates
}

/*
ConcurrentDisjPlus is a macro that extends disjunction to arbitrary arguments
It resolves its arguments in parallel and merges them at the end
Notably, unlike DisjPlus, it does _not_ wrap its goals in Zzz
*/
func ConcurrentDisjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			list := make([]micro.StreamOfStates, len(gs))
			ch := make(chan answer)
			for i, g := range gs {
				go func(index int, goal micro.Goal) {
					ss := goal()(s)
					ch <- answer{i: index, s: ss}
				}(i, g)
			}
			for i := 0; i < len(gs); i++ {
				ans := <-ch
				list[ans.i] = ans.s
			}
			stream := micro.Mplus(list[len(gs)-2], list[len(gs)-1])
			if len(gs) == 2 {
				return stream
			}
			for i := len(gs) - 3; i >= 0; i-- {
				stream = micro.Mplus(list[i], stream)
			}
			return stream
		}
	}
}
