package micro

import (
	"strings"
)

type StreamOfSubstitutions func() (*State, StreamOfSubstitutions)

func (stream StreamOfSubstitutions) String() string {
	buf := []string{}
	var s *State
	for stream != nil {
		s, stream = stream()
		buf = append(buf, s.String())
	}
	return "(" + strings.Join(buf, " ") + ")"
}

func NewSingletonStream(s *State) StreamOfSubstitutions {
	return func() (*State, StreamOfSubstitutions) {
		return s, nil
	}
}

func Suspension(s func() StreamOfSubstitutions) StreamOfSubstitutions {
	return func() (*State, StreamOfSubstitutions) {
		return nil, s()
	}
}

/*
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
*/
// n == -1 results in the whole stream being returned.
func takeStream(n int, s StreamOfSubstitutions) []*State {
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
