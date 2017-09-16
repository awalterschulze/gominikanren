package micro

/*
(define (disj2 g1 g2)
	(lambda (s)
		(appendinf (g1 s) (g2 s))
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
	return func(s Substitution) StreamOfSubstitutions {
		g1s := g1(s)
		g2s := g2(s)
		return appendInf(g1s, g2s)
	}
}

/*
(define (appendInf s1 s2)
	(cond
		(
			(null? s1)
			s2
		)
		(
			(pair? s1)
			(cons
				(car s1)
				(appendInf (cdr s1) s2)
			)
		)
		(else
			(lambda ()
				(appendInf s2 (s1))
			)
		)
	)
)
*/
func appendInf(s1, s2 StreamOfSubstitutions) StreamOfSubstitutions {
	if s1 == nil {
		return s2
	}
	car, cdr := s1()
	if car != nil {
		return func() (Substitution, StreamOfSubstitutions) {
			return car, appendInf(cdr, s2)
		}
	}
	return Suspension(appendInf(s2, cdr))
}
