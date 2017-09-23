package micro

import (
	"fmt"

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
	fmt.Printf("unify %v %v %v\n", u, v, s)
	uu := walk(u, s)
	vv := walk(v, s)
	fmt.Printf("unify walked %v %v %v\n", uu, vv, s)
	if uu.IsVariable() && vv.IsVariable() && uu.Atom.Var.Equal(uu.Atom.Var) {
		fmt.Printf("unify vars %v %v %v\n", uu, vv, s)
		return s, true
	}
	if uu.IsVariable() {
		fmt.Printf("unify exts %v %v %v\n", uu, vv, s)
		return exts(uu, vv, s)
	}
	if vv.IsVariable() {
		fmt.Printf("unify exts %v %v %v\n", vv, uu, s)
		return exts(vv, uu, s)
	}
	if uu.IsPair() && vv.IsPair() {
		fmt.Printf("unify list\n")
		scar, sok := unify(uu.List.Car(), vv.List.Car(), s)
		if !sok {
			fmt.Printf("unify () not car %v %v %v\n", uu, vv, s)
			return nil, false
		}
		return unify(uu.List.Cdr(), vv.List.Cdr(), scar)
	}
	if uu.Equal(vv) {
		return s, true
	}
	fmt.Printf("unify () not a list\n")
	return nil, false
}
