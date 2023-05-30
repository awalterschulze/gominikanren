package comicro

import "context"

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
	return func(ctx context.Context, s *State) StreamOfStates {
		g1s := g1(ctx, s)
		return Bind(ctx, g1s, g2)
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
func Bind(ctx context.Context, s StreamOfStates, g Goal) StreamOfStates {
	if s == nil {
		return nil
	}
	car, cdr := s.CarCdr()
	if car != nil { // not a suspension => procedure? == false
		return Mplus(ctx, g(ctx, car), Bind(ctx, cdr, g))
	}
	return Suspension(ctx, func() StreamOfStates {
		return Bind(ctx, cdr, g)
	})
}
