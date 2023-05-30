package comicro

import (
	"context"
	"sort"
	"strings"
)

type StreamOfStates <-chan *State

// CarCdr returns both the car and cdr of a stream of states
func (stream StreamOfStates) CarCdr() (*State, StreamOfStates) {
	if stream == nil {
		return nil, nil
	}
	s, ok := <-stream
	if !ok {
		return s, nil
	}
	return s, stream
}

// String returns a string representation of a stream of states.
// Warning: If the list is infinite this function will not terminate.
func (stream StreamOfStates) String() string {
	buf := []string{}
	if stream != nil {
		for s := range stream {
			if s != nil {
				buf = append(buf, s.String())
			}
		}
	}
	sort.Strings(buf)
	return "(" + strings.Join(buf, " ") + ")"
}

// ConsStream returns a stream from a head plus a continuation
func ConsStream(ctx context.Context, s *State, proc func() StreamOfStates) StreamOfStates {
	c := make(chan *State, 0)
	go func() {
		defer close(c)
		select {
		case <-ctx.Done():
			return
		case c <- s:
		}

		if proc == nil {
			return
		}
		stream := proc()
		if stream == nil {
			return
		}
		for s := range stream {
			select {
			case <-ctx.Done():
				return
			case c <- s:
			}
		}
	}()
	return c
}

// NewSingletonStream returns the input state as a stream of states containing only the head state.
func NewSingletonStream(ctx context.Context, s *State) StreamOfStates {
	return ConsStream(ctx, s, nil)
}

// Suspension prepends a nil state infront of the input stream of states.
func Suspension(ctx context.Context, proc func() StreamOfStates) StreamOfStates {
	newStream := make(chan *State, 0)
	go func() {
		defer close(newStream)
		select {
		case <-ctx.Done():
			return
		case newStream <- nil:
		}
		if proc == nil {
			return
		}
		stream := proc()
		if stream == nil {
			return
		}
		for state := range stream {
			select {
			case <-ctx.Done():
				return
			case newStream <- state:
			}
		}
	}()
	return newStream
}

/*
Zzz is the macro to add inverse-eta-delay less tediously

(define- syntax Zzz
(syntax- rules ()

	((_ g) (λg (s/c) (λ$ () (g s/c))))))
*/
func Zzz(g Goal) Goal {
	return func(ctx context.Context, s *State) StreamOfStates {
		return Suspension(ctx, func() StreamOfStates {
			return g(ctx, s)
		})
	}
}

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
