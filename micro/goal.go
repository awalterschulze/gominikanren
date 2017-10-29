package micro

import (
	"strconv"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// State is a product of a list of substitutions and a variable counter.
type State struct {
	Substitutions
	Counter int
}

// String returns a string representation of State.
func (s *State) String() string {
	if s.Substitutions == nil {
		return "(" + "()" + " . " + strconv.Itoa(s.Counter) + ")"	
	}
	return "(" + s.Substitutions.String() + " . " + strconv.Itoa(s.Counter) + ")"
}

// EmptyState returns an empty state.
func EmptyState() *State {
	return &State{}
}

// Substitutions is a list of substitutions represented by a sexprs pair.
type Substitutions = *ast.Pair

// Goal is a function that takes a state and returns a stream of states.
type Goal func(*State) StreamOfStates

/*
RunGoal calls a goal with an emptystate and n possible resulting states.

scheme code:

	(define (run-goal n g) 
		(takeInf n (g empty-s))
	)

If n == -1 then all possible states are returned.
*/
func RunGoal(n int, g Goal) []*State {
	ss := g(EmptyState())
	return takeStream(n, ss)
}

// SuccessO is a goal that always returns the input state in the resulting stream of states.
func SuccessO() Goal {
	return func(s *State) StreamOfSubstitutions {
		return NewSingletonStream(s)
	}
}

// FailureO is a goal that always returns an empty stream of states.
func FailureO() Goal {
	return func(s *State) StreamOfSubstitutions {
		return nil
	}
}

/*
EqualO returns a Goal that unifies the input expressions in the output stream.

scheme code:

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
NeverO is a Goal that returns a never ending stream of suspensions.

scheme code:

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
AlwaysO is a goal that returns a never ending stream of success.

scheme code:

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