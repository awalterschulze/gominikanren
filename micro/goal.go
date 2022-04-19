package micro

import (
	"fmt"
	"sort"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// State is a product of a list of substitutions and a variable counter.
type State struct {
	Substitutions
	Counter uint64
}

// String returns a string representation of State.
func (s *State) String() string {
	if s.Substitutions == nil {
		return fmt.Sprintf("(() . %d)", s.Counter)
	}
	return fmt.Sprintf("(%s . %d)", s.Substitutions.String(), s.Counter)
}

// EmptyState returns an empty state.
func EmptyState() *State {
	return &State{}
}

type SubPair struct {
	Key   uint64
	Value *ast.SExpr
}

// Substitutions is a list of substitutions represented by a sexprs pair.
type Substitutions []SubPair

func (s Substitutions) String() string {
	sexprs := make([]*ast.SExpr, len(s))
	sort.Slice(s, func(i, j int) bool { return s[i].Key < s[j].Key })
	for i, pair := range s {
		k, v := pair.Key, pair.Value
		vv := Var(k)
		vvv := ast.Cons(vv, v)
		sexprs[len(s)-1-i] = vvv
	}
	l := ast.NewList(sexprs...).String()
	return l[1 : len(l)-1]
}

// Goal is a function that takes a state and returns a stream of states.
type Goal func(*State) *StreamOfStates

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

// Run behaves like the default miniKanren run command
func Run(n int, g func(*ast.SExpr) Goal) []*ast.SExpr {
	v := Var(0)
	ss := g(v)(&State{nil, 1})
	states := takeStream(n, ss)
	return MKReify(states)
}

// SuccessO is a goal that always returns the input state in the resulting stream of states.
func SuccessO(s *State) *StreamOfStates {
	return NewSingletonStream(s)
}

// FailureO is a goal that always returns an empty stream of states.
func FailureO(s *State) *StreamOfStates {
	return nil
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
	return func(s *State) *StreamOfStates {
		ss, sok := unify(u, v, s.Substitutions)
		if sok {
			return NewSingletonStream(&State{Substitutions: ss, Counter: s.Counter})
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
func NeverO(s *State) *StreamOfStates {
	return Suspension(func() *StreamOfStates {
		return NeverO(s)
	})
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
func AlwaysO(s *State) *StreamOfStates {
	return Suspension(func() *StreamOfStates {
		return Disj(SuccessO, AlwaysO)(s)
	})
}
