package micro

/*
ConjunctionO is a goal that returns a logical AND of the input goals.

scheme code:

	(define (conj g1 g2)
		(lambda_g (s/c)
			(bind (g1 s/c) g2)
		)
	)
*/
func ConjunctionO(gs ...Goal) Goal {
	if len(gs) == 0 {
		return SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := ConjunctionO(gs[1:]...)
	return func() GoalFn {
		return func(s *State) StreamOfStates {
			g1s := g1()(s)
			return Bind(g1s, g2)
		}
	}
}

/*
Bind is the monad bind function for goals.

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
