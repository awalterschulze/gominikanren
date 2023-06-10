package comicro

import (
	"context"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// Goal is a function that takes a state and returns a stream of states.
type Goal func(context.Context, *State, StreamOfStates)

/*
RunGoal calls a goal with an emptystate and n possible resulting states.

scheme code:

	(define (run-goal n g)
		(takeInf n (g empty-s))
	)

If n == -1 then all possible states are returned.
*/
func RunGoal(ctx context.Context, n int, g Goal) []*State {
	ss := NewStreamForGoal(ctx, g, EmptyState())
	return Take(ctx, n, ss)
}

// Run behaves like the default miniKanren run command
func Run(ctx context.Context, n int, g func(*ast.SExpr) Goal) []*ast.SExpr {
	ss := RunStream(ctx, g)
	return Take(ctx, n, ss)
}

// RunStream behaves like the default miniKanren run command, but returns a stream of answers
func RunStream(ctx context.Context, g func(*ast.SExpr) Goal) chan *ast.SExpr {
	v := Var(0)
	ss := NewStreamForGoal(ctx, g(v), &State{nil, 1})
	res := make(chan *ast.SExpr, 0)
	go func() {
		defer close(res)
		s, ok := ss.Read(ctx)
		if !ok {
			return
		}
		r := MKReify(s)
		if ok := Write(ctx, r, res); !ok {
			return
		}
	}()
	return res
}

// SuccessO is a goal that always returns the input state in the resulting stream of states.
func SuccessO(ctx context.Context, s *State, ss StreamOfStates) {
	ss.Write(ctx, s)
}

// FailureO is a goal that always returns an empty stream of states.
func FailureO(ctx context.Context, s *State, ss StreamOfStates) {
	return
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
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		ss.Write(ctx, Unify(s, u, v))
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
func NeverO(ctx context.Context, s *State, ss StreamOfStates) {
	for {
		if ok := ss.Write(ctx, nil); !ok {
			return
		}
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
func AlwaysO(ctx context.Context, s *State, ss StreamOfStates) {
	for {
		if ok := ss.Write(ctx, nil); !ok {
			return
		}
		if ok := ss.Write(ctx, s); !ok {
			return
		}
	}
}
