package comini

import (
	"context"

	"github.com/awalterschulze/gominikanren/comicro"
)

// OnceO is a goal that returns one successful state.
func OnceO(g comicro.Goal) comicro.Goal {
	return func(ctx context.Context, s *comicro.State, ss comicro.StreamOfStates) {
		gs := comicro.NewStreamForGoal(ctx, g, s)
		state, ok := gs.ReadNonNil(ctx)
		if !ok {
			return
		}
		ss.Write(ctx, state)
	}
}
