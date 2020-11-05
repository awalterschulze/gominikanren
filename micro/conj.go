package micro

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
	return func() GoalFn {
		return func(s *State) StreamOfStates {
			g1s := g1()(s)
			return Bind(g1s, g2)
		}
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
*/
func Bind(s StreamOfStates, g Goal) StreamOfStates {
	if s == nil {
		return nil
	}
	car, cdr := s()
	if car != nil { // not a suspension => procedure? == false
		return Mplus(
			g()(car),
			Bind(
				cdr,
				g,
			),
		)
	}
	return Suspension(func() StreamOfStates {
		return Bind(
			cdr,
			g,
		)
	})
}
