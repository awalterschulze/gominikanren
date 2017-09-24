package micro

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(define (walk v s)
	(let
		(
			(a
				(and
					(var? v)
					(assv v s)
				)
			)
		)
		(cond
			(
				(pair? a)
				(walk (cdr a) s)
			)
			(else v)
		)
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
		return walk(a.Cdr(), s)
	}
	return v
}

/*
assv either produces the first association in l that has v as its car using eqv,
or produces ok = false if l has no such association.
*/
func assv(v *ast.Variable, s Substitution) (*ast.SExpr, bool) {
	// fmt.Printf("(assv %v %v)\n", v, s)
	if s.IsNil() {
		return nil, false
	}
	pair := s.Head()
	// fmt.Printf("head = %v\n", pair)
	if pair.IsAssosiation() {
		left := pair.Car()
		// fmt.Printf("left %v\n", left)
		if left.IsVariable() {
			// fmt.Printf("isvar %v\n", left)
			if v.Equal(left.Atom.Var) {
				fmt.Printf("got pair %v leftid:%v rightid:%v\n", pair, v.ID, left.Atom.Var.ID)
				return pair, true
			} else {
				// fmt.Printf("not equal to %v != %v\n", left, v)
			}
		}
	}
	tail := s.Tail()
	if tail.List == nil {
		return nil, false
	}
	return assv(v, tail.List)
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
		car := vv.List.Head()
		cdr := vv.List.Tail()
		// fmt.Printf("input: %v, car: %v, cdr: %v\n", vv, car, cdr)
		wcar := walkStar(car, s)
		wcdr := walkStar(cdr, s)
		var w *ast.SExpr
		if wcdr.List != nil {
			w = ast.Prepend(vv.List.IsQuoted(), wcar, wcdr.List)
		} else {
			w = ast.NewList(vv.List.IsQuoted(), wcar, wcdr)
		}
		// fmt.Printf("output: %v, car: %v, cdr: %v\n", w, wcar, wcdr)
		return w
	}
	return vv
}
