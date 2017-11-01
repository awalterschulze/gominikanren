package micro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
exts either extends a substitution s with an association between the variable x and the value v ,
or it produces (ok = false)s
if extending the substitution with the pair `(,x . ,v) would create a cycle.

scheme code:

	(define (ext-s x v s)
		(cond
			(
				(occurs? x v s)
				#f
			)
			(else
				(cons
					`(,x . ,v)
					s
				)
			)
		)
	)
*/
func exts(x *ast.Variable, v *ast.SExpr, s Substitutions) (Substitutions, bool) {
	if occurs(x, v, s) {
		return nil, false
	}
	return append([]*Substitution{&Substitution{Var: x.Name, Value: v}}, s...), true
}

/*
scheme code:

	(define (occurs? x v s)
		(let
			(
				(v
					(walk v s)
				)
			)
			(cond
				(
					(var? v)
					(eqv? v x)
				)
				(
					(pair? v)
					(or
						(occurs? x (car v) s)
						(occurs? x (cdr v) s)
					)
				)
				(else
					#f
				)
			)
		)
	)
*/
func occurs(x *ast.Variable, v *ast.SExpr, s Substitutions) bool {
	vv := v
	if v.IsVariable() {
		vv = walk(v.Atom.Var, s)
	}
	if vv.IsVariable() {
		return vv.Atom.Var.Equal(x)
	}
	if vv.IsPair() {
		return occurs(x, vv.Car(), s) || occurs(x, vv.Cdr(), s)
	}
	return false
}
