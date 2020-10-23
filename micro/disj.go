package micro

/*
DisjointO is a goal that returns a logical OR of the input goals.

scheme code:

	(define (disj g1 g2)
		(lambda_g (s/c)
			(mplus (g1 s/c) (g2 s/c))
		)
	)
*/
func DisjointO(g1, g2 Goal) Goal {
	return func() GoalFn {
		return func(s *State) StreamOfStates {
			g1s := g1()(s)
			g2s := g2()(s)
			return Mplus(g1s, g2s)
		}
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
func Mplus(g1, g2 StreamOfStates) StreamOfStates {
	if g1 == nil {
		return g2
	}
	car, cdr := g1()
	if car != nil { // not a suspension => procedure? == false
		return func() (*State, StreamOfStates) {
			return car, Mplus(cdr, g2)
		}
	}
	return Suspension(func() StreamOfStates {
		return Mplus(g2, cdr)
	})
}
