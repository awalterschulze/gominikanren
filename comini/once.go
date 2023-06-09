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
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		gs := comicro.NewStreamForGoal(ctx, g, s)
		onceLoop(ctx, gs, ss)
	}
}

func onceLoop(ctx context.Context, ss comicro.StreamOfStates, res comicro.StreamOfStates) {
	if ss == nil {
		return
	}
	var rest comicro.StreamOfStates = nil
	state, ok := ss.Read(ctx)
	if ok {
		rest = ss
	}
	if state != nil {
		res.Write(ctx, state)
		return
	}
	onceLoop(ctx, rest, res)
}
