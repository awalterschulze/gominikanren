package comini

import (
	"context"
	"sync"

	"github.com/awalterschulze/gominikanren/comicro"
)

/*
DisjPlus is a macro that extends disjunction to arbitrary arguments

(define- syntax disj+
(syntax-rules ()

	((_ g) (Zzz g))
	((_ g0 g . . . ) (disj (Zzz g0) (disj+ g . . . )))))
*/
func DisjPlus(gs ...comicro.Goal) comicro.Goal {
	if len(gs) == 0 {
		return comicro.FailureO
	}
	if len(gs) == 1 {
		return comicro.Zzz(gs[0])
	}
	return func(ctx context.Context, s *comicro.State) comicro.StreamOfStates {
		ss := toStreams(ctx, gs, s)
		return Mplus(ctx, ss)
	}
}

func toStreams(ctx context.Context, gs []comicro.Goal, s *comicro.State) []comicro.StreamOfStates {
	ss := make([]comicro.StreamOfStates, len(gs))
	for i, g := range gs {
		ss[i] = g(ctx, s)
	}
	return ss
}

// DisjPlusNoZzz is the disj+ macro without wrapping
// its goals in a suspension (zzz)
func DisjPlusNoZzz(gs ...comicro.Goal) comicro.Goal {
	if len(gs) == 0 {
		return comicro.FailureO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	return func(ctx context.Context, s *comicro.State) comicro.StreamOfStates {
		ss := toStreams(ctx, gs, s)
		return Mplus(ctx, ss)
	}
}

func Mplus(ctx context.Context, ss []comicro.StreamOfStates) comicro.StreamOfStates {
	ss = filterNulls(ss)
	if len(ss) == 0 {
		return nil
	}
	if len(ss) == 1 {
		return ss[0]
	}
	s := comicro.NewEmptyStream()
	go func() {
		defer close(s)
		wait := sync.WaitGroup{}
		wait.Add(len(ss))
		for _, s1 := range ss {
			f := func(s11 comicro.StreamOfStates) func() {
				return func() {
					defer wait.Done()
					comicro.WriteStreamTo(ctx, s11, s)
				}
			}(s1)
			go f()
		}
		wait.Wait()
	}()
	return s
}

func filterNulls(ss []comicro.StreamOfStates) []comicro.StreamOfStates {
	res := []comicro.StreamOfStates{}
	for _, s := range ss {
		if s != nil {
			res = append(res, s)
		}
	}
	return res
}
