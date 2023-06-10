package comini

import (
	"context"

	"github.com/awalterschulze/gominikanren/comicro"
)

/*
ConjPlus is a macro that extends conjunction to arbitrary arguments

(define- syntax conj+
(syntax- rules ()

	((_ g) (Zzz g))
	((_ g0 g . . . ) (conj (Zzz g0) (conj+ g . . . )))))
*/
func Conjs(gs ...comicro.Goal) comicro.Goal {
	if len(gs) == 0 {
		return comicro.SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := Conjs(gs[1:]...)
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		g1s := comicro.NewStreamForGoal(ctx, g1, s)
		comicro.Bind(ctx, g1s, g2, ss)
	}
}
