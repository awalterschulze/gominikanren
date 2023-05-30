package comini

import (
	"github.com/awalterschulze/gominikanren/comicro"
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
func OnceO(g comicro.Goal) comicro.Goal {
	return func(s *comicro.State) comicro.StreamOfStates {
		return onceLoop(g(s))
	}
}

func onceLoop(ss comicro.StreamOfStates) comicro.StreamOfStates {
	if ss == nil {
		return nil
	}
	car, cdr := ss.CarCdr()
	if car != nil {
		return comicro.NewSingletonStream(car)
	}
	return comicro.Suspension(func() comicro.StreamOfStates {
		return onceLoop(cdr)
	})
}
