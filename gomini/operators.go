package gomini

import (
	"context"
	"sync"
)

// EqualO returns a Goal that unifies the input expressions in the output stream.
func EqualO[A any](x, y A) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		u := Unify(s, x, y)
		if u != nil {
			ss.Write(ctx, u)
		}
	}
}

// SuccessO is a goal that always returns the input state in the resulting stream of states.
func SuccessO(ctx context.Context, s *State, ss StreamOfStates) {
	ss.Write(ctx, s)
}

// FailureO is a goal that always returns an empty stream of states.
func FailureO(ctx context.Context, s *State, ss StreamOfStates) {
	return
}

// Disj is a goal that returns a logical OR of the input goals.
func Disj(g1, g2 Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		wait := sync.WaitGroup{}
		Go(ctx, &wait, func() { g1(ctx, s, ss) })
		Go(ctx, &wait, func() { g2(ctx, s, ss) })
		wait.Wait()
	}
}

// Mplus is responsible for merging streams
func Mplus(ctx context.Context, s1, s2, res StreamOfStates) {
	wait := sync.WaitGroup{}
	Go(ctx, &wait, func() { WriteStreamTo(ctx, s1, res) })
	Go(ctx, &wait, func() { WriteStreamTo(ctx, s2, res) })
	wait.Wait()
}

// Conj is a goal that returns a logical AND of the input goals.
func Conj(g1, g2 Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		g1s := NewStreamForGoal(ctx, g1, s)
		Bind(ctx, g1s, g2, ss)
	}
}

// Bind is the monad bind function for state.
func Bind(ctx context.Context, stream StreamOfStates, g Goal, res StreamOfStates) {
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
func Exists[A any](f func(A) Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		var typ A
		s1, v := NewVar(s, typ)
		f(v)(ctx, s1, ss)
	}
}

// Disjs is a macro that extends disjunction to arbitrary arguments
func Disjs(gs ...Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		wait := sync.WaitGroup{}
		for i := range gs {
			f := func(i int) func() { return func() { gs[i](ctx, s, ss) } }(i)
			Go(ctx, &wait, f)
		}
		wait.Wait()
	}
}

// Conjs is a macro that extends conjunction to arbitrary arguments
func Conjs(gs ...Goal) Goal {
	if len(gs) == 0 {
		return SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := Conjs(gs[1:]...)
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		g1s := NewStreamForGoal(ctx, g1, s)
		Bind(ctx, g1s, g2, ss)
	}
}
