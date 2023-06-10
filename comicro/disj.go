package comicro

import (
	"context"
	"sync"
)

// Disj is a goal that returns a logical OR of the input goals.
func Disj(g1, g2 Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		wait := sync.WaitGroup{}
		Go(ctx, &wait, func() { g1(ctx, s, ss) })
		Go(ctx, &wait, func() { g2(ctx, s, ss) })
		wait.Wait()
	}
}

// Mplus is responsible for merging streams
func Mplus(ctx context.Context, s1, s2, res StreamOfStates) {
	wait := sync.WaitGroup{}
	Go(ctx, &wait, func() { WriteStreamTo(ctx, s1, res) })
	Go(ctx, &wait, func() { WriteStreamTo(ctx, s2, res) })
	wait.Wait()
}
