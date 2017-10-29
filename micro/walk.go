package micro

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(define (walk u s)
	(let (
		(pr (and
				(var? u)
				(assp (λ (v) (var=? u v)) s)
		)))
		(if pr (walk (cdr pr) s) u)
	)
)
*/
func walk(v *ast.SExpr, s Substitution) *ast.SExpr {
	// fmt.Printf("(walk %v %v)\n", v, s)
	if !v.IsVariable() {
		return v
	}
	a, ok := assv(v.Atom.Var, s)
	if !ok {
		return v
	}
	if a.IsPair() {
		return walk(a.Pair.Cdr, s)
	}
	return v
}

/*
assv either produces the first association in s that has v as its car using eqv,
or produces ok = false if l has no such association.

for example: assv v s <==> (assp (λ (v) (var=? u v)) s))
*/
func assv(v *ast.Variable, s Substitution) (*ast.SExpr, bool) {
	// fmt.Printf("(assv %v %v)\n", v, s)
	if s == nil {
		return nil, false
	}
	pair := s.Car
	// fmt.Printf("head = %v\n", pair)
	if pair.IsPair() {
		left := pair.Car()
		// fmt.Printf("left %v\n", left)
		if left.IsVariable() {
			// fmt.Printf("isvar %v\n", left)
			if v.Equal(left.Atom.Var) {
				fmt.Printf("got pair %v\n", pair)
				return pair, true
			} else {
				// fmt.Printf("not equal to %v != %v\n", left, v)
			}
		}
	}
	tail := s.Cdr
	if tail.IsPair() {
		return assv(v, tail.Pair)
	}
	return nil, false
}

/*
(define (walkStar v s)
	(let
		(
			(v (walk v s))
		)
		(cond
			(
				(var? v)
				v
			)
			(
				(pair? v)
				(cons
					(walkStar (car v) s)
					(walkStar (cdr v) s)
				)
			)
			(else v)
		)
	)
)
*/
func walkStar(v *ast.SExpr, s Substitution) *ast.SExpr {
	vv := walk(v, s)
	if vv.IsVariable() {
		return vv
	}
	if vv.IsPair() {
		carv := vv.Pair.Car
		cdrv := vv.Pair.Cdr
		// fmt.Printf("input: %v, car: %v, cdr: %v\n", vv, car, cdr)
		wcar := walkStar(carv, s)
		wcdr := walkStar(cdrv, s)
		w := ast.Cons(wcar, wcdr)
		// fmt.Printf("output: %v, car: %v, cdr: %v\n", w, wcar, wcdr)
		return w
	}
	return vv
}
