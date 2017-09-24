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
	fmt.Printf("unify:u=%v;v=%v;s=%v\n", u, v, s)
	uu := walk(u, s)
	vv := walk(v, s)
	fmt.Printf("unify: walked %v %v %v\n", uu, vv, s)
	if uu.IsVariable() && vv.IsVariable() && uu.Atom.Var.Equal(uu.Atom.Var) {
		fmt.Printf("unify: equal vars %v %v %v\n", uu, vv, s)
		return s, true
	}
	if uu.IsVariable() {
		fmt.Printf("unify: added left var %v %v %v\n", uu, vv, s)
		return exts(uu, vv, s)
	}
	if vv.IsVariable() {
		fmt.Printf("unify: added right var %v %v %v\n", uu, vv, s)
		return exts(vv, uu, s)
	}
	if uu.IsPair() && vv.IsPair() {
		fmt.Printf("unify pairs:u=%v;v=%v;s=%v\n", uu, vv, s)
		uucar, vvcar := uu.Car(), vv.Car()
		fmt.Printf("unify: cars %v %v %v\n", uucar, vvcar, s)
		scar, sok := unify(uucar, vvcar, s)
		if !sok {
			return nil, false
		}
		uucdr, vvcdr := uu.Cdr(), vv.Cdr()
		fmt.Printf("unify: cdr of %v %v\n", uu, uucdr)
		fmt.Printf("unify: cdr of %v %v\n", vv, vvcdr)
		return unify(uucdr, vvcdr, scar)
	}
	if uu.Equal(vv) {
		fmt.Printf("unify: equal\n")
		return s, true
	}
	fmt.Printf("unify: FALSE %v %v\n", uu, vv)
	return nil, false
}
