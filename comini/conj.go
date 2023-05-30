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
func ConjPlus(gs ...comicro.Goal) comicro.Goal {
	if len(gs) == 0 {
		return comicro.SuccessO
	}
	if len(gs) == 1 {
		return comicro.Zzz(gs[0])
	}
	g1 := comicro.Zzz(gs[0])
	g2 := ConjPlus(gs[1:]...)
	return func(ctx context.Context, s *comicro.State) comicro.StreamOfStates {
		g1s := g1(ctx, s)
		return comicro.Bind(ctx, g1s, g2)
	}
}

// ConjPlusNoZzz is the conj+ macro without wrapping
// its goals in a suspension (zzz)
func ConjPlusNoZzz(gs ...comicro.Goal) comicro.Goal {
	if len(gs) == 0 {
		return comicro.SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := ConjPlusNoZzz(gs[1:]...)
	return func(ctx context.Context, s *comicro.State) comicro.StreamOfStates {
		g1s := g1(ctx, s)
		return comicro.Bind(ctx, g1s, g2)
	}
}
