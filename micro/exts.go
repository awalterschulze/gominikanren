package micro

import "github.com/awalterschulze/gominikanren/sexpr/ast"

/*
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

exts either extends a substitution s with an association between the variable x and the value v ,
or it produces (ok = false)s
if extending the substitution with the pair `(,x . ,v) would create a cycle.
*/
func exts(x, v *ast.SExpr, s Substitution) (Substitution, bool) {
	if occurs(x, v, s) {
		return nil, false
	}
	pair := ast.NewAssosiation(x, v)
	pair.List.Quoted = ""
	list := ast.Prepend(s.IsQuoted(), pair, s)
	return list.List, true
}

/*
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
func occurs(x, v *ast.SExpr, s Substitution) bool {
	vv := walk(v, s)
	if vv.IsVariable() {
		return eqv(vv, x)
	}
	if vv.List != nil {
		return occurs(x, vv.List.Car(), s) || occurs(x, vv.List.Cdr(), s)
	}
	return false
}
