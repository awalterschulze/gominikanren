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
func Bind(ctx context.Context, stream StreamOfStates, g Goal) StreamOfStates {
	if stream == nil {
		return nil
	}
	state, ok := stream.Read(ctx)
	var rest StreamOfStates = nil
	if ok {
		rest = stream
	}
	if state != nil { // not a suspension => procedure? == false
		return Mplus(ctx, g(ctx, state), Bind(ctx, rest, g))
	}
	return Suspension(ctx, func() StreamOfStates {
		return Bind(ctx, rest, g)
	})
}
