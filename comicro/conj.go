package comicro

import (
	"context"
	"sync"
)

/*
Conj is a goal that returns a logical AND of the input goals.

scheme code:

	(define (conj g1 g2)
		(lambda_g (s/c)
			(bind (g1 s/c) g2)
		)
	)
*/
func Conj(g1, g2 Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		g1s := NewStreamForGoal(ctx, g1, s)
		Bind(ctx, g1s, g2, ss)
	}
}

/*
Bind is the monad bind function for state.

scheme code:

	(define (bind $ g)
		(cond
			((null? $) mzero)
			((procedure? $) (lambda_$ () (bind ($) g)))
			(else (mplus (g (car $)) (bind (cdr $) g)))
		)
	)

not a suspension => procedure? == false
*/
func Bind(ctx context.Context, stream StreamOfStates, g Goal, res StreamOfStates) {
	wait := sync.WaitGroup{}
	for {
		state, ok := stream.ReadNonNull(ctx)
		if !ok {
			break
		}
		Go(ctx, &wait, func() { g(ctx, state, res) })
	}
	wait.Wait()
}
