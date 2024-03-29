package micro

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
	return func(s *State) *StreamOfStates {
		g1s := g1(s)
		g2s := g2(s)
		return Mplus(g1s, g2s)
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
*/
func Mplus(s1, s2 *StreamOfStates) *StreamOfStates {
	if s1 == nil {
		return s2
	}
	car, cdr := s1.CarCdr()
	if car != nil { // not a suspension => procedure? == false
		return NewStream(car, func() *StreamOfStates {
			return Mplus(cdr, s2)
		})
	}
	return Suspension(func() *StreamOfStates {
		return Mplus(s2, cdr)
	})
}
