package gomini

import (
	"context"
	"sort"
	"strings"
)

type Stream chan *State

func NewEmptyStream() Stream {
	return make(chan *State, 0)
}

func (ss Stream) Read(ctx context.Context) (*State, bool) {
	return readFromStream(ctx, ss)
}

func readFromStream[A any](ctx context.Context, ss <-chan A) (A, bool) {
	var zero A
	select {
	case a, ok := <-ss:
		return a, ok
	case <-ctx.Done():
	}
	return zero, false
}

func (ss Stream) Write(ctx context.Context, s *State) bool {
	return writeToStream(ctx, s, ss)
}

func writeToStream[A any](ctx context.Context, a A, ss chan<- A) bool {
	select {
	case <-ctx.Done():
		return false
	case ss <- a:
		return true
	}
}

func (ss Stream) Close() {
	close(ss)
}

func writeStreamTo[A any](ctx context.Context, src <-chan A, dst chan<- A) {
	identity := func(a A) A {
		return a
	}
	mapOverStream(ctx, src, identity, dst)
}

func mapOverStream[A, B any](ctx context.Context, src <-chan A, f func(A) B, dst chan<- B) {
	for {
		a, ok := readFromStream(ctx, src)
		if !ok {
			return
		}
		b := f(a)
		if ok := writeToStream(ctx, b, dst); !ok {
			return
		}
	}
}

// String returns a string representation of a stream of states.
// Warning: If the list is infinite this function will not terminate.
func (ss Stream) String() string {
	buf := []string{}
	if ss != nil {
		for {
			s, ok := ss.Read(context.Background())
			if !ok {
				break
			}
			buf = append(buf, s.String())
		}
	}
	sort.Strings(buf)
	return "(" + strings.Join(buf, " ") + ")"
}

func NewStreamForGoal(ctx context.Context, g Goal, s *State) Stream {
	ss := NewEmptyStream()
	Go(ctx, nil, func() {
		defer close(ss)
		g(ctx, s, ss)
	})
	return ss
}

func take[A any](ctx context.Context, n int, ss chan A) []A {
	if n == 0 {
		return nil
	}
	if ss == nil {
		return nil
	}
	var as []A
	if n > 0 {
		as = make([]A, 0, n)
	}
	i := 0
	for {
		if i == n {
			break
		}
		a, ok := readFromStream(ctx, ss)
		if !ok {
			break
		}
		as = append(as, a)
		i++
	}
	return as
}
