package comicro

import (
	"context"
)

// CallFresh expects a function that expects a variable and returns a Goal.
func CallFresh(f func(Var) Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		s1, v := s.NewVar()
		f(v)(ctx, s1, ss)
	}
}
