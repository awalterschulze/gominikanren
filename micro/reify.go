package micro

import (
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(define (reify-name n)
	(string -> symbol
		(string-append
			"_"
			(number -> string n)
		)
	)
)
*/
func reifyName(n int) string {
	return "_" + strconv.Itoa(n)
}

/*
(define (reify-s v r)
	(let
		(
			(v (walk v r))
		)
		(cond
			(
				(var? v)
				(let
					(
						(n (length r))
					)
					(let
						(
							(rn (reify-name n))
						)
						(ext-s v rn r)
					)
				)
			)
			(
				(pair? v)
				(let
					(
						(r
							(reify-s (car v) r)
						)
					)
					(reify-s (cdr v) r)
				)
			)
			(else r)
		)
	)
)
*/
func reifyS(v *ast.SExpr) Substitution {
	return reifys(v, EmptySubstitution())
}

// reifys expects a value and initially an empty substitution.
func reifys(v *ast.SExpr, r Substitution) Substitution {
	vv := walk(v, r)
	if vv.IsVariable() {
		n := len(r.Items)
		rn := reifyName(n)
		s, ok := exts(vv, ast.NewSymbol(rn), r)
		if !ok {
			panic("impossible, since the variable is fresh")
		}
		return s
	}
	if vv.IsPair() {
		car := vv.List.Head()
		rr := reifys(car, r)
		cdr := vv.List.Tail()
		return reifys(cdr, rr)
	}
	return r
}

/*
(define (reify v)
	(lambda (s)
		(let
			(
				(v (walkStar v s))
			)
			(let
				(
					(r (reify-s v empty-s))
				)
				(walkStar v r)
			)
		)
	)
)
*/
func Reify(v *ast.SExpr) func(Substitution) *ast.SExpr {
	return func(s Substitution) *ast.SExpr {
		vv := walkStar(v, s)
		r := reifyS(vv)
		return walkStar(vv, r)
	}
}
