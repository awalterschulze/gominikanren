package comini

import (
	"context"

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
	return func(ctx context.Context, s *comicro.State) comicro.StreamOfStates {
		return onceLoop(ctx, g(ctx, s))
	}
}

func onceLoop(ctx context.Context, ss comicro.StreamOfStates) comicro.StreamOfStates {
	if ss == nil {
		return nil
	}
	car, cdr := ss.CarCdr()
	if car != nil {
		return comicro.NewSingletonStream(ctx, car)
	}
	return comicro.Suspension(ctx, func() comicro.StreamOfStates {
		return onceLoop(ctx, cdr)
	})
}
