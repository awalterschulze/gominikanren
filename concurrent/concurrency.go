package concurrent

import (
	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
)

type answer struct {
	i int
	s micro.StreamOfStates
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
}

/*
ConjPlus is a macro that extends disjunction to arbitrary arguments
It resolves its arguments in parallel and binds them at the end
Can fail early if one of the goals returns nil, signalling failure
Notably, unlike mini.DisjPlus, it does _not_ wrap its goals in Zzz
*/
func ConjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			ch := make(chan answer)
			go func() {
				for i, g := range gs {
					go func(index int, goal micro.Goal) {
						ss := goal()(s)
						ch <- answer{s: ss}
					}(i, g)
				}
			}()
			ch2 := make(chan micro.StreamOfStates)
			go func() {
				g1s := gs[0]()(s)
				g2 := mini.ConjPlus(gs[1:]...)
				ch2 <- micro.Bind(g1s, g2)
			}()
			for {
				select {
				case ans := <-ch:
					if ans.s == nil {
						return nil
					}
				case stream := <-ch2:
					return stream
				}
			}
		}
	}
}
