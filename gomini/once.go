package gomini

import (
	"context"
)

// OnceO is a goal that returns one successful state.
func OnceO(g Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		gs := NewStreamForGoal(ctx, g, s)
		state, ok := gs.ReadNonNil(ctx)
		if !ok {
			return
		}
		ss.Write(ctx, state)
	}
}
