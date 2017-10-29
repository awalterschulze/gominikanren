package micro

/*
(define (conj g1 g2)
	(lambda_g (s/c)
		(bind (g1 s/c) g2)
	)
)
*/
func ConjunctionO(gs ...Goal) Goal {
	if len(gs) == 0 {
		return SuccessO()
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := ConjunctionO(gs[1:]...)
	return func(s *State) StreamOfSubstitutions {
		g1s := g1(s)
		return bind(g1s, g2)
	}
}

/*
(define (bind $ g)
	(cond
		((null? $) mzero)
		((procedure? $) (lambda_$ () (bind ($) g)))
		(else (mplus (g (car $)) (bind (cdr $) g)))
	)
)
*/
func bind(s StreamOfSubstitutions, g Goal) StreamOfSubstitutions {
	if s == nil {
		return nil
	}
	car, cdr := s()
	if car != nil { // not a suspension => procedure? == false
		return mplus(
			g(car),
			bind(
				cdr,
				g,
			),
		)
	}
	return Suspension(func() StreamOfSubstitutions {
		return bind(
			cdr,
			g,
		)
	})
}
