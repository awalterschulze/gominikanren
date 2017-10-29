package micro

import (
	"strconv"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

type State struct {
	Substitution
	Counter int
}

func (s *State) String() string {
	if s.Substitution == nil {
		return "(" + "()" + " . " + strconv.Itoa(s.Counter) + ")"	
	}
	return "(" + s.Substitution.String() + " . " + strconv.Itoa(s.Counter) + ")"
}

func EmptyState() *State {
	return &State{}
}

type Substitution = *ast.Pair

type Goal func(*State) StreamOfSubstitutions

/*
(define (run-goal n g) 
	(takeInf n (g empty-s))
)
*/
func RunGoal(n int, g Goal) []*State {
	ss := g(EmptyState())
	return takeStream(n, ss)
}

func SuccessO() Goal {
	return func(s *State) StreamOfSubstitutions {
		return NewSingletonStream(s)
	}
}

func FailureO() Goal {
	return func(s *State) StreamOfSubstitutions {
		return nil
	}
}

/*
(define (= u v)
	(lambda_g (s/c) 
		(let ((s (unify u v (car s/c))))
			(if s (unit `(,s . ,(cdr s/c))) mzero)
		)
	)
)
*/
func EqualO(u, v *ast.SExpr) Goal {
	return func(s *State) StreamOfSubstitutions {
		ss, sok := unify(u, v, s.Substitution)
		if sok {
			return NewSingletonStream(&State{Substitution: ss, Counter: s.Counter})
		}
		return nil
	}
}

/*
(define (nevero) 
	(lambda (s)
		(lambda () 
			((nevero) s)
		)
	)
)
*/
func NeverO() Goal {
	return func(s *State) StreamOfSubstitutions {
		return Suspension(func() StreamOfSubstitutions {
			return NeverO()(s)
		})
	}
}

/*
(define (alwayso) 
	(lambda (s)
		(lambda ()
			(
				(disj 
					S 
					(alwayso)
				)
				s
			)
		)
	)
)
*/
func AlwaysO() Goal {
	return func(s *State) StreamOfSubstitutions {
		return Suspension(func() StreamOfSubstitutions {
			return DisjointO(
				SuccessO(),
				AlwaysO(),
			)(s)
		})
	}
}