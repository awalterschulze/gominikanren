package comicro

import (
	"context"
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// Var creates a new variable as the string vC
func Var(c uint64) *ast.SExpr {
	return ast.NewVar(fmt.Sprintf("v%d", c), c)
}

/*
CallFresh expects a function that expects a variable and returns a Goal.

scheme code:

	(define (call/fresh f)
		(lambda_g (s/c)
			(let (
				(c (cdr s/c))
			)
			(
				(f (var c))
				`(,(car s/c) . ,(+ c 1)))
			)
		)
	)
*/
func CallFresh(f func(*ast.SExpr) Goal) Goal {
	return func(ctx context.Context, s *State) StreamOfStates {
		v := Var(s.Counter)
		ss := &State{s.Substitutions, s.Counter + 1}
		return f(v)(ctx, ss)
	}
}
