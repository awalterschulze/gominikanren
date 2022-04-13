package micro

import (
	"strings"
)

// StreamOfStates is a tuple (*State, *StreamOfStates)
// the first call to stream.Cdr() stores its result in mem
type StreamOfStates struct {
	state *State
	proc  func() *StreamOfStates
	mem   *StreamOfStates
}

func (s *StreamOfStates) CarCdr() (*State, *StreamOfStates) {
	if s.proc == nil {
		return s.state, nil
	}
	if s.mem == nil {
		s.mem = s.proc()
	}
	return s.state, s.mem
}

// String returns a string representation of a stream of states.
// Warning: If the list is infinite this function will not terminate.
func (stream *StreamOfStates) String() string {
	buf := []string{}
	var s *State
	for stream != nil {
		s, stream = stream.CarCdr()
		buf = append(buf, s.String())
	}
	return "(" + strings.Join(buf, " ") + ")"
}

func NewStream(s *State, proc func() *StreamOfStates) *StreamOfStates {
	return &StreamOfStates{state: s, proc: proc}
}

// NewSingletonStream returns the input state as a stream of states containing only the head state.
func NewSingletonStream(s *State) *StreamOfStates {
	return NewStream(s, nil)
}

// Suspension prepends a nil state infront of the input stream of states.
func Suspension(s func() *StreamOfStates) *StreamOfStates {
	return &StreamOfStates{state: nil, proc: s}
}

/*
Zzz is the macro to add inverse-eta-delay less tediously

(define- syntax Zzz
(syntax- rules ()
    ((_ g) (λg (s/c) (λ$ () (g s/c))))))
*/
func Zzz(g Goal) Goal {
	return func() GoalFn {
		return func(s *State) *StreamOfStates {
			return Suspension(func() *StreamOfStates {
				return g()(s)
			})
		}
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
func takeStream(n int, s *StreamOfStates) []*State {
	if n == 0 {
		return nil
	}
	if s == nil {
		return nil
	}
	car, cdr := s.CarCdr()
	if car != nil {
		ss := takeStream(n-1, cdr)
		return append([]*State{car}, ss...)
	}
	return takeStream(n, cdr)
}
