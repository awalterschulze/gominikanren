package micro

/*
(define (ifte g1 g2 g3)
	(lambda (s)
		(let loop
			(
				(s1 (g1 s))
			)
			(cond
				(
					(null? s1)
					(g3 s)
				)
				(
					(pair? s1)
					(append-map1 g2 s1)
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
func IfThenElseO(g1, g2, g3 Goal) Goal {
	return func(s Substitution) StreamOfSubstitutions {
		return ifThenElseLoop(g2, g3, s, g1(s))
	}
}

func ifThenElseLoop(g2, g3 Goal, s Substitution, sinf StreamOfSubstitutions) StreamOfSubstitutions {
	if sinf == nil {
		return g3(s)
	}
	car, cdr := sinf()
	if car != nil {
		return appendMapInf(g2, sinf)
	}
	return Suspension(func() StreamOfSubstitutions {
		return ifThenElseLoop(g2, g3, s, cdr)
	})
}

/*
To explain let loop
https://stackoverflow.com/questions/31909121/how-does-the-named-let-in-the-form-of-a-loop-work

These two expressions are equivalent:

(define (number->list n)
	(let loop ((n n) (acc '()))
	(if (< n 10)
		(cons n acc)
		(loop (quotient n 10)
				(cons (remainder n 10) acc)))))

(define (number->list n)
	(define (loop n acc)
		(if (< n 10)
			(cons n acc)
			(loop (quotient n 10)
				(cons (remainder n 10) acc))))
	(loop n '()))

let loop not only declares a function, called loop, but also calls it, in the same line.
*/
