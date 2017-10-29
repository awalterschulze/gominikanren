package micro

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
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
// call/fresh second arguments is a function that expects a varaible and returns a Goal.
func CallFresh(f func(v *ast.SExpr) Goal) Goal {
	return func(s *State) StreamOfSubstitutions {
		v := ast.NewVariable(fmt.Sprintf("v%d", s.Counter))
		ss := &State{s.Substitution, s.Counter + 1}
		return f(v)(ss)
	}
}
