package comicro

import (
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
	var s *State
	for stream != nil {
		s, stream = stream.CarCdr()
		if s == nil {
			continue
		}
		buf = append(buf, s.String())
	}
	return "(" + strings.Join(buf, " ") + ")"
}

// ConsStream returns a stream from a head plus a continuation
func ConsStream(s *State, proc func() StreamOfStates) StreamOfStates {
	c := make(chan *State, 0)
	go func() {
		defer close(c)
		c <- s
		if proc == nil {
			return
		}
		stream := proc()
		for s := range stream {
			c <- s
		}
	}()
	return c
}

// NewSingletonStream returns the input state as a stream of states containing only the head state.
func NewSingletonStream(s *State) StreamOfStates {
	return ConsStream(s, nil)
}

// Suspension prepends a nil state infront of the input stream of states.
func Suspension(proc func() StreamOfStates) StreamOfStates {
	c := make(chan *State, 0)
	go func() {
		defer close(c)
		// we need to send two nils, otherwise when we cons execution doesn't really get suspended, but is executed immediately after the first read.
		c <- nil
		if proc == nil {
			return
		}
		stream := proc()
		if stream == nil {
			return
		}
		for s := range stream {
			c <- s
		}
	}()
	return c
}

/*
Zzz is the macro to add inverse-eta-delay less tediously

(define- syntax Zzz
(syntax- rules ()

	((_ g) (λg (s/c) (λ$ () (g s/c))))))
*/
func Zzz(g Goal) Goal {
	return func(s *State) StreamOfStates {
		return Suspension(func() StreamOfStates {
			return g(s)
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
	res := make([]*State, 0, n)
	i := 0
	for {
		if i == n {
			break
		}
		a, ok := <-s
		if !ok {
			break
		}
		if a != nil {
			res = append(res, a)
			i++
		}
	}
	return res
}
