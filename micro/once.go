package micro

/*
(define (once g)
	(lambda (s)
		(let loop
			(
				(s1 (g s))
			)
			(cond
				(
					(null? s1)
					()
				)
				(
					(pair? s1)
					(cons (car s1) ())
				)
				(else
					(lambda ()
						(loop (s1))
					)
				)
			)
		)
	)
)
*/
func OnceO(g Goal) Goal {
	return func(s Substitution) StreamOfSubstitutions {
		return onceLoop(g(s))
	}
}

func onceLoop(ss StreamOfSubstitutions) StreamOfSubstitutions {
	if ss == nil {
		return nil
	}
	car, cdr := ss()
	if car != nil {
		return NewSingletonStream(car)
	}
	return Suspension(
		func() StreamOfSubstitutions {
			return onceLoop(cdr)
		},
	)
}
