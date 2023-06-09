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
		return gs[0]
	}
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		sss := toStreams(ctx, gs, s)
		Mplus(ctx, sss, ss)
	}
}

func toStreams(ctx context.Context, gs []comicro.Goal, s *comicro.State) []comicro.StreamOfStates {
	ss := make([]comicro.StreamOfStates, len(gs))
	for i := range gs {
		ss[i] = comicro.NewStreamForGoal(ctx, gs[i], s)
	}
	return ss
}

func Mplus(ctx context.Context, sss []comicro.StreamOfStates, ss comicro.StreamOfStates) {
	sss = filterNulls(sss)
	if len(sss) == 0 {
		return
	}
	wait := sync.WaitGroup{}
	wait.Add(len(sss))
	for _, s1 := range sss {
		f := func(s11 comicro.StreamOfStates) func() {
			return func() {
				defer wait.Done()
				comicro.WriteStreamTo(ctx, s11, ss)
			}
		}(s1)
		go f()
	}
	wait.Wait()
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
