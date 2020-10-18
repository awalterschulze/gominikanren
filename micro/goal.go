package micro

import (
    "sort"
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
type Substitutions map[string]*ast.SExpr

func (s Substitutions) String() string {
    keys := make([]string, len(s))
    i := 0
    for k, _ := range s {
        keys[i] = k
        i++
    }
    sort.Strings(keys)
    sexprs := make([]*ast.SExpr, len(s))
    for i, k := range keys {
        v := s[k]
        sexprs[len(s)-1-i] = ast.Cons(&ast.SExpr{Atom: &ast.Atom{Var: &ast.Variable{Name: k}}}, v)
    }
    l := ast.NewList(sexprs...).String()
	return l[1 : len(l)-1]
}

/*
func (s Substitution) String() string {
	return ast.Cons(&ast.SExpr{Atom: &ast.Atom{Var: &ast.Variable{Name: s.Var}}}, s.Value).String()
}
*/

// GoalFn is a function that takes a state and returns a stream of states.
type GoalFn func(*State) StreamOfStates
// Goal is a function that returns a GoalFn. Used for lazy evaluation
type Goal func() GoalFn

/*
RunGoal calls a goal with an emptystate and n possible resulting states.

scheme code:

	(define (run-goal n g)
		(takeInf n (g empty-s))
	)

If n == -1 then all possible states are returned.
*/
func RunGoal(n int, g Goal) []*State {
	ss := g()(EmptyState())
	return takeStream(n, ss)
}

// SuccessO is a goal that always returns the input state in the resulting stream of states.
func SuccessO() GoalFn {
	return func(s *State) StreamOfStates {
		return NewSingletonStream(s)
	}
}

// FailureO is a goal that always returns an empty stream of states.
func FailureO() GoalFn {
	return func(s *State) StreamOfStates {
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
    return func() GoalFn {
	return func(s *State) StreamOfStates {
		ss, sok := unify(u, v, s.Substitutions)
		if sok {
			return NewSingletonStream(&State{Substitutions: ss, Counter: s.Counter})
		}
		return nil
	}
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
    return func() GoalFn {
        return neverO()
    }
}

func neverO() GoalFn {
    return DefRel(neverO)
}

func DefRel(f func() GoalFn) GoalFn {
    return func(s *State) StreamOfStates {
        return Suspension(func() StreamOfStates {
            return f()(s)
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
    return func() GoalFn {
        return alwaysO()
    }
}

func alwaysO() GoalFn {
    return DefRel(DisjointO(SuccessO, alwaysO))
}
