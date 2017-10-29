package micro

/*
(define (disj g1 g2)
	(lambda_g (s/c)
		(mplus (g1 s/c) (g2 s/c))
	)
)
*/
func DisjointO(gs ...Goal) Goal {
	if len(gs) == 0 {
		return FailureO()
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := DisjointO(gs[1:]...)
	return func(s *State) StreamOfSubstitutions {
		g1s := g1(s)
		g2s := g2(s)
		return mplus(g1s, g2s)
	}
}

/*
(define (mplus $1 $2)
	(cond
		((null? $1) $2)
		((procedure? $1) (λ_$ () (mplus $2 ($1))))
		(else (cons (car $1) (mplus (cdr $1) $2)))
	)
)
*/
func mplus(g1, g2 StreamOfSubstitutions) StreamOfSubstitutions {
	if g1 == nil {
		return g2
	}
	car, cdr := g1()
	if car != nil { // not a suspension => procedure? == false
		return func() (*State, StreamOfSubstitutions) {
			return car, mplus(cdr, g2)
		}
	}
	return Suspension(func() StreamOfSubstitutions {
		return mplus(g2, cdr)
	})
}
