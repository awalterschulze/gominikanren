package micro

import (
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
scheme code:

	(define (reify-name n)
		(string->symbol
			(string-append "_" "." (number->string n))
		)
	)
*/
func reifyName(n int) *ast.SExpr {
	return ast.NewSymbol("_" + strconv.Itoa(n))
}

// reifys expects a value and initially an empty substitution.
func reifyS(v *ast.SExpr) Substitutions {
	return reifys(v, nil)
}

/*
scheme code:

	(define (reify-s v s)
		(let
			((v (walk v s)))
			(cond
				(
					(var? v)
					(let
						((n (reify-name (length s))))
						(cons `(,v . ,n) s)
					)
				)
				(
					(pair? v)
					(reify-s (cdr v) (reify-s (car v) s))
				)
				(else s)
			)
		)
	)
*/
func reifys(v *ast.SExpr, s Substitutions) Substitutions {
	vv := walk(v, s)
	if vv.IsVariable() {
		n := reifyName(length(s))
		if s == nil {
			return ast.Cons(ast.Cons(vv, n), nil).Pair
		}
		return ast.Cons(ast.Cons(vv, n), &ast.SExpr{Pair: s}).Pair
	}
	if vv.IsPair() {
		car := vv.Car()
		cars := reifys(car, s)
		cdr := vv.Cdr()
		return reifys(cdr, cars)
	}
	return s
}

// length returns the length of the list of substitutions.
func length(s Substitutions) int {
	if s == nil {
		return 0
	}
	if s.Cdr == nil {
		return 1
	}
	if s.Cdr.Pair == nil {
		return 2
	}
	return 1 + length(s.Cdr.Pair)
}

// ReifyVarFromState is a curried function that reifies the input variable for the given input state.
func ReifyVarFromState(v *ast.SExpr) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(v, s.Substitutions)
		r := reifyS(vv)
		return walkStar(vv, r)
	}
}
