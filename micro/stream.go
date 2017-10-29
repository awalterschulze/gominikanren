package micro

import (
	"strings"
)

// StreamOfStates is a function that returns a tuple containing the head state and the tail of the stream of states.
type StreamOfStates func() (*State, StreamOfStates)

// String returns a string representation of a stream of states.
// Warning: If the list is infinite this function will not terminate.
func (stream StreamOfStates) String() string {
	buf := []string{}
	var s *State
	for stream != nil {
		s, stream = stream()
		buf = append(buf, s.String())
	}
	return "(" + strings.Join(buf, " ") + ")"
}

// NewSingletonStream returns the input state as a stream of states containing only the head state.
func NewSingletonStream(s *State) StreamOfStates {
	return func() (*State, StreamOfStates) {
		return s, nil
	}
}

// Suspension prepends a nil state infront of the input stream of states.
func Suspension(s func() StreamOfStates) StreamOfStates {
	return func() (*State, StreamOfStates) {
		return nil, s()
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
	car, cdr := s()
	if car != nil {
		ss := takeStream(n-1, cdr)
		return append([]*State{car}, ss...)
	}
	return takeStream(n, cdr)
}
