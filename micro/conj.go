package micro

/*
(define (conj g1 g2)
	(lambda (s)
		(append-mapInf g2 (g1 s))
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
	return func(s Substitution) StreamOfSubstitutions {
		g1s := g1(s)
		return Bind(g2, g1s)
	}
}

/*
(define (append-mapInf g s)
	(cond
		(
			(null? s)
			'()
		)
		(
			(pair? s)
			(appendInf
				(g (car s))
				(append-mapInf
					g
					(cdr s)
				)
			)
		)
		(else
			(lambda ()
				(append-mapInf g (s))
			)
		)
	)
)
*/
func Bind(g Goal, s StreamOfSubstitutions) StreamOfSubstitutions {
	if s == nil {
		return nil
	}
	car, cdr := s()
	if car != nil {
		return appendStream(
			g(car),
			Bind(
				g,
				cdr,
			),
		)
	} else {
		return Suspension(func() StreamOfSubstitutions {
			return Bind(
				g,
				cdr,
			)
		})
	}
}
