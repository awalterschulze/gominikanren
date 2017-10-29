package micro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
scheme code:

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
func walk(v *ast.SExpr, s Substitutions) *ast.SExpr {
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

for example:

	assv v s <==> (assp (λ (v) (var=? u v)) s))
*/
func assv(v *ast.Variable, s Substitutions) (*ast.SExpr, bool) {
	if s == nil {
		return nil, false
	}
	pair := s.Car
	if pair.IsPair() {
		left := pair.Car()
		if left.IsVariable() {
			if v.Equal(left.Atom.Var) {
				return pair, true
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
scheme code:

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
func walkStar(v *ast.SExpr, s Substitutions) *ast.SExpr {
	vv := walk(v, s)
	if vv.IsVariable() {
		return vv
	}
	if vv.IsPair() {
		carv := vv.Pair.Car
		cdrv := vv.Pair.Cdr
		wcar := walkStar(carv, s)
		wcdr := walkStar(cdrv, s)
		w := ast.Cons(wcar, wcdr)
		return w
	}
	return vv
}
