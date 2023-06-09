package comicro

import (
	"context"
	"sync"
)

/*
Disj is a goal that returns a logical OR of the input goals.

scheme code:

	(define (disj g1 g2)
		(lambda_g (s/c)
			(mplus (g1 s/c) (g2 s/c))
		)
	)
*/
func Disj(g1, g2 Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		g1s := NewStreamForGoal(ctx, g1, s)
		g2s := NewStreamForGoal(ctx, g2, s)
		Mplus(ctx, g1s, g2s, ss)
	}
}

/*
Mplus is responsible for merging streams

scheme code:

	(define (mplus $1 $2)
		(cond
			((null? $1) $2)
			((procedure? $1) (λ_$ () (mplus $2 ($1))))
			(else (cons (car $1) (mplus (cdr $1) $2)))
		)
	)

An alternative implementation could swap the goals for a more fair distribution:

scheme code:

	(define (mplus $1 $2)
		(cond
			((null? $1) $2)
			((procedure? $1) (λ_$ () (mplus $2 ($1))))
			(else (cons (car $1) (mplus $2 (cdr $1))))
		)
	)

not a suspension => procedure? == false
*/
func Mplus(ctx context.Context, s1, s2, res StreamOfStates) {
	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		WriteStreamTo(ctx, s1, res)
	}()
	go func() {
		defer wait.Done()
		WriteStreamTo(ctx, s2, res)
	}()
	wait.Wait()
}
