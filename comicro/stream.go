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
	return Read(ctx, ss)
}

func Read[A any](ctx context.Context, ss <-chan A) (A, bool) {
	var nilA A
	if ss == nil {
		return nilA, false
	}
	select {
	case s, ok := <-ss:
		return s, ok
	case <-ctx.Done():
	}
	return nilA, false
}

func (ss StreamOfStates) ReadNonNull(ctx context.Context) (*State, bool) {
	return ReadNonNull(ctx, ss)
}

func ReadNonNull[A comparable](ctx context.Context, ss chan A) (A, bool) {
	var nilA A
	for {
		s, ok := Read(ctx, ss)
		if !ok {
			return nilA, false
		}
		if s != nilA {
			return s, true
		}
	}
}

func (ss StreamOfStates) Write(ctx context.Context, s *State) bool {
	return Write(ctx, s, ss)
}

func Write[A any](ctx context.Context, s A, ss chan<- A) bool {
	if ss == nil {
		return false
	}
	select {
	case <-ctx.Done():
		return false
	case ss <- s:
		return true
	}
}

func (ss StreamOfStates) Close() {
	close(ss)
}

func WriteStreamTo[A any](ctx context.Context, src <-chan A, dst chan<- A) {
	for {
		s, ok := Read(ctx, src)
		if !ok {
			return
		}
		if ok := Write(ctx, s, dst); !ok {
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
	go func() {
		defer close(ss)
		ss.Write(ctx, s)
	}()
	return ss
}

func NewStreamForGoal(ctx context.Context, g Goal, s *State) StreamOfStates {
	ss := NewEmptyStream()
	go func() {
		defer close(ss)
		g(ctx, s, ss)
	}()
	return ss
}

/*
Zzz is the macro to add inverse-eta-delay less tediously

/*
takeStream returns the first n states from the stream of states as a list.

scheme code:

	(define (takeInf n s)
		(cond
			(
				(and
					n
					(zero? n)
				)
				()
			)
			(
				(null? s)
				()
			)
			(
				(pair? s)
				(cons
					(car s)
					(takeInf
						(and n (sub1 n))
						(cdr s)
					)
				)
			)
			(else
				(takeInf n (s))
			)
		)
	)

If n == -1 results in the whole stream being returned.
*/
func takeStream[A comparable](n int, s chan A) []A {
	if n == 0 {
		return nil
	}
	if s == nil {
		return nil
	}
	var res []A
	if n < 0 {
		res = make([]A, 0)
	} else {
		res = make([]A, 0, n)
	}
	i := 0
	var nilA A
	for {
		if i == n {
			break
		}
		a, ok := <-s
		if !ok {
			break
		}
		if a == nilA {
			continue
		}
		res = append(res, a)
		i++
	}
	return res
}
