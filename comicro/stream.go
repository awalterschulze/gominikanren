package comicro

import (
	"context"
	"sort"
	"strings"
)

type StreamOfStates chan *State

func NewEmptyStream() StreamOfStates {
	return make(chan *State, 0)
}

func (ss StreamOfStates) Read(ctx context.Context) (*State, bool) {
	return ReadFromStream(ctx, ss)
}

func ReadFromStream[A any](ctx context.Context, c <-chan A) (A, bool) {
	var zero A
	if c == nil {
		return zero, false
	}
	select {
	case a, ok := <-c:
		return a, ok
	case <-ctx.Done():
	}
	return zero, false
}

func (ss StreamOfStates) ReadNonNil(ctx context.Context) (*State, bool) {
	return ReadNonNilFromStream(ctx, ss)
}

func ReadNonNilFromStream[A comparable](ctx context.Context, c <-chan A) (A, bool) {
	var zero A
	for {
		a, ok := ReadFromStream(ctx, c)
		if !ok {
			return zero, false
		}
		if a != zero {
			return a, true
		}
	}
}

func (ss StreamOfStates) Write(ctx context.Context, s *State) bool {
	return WriteToStream(ctx, s, ss)
}

func WriteToStream[A any](ctx context.Context, a A, c chan<- A) bool {
	if c == nil {
		return false
	}
	select {
	case <-ctx.Done():
		return false
	case c <- a:
		return true
	}
}

func (ss StreamOfStates) Close() {
	close(ss)
}

func WriteStreamTo[A any](ctx context.Context, src <-chan A, dst chan<- A) {
	MapStream(ctx, src, func(a A) A {
		return a
	}, dst)
}

func MapStream[A, B any](ctx context.Context, src <-chan A, f func(A) B, dst chan<- B) {
	for {
		a, ok := ReadFromStream(ctx, src)
		if !ok {
			return
		}
		b := f(a)
		if ok := WriteToStream(ctx, b, dst); !ok {
			return
		}
	}
}

func MapOverNonNilStream[A comparable, B any](ctx context.Context, src <-chan A, f func(A) B, dst chan<- B) {
	for {
		a, ok := ReadNonNilFromStream(ctx, src)
		if !ok {
			return
		}
		b := f(a)
		if ok := WriteToStream(ctx, b, dst); !ok {
			return
		}
	}
}

// String returns a string representation of a stream of states.
// Warning: If the list is infinite this function will not terminate.
func (ss StreamOfStates) String() string {
	buf := []string{}
	if ss != nil {
		for s := range ss {
			if s != nil {
				buf = append(buf, s.String())
			}
		}
	}
	sort.Strings(buf)
	return "(" + strings.Join(buf, " ") + ")"
}

// NewSingletonStream returns the input state as a stream of states containing only the head state.
func NewSingletonStream(ctx context.Context, s *State) StreamOfStates {
	ss := NewEmptyStream()
	Go(ctx, nil, func() {
		defer close(ss)
		ss.Write(ctx, s)
	})
	return ss
}

func NewStreamForGoal(ctx context.Context, g Goal, s *State) StreamOfStates {
	ss := NewEmptyStream()
	Go(ctx, nil, func() {
		defer close(ss)
		g(ctx, s, ss)
	})
	return ss
}

func Take[A comparable](ctx context.Context, n int, c chan A) []A {
	if n == 0 {
		return nil
	}
	if c == nil {
		return nil
	}
	var as []A
	if n > 0 {
		as = make([]A, 0, n)
	}
	i := 0
	var nilA A
	for {
		if i == n {
			break
		}
		a, ok := ReadNonNilFromStream(ctx, c)
		if !ok {
			break
		}
		if a == nilA {
			continue
		}
		as = append(as, a)
		i++
	}
	return as
}
