package micro

import (
	"strings"
)

type StreamOfSubstitutions func() (Substitution, StreamOfSubstitutions)

func (stream StreamOfSubstitutions) String() string {
	buf := []string{}
	var s Substitution
	for stream != nil {
		s, stream = stream()
		buf = append(buf, s.String())
	}
	return "(" + strings.Join(buf, " ") + ")"
}

func NewSingletonStream(s Substitution) StreamOfSubstitutions {
	return func() (Substitution, StreamOfSubstitutions) {
		return s, nil
	}
}

func Suspension(s func() StreamOfSubstitutions) StreamOfSubstitutions {
	return func() (Substitution, StreamOfSubstitutions) {
		return nil, s()
	}
}

/*
(define (appendInf s1 s2)
	(cond
		(
			(null? s1)
			s2
		)
		(
			(pair? s1)
			(cons
				(car s1)
				(appendInf (cdr s1) s2)
			)
		)
		(else
			(lambda ()
				(appendInf s2 (s1))
			)
		)
	)
)
*/
func appendInf(s1, s2 StreamOfSubstitutions) StreamOfSubstitutions {
	if s1 == nil {
		return s2
	}
	car, cdr := s1()
	if car != nil {
		return func() (Substitution, StreamOfSubstitutions) {
			return car, appendInf(cdr, s2)
		}
	}
	return Suspension(func() StreamOfSubstitutions {
		return appendInf(s2, cdr)
	})
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
func takeInf(n int, s StreamOfSubstitutions) []Substitution {
	if n == 0 {
		return nil
	}
	if s == nil {
		return nil
	}
	car, cdr := s()
	if car != nil {
		ss := takeInf(n-1, cdr)
		return append([]Substitution{car}, ss...)
	}
	return takeInf(n, cdr)
}
