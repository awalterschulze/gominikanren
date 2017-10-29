package micro

import (
	"fmt"
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
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
func reifyS(v *ast.SExpr) Substitution {
	return reifys(v, nil)
}

/*
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
func reifys(v *ast.SExpr, s Substitution) Substitution {
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

func length(s Substitution) int {
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

/*
(define (reify-state/1st-var s/c)
	(let
		((v (walkStar (var 0) (car s/c))))
		(walkStar v (reify-s v '()))
	)
)
*/
func reifyState1stVar(s *State) *ast.SExpr {
	v := ast.NewVariable(fmt.Sprintf("v%d", 0))
	return ReifyVarFromState(v)(s)
}

/*
(define (mK-reify s/c*)
	(map reify-state/1st-var s/c*)
)
*/
func mKreify(ss []*State) []*ast.SExpr {
	return deriveFmapReify(reifyState1stVar, ss)
}

func ReifyVarFromState(v *ast.SExpr) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(v, s.Substitution)
		r := reifyS(vv)
		return walkStar(vv, r)
	}
}
