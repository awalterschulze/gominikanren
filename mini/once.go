package mini

import (
	"github.com/awalterschulze/gominikanren/micro"
)

/*
OnceO is a goal that returns one successful state.

scheme code:

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
func OnceO(g micro.Goal) micro.Goal {
	return func(s *micro.State) micro.StreamOfStates {
		return onceLoop(g(s))
	}
}

func onceLoop(ss micro.StreamOfStates) micro.StreamOfStates {
	if ss == nil {
		return nil
	}
	car, cdr := ss()
	if car != nil {
		return micro.NewSingletonStream(car)
	}
	return micro.Suspension(
		func() micro.StreamOfStates {
			return onceLoop(cdr)
		},
	)
}
