package micro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
walk has been modified to assume it is getting a variable, but here is the original scheme code:

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
func walk(v *ast.Variable, s Substitutions) *ast.SExpr {
	value, ok := assv(v, s)
	if !ok {
		return &ast.SExpr{Atom: &ast.Atom{Var: v}}
	}
	if !value.IsVariable() {
		return value
	}
	return walk(value.Atom.Var, s)
}

/*
assv either produces the first association in s that has v as its car using eqv,
or produces ok = false if l has no such association.

for example:

	assv v s <==> (assp (λ (v) (var=? u v)) s))
*/
func assv(v *ast.Variable, ss Substitutions) (*ast.SExpr, bool) {
	if ss == nil {
		return nil, false
	}
	value, ok := ss[v.Name]
	if !ok {
		return nil, false
	}
	return value, true
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
	vv := v
	if v.IsVariable() {
		vv = walk(v.Atom.Var, s)
	}
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
