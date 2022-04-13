package concurrent

import (
	"github.com/awalterschulze/gominikanren/micro"
)

type answer struct {
	i int
	s *micro.StreamOfStates
}

/*
DisjPlus is a macro that extends disjunction to arbitrary arguments
It resolves its arguments in parallel and merges them at the end
Notably, unlike mini.DisjPlus, it does _not_ wrap its goals in Zzz
*/
func DisjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	return func(s *micro.State) *micro.StreamOfStates {
		list := make([]*micro.StreamOfStates, len(gs))
		ch := make(chan answer)
		for i, g := range gs {
			go func(index int, goal micro.Goal) {
				ss := goal(s)
				ch <- answer{i: index, s: ss}
			}(i, g)
		}
		for range gs {
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

// DisjPlusZzz is concurrent disj+ that _does_ wrap its
// goal arguments in suspension (zzz)
func DisjPlusZzz(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return micro.Zzz(gs[0])
	}
	return func(s *micro.State) *micro.StreamOfStates {
		list := make([]*micro.StreamOfStates, len(gs))
		ch := make(chan answer)
		for i, g := range gs {
			go func(index int, goal micro.Goal) {
				ss := micro.Zzz(goal)(s)
				ch <- answer{i: index, s: ss}
			}(i, g)
		}
		for range gs {
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

// DisjPlusNoOrder throws away order of evaluation of goals
// in order to be faster. It will still evaluate all equally
func DisjPlusNoOrder(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	return func(s *micro.State) *micro.StreamOfStates {
		ch := make(chan *micro.StreamOfStates)
		for _, g := range gs {
			go func(goal micro.Goal) {
				ch <- goal(s)
			}(g)
		}
		stream := <-ch
		for i := 0; i < len(gs)-1; i++ {
			stream = micro.Mplus(<-ch, stream)
		}
		return stream
	}
}
