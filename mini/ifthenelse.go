package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
)

/*
IfThenElseO is a goal that evaluates the second goal if the first goal is successful,
otherwise it evaluates the third goal.

scheme code:

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

To explain let loop:
https://stackoverflow.com/questions/31909121/how-does-the-named-let-in-the-form-of-a-loop-work

These two expressions are equivalent:

scheme code:

	(define (number->list n)
		(let loop ((n n) (acc '()))
		(if (< n 10)
			(cons n acc)
			(loop (quotient n 10)
					(cons (remainder n 10) acc)))))

scheme code:

	(define (number->list n)
		(define (loop n acc)
			(if (< n 10)
				(cons n acc)
				(loop (quotient n 10)
					(cons (remainder n 10) acc))))
		(loop n '()))

let loop not only declares a function, called loop, but also calls it, in the same line.
*/
func IfThenElseO(g1, g2, g3 micro.Goal) micro.Goal {
	return func(s *micro.State) micro.StreamOfStates {
		return ifThenElseLoop(g2, g3, s, g1(s))
	}
}

func ifThenElseLoop(g2, g3 micro.Goal, s *micro.State, sinf micro.StreamOfStates) micro.StreamOfStates {
	if sinf == nil {
		return g3(s)
	}
	car, cdr := sinf()
	if car != nil {
		return micro.Bind(sinf, g2)
	}
	return micro.Suspension(func() micro.StreamOfStates {
		return ifThenElseLoop(g2, g3, s, cdr)
	})
}
