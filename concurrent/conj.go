package concurrent

import (
	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
)

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
		return func(s *micro.State) *micro.StreamOfStates {
			ch := make(chan answer)
			go func() {
				for i, g := range gs {
					go func(index int, goal micro.Goal) {
						ss := goal()(s)
						ch <- answer{s: ss}
					}(i, g)
				}
			}()
			ch2 := make(chan *micro.StreamOfStates)
			go func() {
				g1s := gs[0]()(s)
				g2 := mini.ConjPlusNoZzz(gs[1:]...)
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

// ConjPlusZzz is concurrent conj+ that _does_ wrap its
// goal arguments in suspension (zzz)
func ConjPlusZzz(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.SuccessO
	}
	if len(gs) == 1 {
		return micro.Zzz(gs[0])
	}
	return func() micro.GoalFn {
		return func(s *micro.State) *micro.StreamOfStates {
			ch := make(chan answer)
			go func() {
				for i, g := range gs {
					go func(index int, goal micro.Goal) {
						ss := micro.Zzz(goal)()(s)
						ch <- answer{s: ss}
					}(i, g)
				}
			}()
			ch2 := make(chan *micro.StreamOfStates)
			go func() {
				g1s := micro.Zzz(gs[0])()(s)
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
