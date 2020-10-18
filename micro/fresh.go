package micro

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
CallFresh expects a function that expects a varaible and returns a Goal.

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
func CallFresh(f func(v *ast.SExpr) Goal) Goal {
    return func() GoalFn {
	return func(s *State) StreamOfStates {
		v := ast.NewVariable(fmt.Sprintf("v%d", s.Counter))
		ss := &State{s.Substitutions, s.Counter + 1}
		return f(v)()(ss)
	}
    }
}
