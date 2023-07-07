package gomini

import (
	"context"
)

// Goal is a function that takes a state and returns a stream of states.
type Goal func(context.Context, *State, Stream)

// Run behaves like the default miniKanren run command
func Run[A any](ctx context.Context, n int, s *State, g func(A) Goal) []any {
	ss := RunStream(ctx, s, g)
	return Take(ctx, n, ss)
}

// RunStream behaves like the default miniKanren run command, but returns a stream of answers
func RunStream[A any](ctx context.Context, s *State, g func(A) Goal) chan any {
	var v A
	s, v = NewVar(s, v)
	key, _ := s.CastVar(v)
	ss := NewStreamForGoal(ctx, g(v), s)
	res := make(chan any, 0)
	// rw finds rewrites for the first introduced variable.
	// This means it rewrites all substitutions for the first introduced variable.
	// For any variables without substitutions, it adds a placeholder value.
	rw := func(s *State) any {
		return rewrite(key, s)
	}
	go func() {
		defer close(res)
		MapOverStream(ctx, ss, rw, res)
	}()
	return res
}
