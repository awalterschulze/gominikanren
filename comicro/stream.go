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
	if ss == nil {
		return nil, false
	}
	select {
	case s, ok := <-ss:
		return s, ok
	case <-ctx.Done():
	}
	return nil, false
}

func (ss StreamOfStates) Write(ctx context.Context, s *State) bool {
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

func WriteStreamTo(ctx context.Context, src StreamOfStates, dst StreamOfStates) {
	for {
		s, ok := src.Read(ctx)
		if !ok {
			return
		}
		if ok := dst.Write(ctx, s); !ok {
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
func takeStream(n int, s StreamOfStates) []*State {
	if n == 0 {
		return nil
	}
	if s == nil {
		return nil
	}
	var res []*State
	if n < 0 {
		res = []*State{}
	} else {
		res = make([]*State, 0, n)
	}
	i := 0
	for {
		if i == n {
			break
		}
		a, ok := <-s
		if !ok {
			break
		}
		if a == nil {
			continue
		}
		res = append(res, a)
		i++
	}
	return res
}
