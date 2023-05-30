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
	return func(ctx context.Context, s *State) StreamOfStates {
		g1s := g1(ctx, s)
		g2s := g2(ctx, s)
		return Mplus(ctx, g1s, g2s)
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
func Mplus(ctx context.Context, s1, s2 StreamOfStates) StreamOfStates {
	if s1 == nil {
		return s2
	}
	if s2 == nil {
		return s1
	}
	s := make(chan *State, 0)
	go func() {
		defer close(s)
		wait := sync.WaitGroup{}
		wait.Add(2)
		go func() {
			for a1 := range s1 {
				select {
				case <-ctx.Done():
					return
				case s <- a1:
				}
			}
			wait.Done()
		}()
		go func() {
			for a2 := range s2 {
				select {
				case <-ctx.Done():
					return
				case s <- a2:
				}
			}
			wait.Done()
		}()
		wait.Wait()
	}()
	return s
}
