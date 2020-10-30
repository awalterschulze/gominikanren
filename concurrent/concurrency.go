package concurrent

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
func DisjPlus(gs ...micro.Goal) micro.Goal {
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
            for _ = range gs {
				ans := <-ch
				list[ans.i] = ans.s
			}
            stream := list[len(gs)-1]
			for i := len(gs) - 2; i >= 0; i-- {
				stream = micro.Mplus(list[i], stream)
			}
			return stream
		}
	}
}
