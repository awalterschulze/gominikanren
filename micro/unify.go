package micro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(define (unify u v s)
	(let
		(
			(u (walk u s))
			(v (walk v s))
		)
		(cond
			(
				(eqv? u v)
				s
			)
			(
				(var? u)
				(ext-s u v s)
			)
			(
				(var? v)
				(ext-s v u s)
			)
			(
				(and
					(pair? u)
					(pair? v)
				)
				(let
					(
						(s
							(unify (car u) (car v) s)
						)
					)
					(and
						s
						(unify (cdr u) (cdr v) s)
					)
				)
			)
			(else
				#f
			)
		)
	)
)

unify returns either (ok = false) or the substitution s extended with zero or more associations,
where cycles in substitutions can lead to (ok = false)
*/
func unify(u, v *ast.SExpr, s Substitution) (Substitution, bool) {
	uu := walk(u, s)
	vv := walk(v, s)
	if uu.IsVariable() && vv.IsVariable() && uu.Atom.Var.Equal(uu.Atom.Var) {
		return s, true
	}
	if uu.IsVariable() {
		return exts(uu, vv, s)
	}
	if vv.IsVariable() {
		return exts(vv, uu, s)
	}
	if uu.IsPair() && vv.IsPair() {
		scar, sok := unify(uu.List.Car(), vv.List.Car(), s)
		if !sok {
			return nil, false
		}
		return unify(uu.List.Cdr(), vv.List.Cdr(), scar)
	}
	if uu.Equal(vv) {
		return s, true
	}
	return nil, false
}
