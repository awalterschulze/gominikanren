package comini

import (
	"context"

	"github.com/awalterschulze/gominikanren/comicro"
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
func IfThenElseO(cond, thn, els comicro.Goal) comicro.Goal {
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		conds := comicro.NewStreamForGoal(ctx, cond, s)
		ifThenElseO(ctx, conds, thn, els, s, ss)
	}
}

func ifThenElseO(ctx context.Context, conds comicro.StreamOfStates, thn, els comicro.Goal, s *comicro.State, res comicro.StreamOfStates) {
	headState, ok := conds.ReadNonNull(ctx)
	if !ok {
		els(ctx, s, res)
		return
	}
	heads := comicro.NewStreamForGoal(ctx, thn, headState)
	rests := comicro.NewEmptyStream()
	comicro.Go(ctx, nil, func() {
		defer close(rests)
		comicro.Bind(ctx, conds, thn, rests)
	})
	comicro.Mplus(ctx, heads, rests, res)
}
