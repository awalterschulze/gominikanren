package comicro

import (
	"context"
)

// Exists expects a function that expects a variable and returns a Goal.
func Exists[A any](f func(*A) Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		var typ *A
		s1, v := NewVar(s, typ)
		f(v)(ctx, s1, ss)
	}
}
