package comicro

import (
	"context"
)

// Goal is a function that takes a state and returns a stream of states.
type Goal func(context.Context, *State, StreamOfStates)

// RunGoal calls a goal with an emptystate and n possible resulting states.
func RunGoal(ctx context.Context, n int, g Goal) []*State {
	ss := NewStreamForGoal(ctx, g, NewEmptyState())
	return Take(ctx, n, ss)
}

// Run behaves like the default miniKanren run command
func Run[A any](ctx context.Context, n int, create varCreator, g func(A) Goal) []any {
	ss := RunStream(ctx, create, g)
	return Take(ctx, n, ss)
}

// RunStream behaves like the default miniKanren run command, but returns a stream of answers
func RunStream[A any](ctx context.Context, create varCreator, g func(A) Goal) chan any {
	s := NewEmptyState()
	var v A
	s, v = NewVar(s, v)
	ss := NewStreamForGoal(ctx, g(v), s)
	res := make(chan any, 0)
	go func() {
		defer close(res)
		MapOverNonNilStream(ctx, ss, func(state *State) any { return Reify(state, create) }, res)
	}()
	return res
}

// SuccessO is a goal that always returns the input state in the resulting stream of states.
func SuccessO(ctx context.Context, s *State, ss StreamOfStates) {
	ss.Write(ctx, s)
}

// FailureO is a goal that always returns an empty stream of states.
func FailureO(ctx context.Context, s *State, ss StreamOfStates) {
	return
}

// EqualO returns a Goal that unifies the input expressions in the output stream.
func EqualO[A any](x, y A) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		ss.Write(ctx, Unify(s, x, y))
	}
}

// NeverO is a Goal that returns a never ending stream of suspensions.
func NeverO(ctx context.Context, s *State, ss StreamOfStates) {
	for {
		if ok := ss.Write(ctx, nil); !ok {
			return
		}
	}
}

// AlwaysO is a goal that returns a never ending stream of success.
func AlwaysO(ctx context.Context, s *State, ss StreamOfStates) {
	for {
		if ok := ss.Write(ctx, nil); !ok {
			return
		}
		if ok := ss.Write(ctx, s); !ok {
			return
		}
	}
}
