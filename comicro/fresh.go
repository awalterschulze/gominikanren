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
		ss := s.AddCounter()
		return f(v)(ctx, ss)
	}
}

// calls n new fresh vars, it is up to the user to use them
// (ie by assigning them to go vars in f)
// not ideal, but better than 5 nested callfresh calls...
func Fresh(n int, f func(...*ast.SExpr) Goal) Goal {
    return func(ctx context.Context, s *State) StreamOfStates {
        ss := s
        vars := make([]*ast.SExpr, n)
        for i:=0; i<n; i++ {
            v := Var(ss.Counter)
            vars[i] = v
            ss = ss.AddCounter()
        }
        return f(vars...)(ctx, ss)
    }
}
