package gomini

import (
	"context"
	"sync"
)

// EqualO returns a Goal that unifies the input expressions in the output stream.
func EqualO[A any](x, y A) Goal {
	return func(ctx context.Context, s *State, ss Stream) {
		unifiedState := unify(x, y, s)
		if unifiedState == nil {
			return
		}
		ss.Write(ctx, unifiedState)
	}
}

// SuccessO is a goal that always returns the input state in the resulting stream of states.
func SuccessO(ctx context.Context, s *State, ss Stream) {
	ss.Write(ctx, s)
}

// FailureO is a goal that always returns an empty stream of states.
func FailureO(ctx context.Context, s *State, ss Stream) {
	return
}

// Disj is a goal that returns a logical OR of the input goals.
func disj2(g1, g2 Goal) Goal {
	return func(ctx context.Context, s *State, ss Stream) {
		wait := sync.WaitGroup{}
		Go(ctx, &wait, func() { g1(ctx, s, ss) })
		Go(ctx, &wait, func() { g2(ctx, s, ss) })
		wait.Wait()
	}
}

// Mplus is responsible for merging streams
func Mplus(ctx context.Context, s1, s2, res Stream) {
	wait := sync.WaitGroup{}
	Go(ctx, &wait, func() { writeStreamTo(ctx, s1, res) })
	Go(ctx, &wait, func() { writeStreamTo(ctx, s2, res) })
	wait.Wait()
}

// Conj is a goal that returns a logical AND of the input goals.
func conj2(g1, g2 Goal) Goal {
	return func(ctx context.Context, s *State, ss Stream) {
		g1s := NewStreamForGoal(ctx, g1, s)
		Bind(ctx, g1s, g2, ss)
	}
}

//	Original: func Bind(s *StreamOfStates, g Goal) *StreamOfStates {
//	          where type Goal = func(*State) *StreamOfStates
//
// Haskell:  (>>=)  :: m a -> (  a -> m b) -> m b

// Bind is the monad bind function for state.
func Bind(ctx context.Context, stream Stream, g Goal, res Stream) {
	wait := sync.WaitGroup{}
	for {
		state, ok := stream.Read(ctx)
		if !ok {
			break
		}
		Go(ctx, &wait, func() { g(ctx, state, res) })
	}
	wait.Wait()
}

// Exists expects a function that expects a variable and returns a Goal.
func ExistO[A any](f func(A) Goal) Goal {
	return func(ctx context.Context, s *State, ss Stream) {
		var v A
		s, v = NewVar[A](s)
		f(v)(ctx, s, ss)
	}
}

// Disjs is a macro that extends disjunction to arbitrary arguments
func DisjO(gs ...Goal) Goal {
	if len(gs) == 2 {
		return disj2(gs[0], gs[1])
	}
	return func(ctx context.Context, s *State, ss Stream) {
		wait := sync.WaitGroup{}
		for i := range gs {
			f := func(i int) func() { return func() { gs[i](ctx, s, ss) } }(i)
			Go(ctx, &wait, f)
		}
		wait.Wait()
	}
}

// Conjs is a macro that extends conjunction to arbitrary arguments
func ConjO(gs ...Goal) Goal {
	if len(gs) == 0 {
		return SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	if len(gs) == 2 {
		return conj2(gs[0], gs[1])
	}
	g1 := gs[0]
	g2 := ConjO(gs[1:]...)
	return func(ctx context.Context, s *State, ss Stream) {
		g1s := NewStreamForGoal(ctx, g1, s)
		Bind(ctx, g1s, g2, ss)
	}
}
