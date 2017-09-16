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
	fmt.Printf("(walk %v %v)\n", v, s)
	if !v.IsVariable() {
		return v
	}
	a, ok := assv(v.Atom.Var, s)
	if !ok {
		return v
	}
	if a.List == nil {
		return v
	}
	return walk(a.List.Cdr(), s)
}

/*
assv either produces the first association in l that has v as its car using eqv,
or produces ok = false if l has no such association.
*/
func assv(v *ast.Variable, s Substitution) (*ast.SExpr, bool) {
	fmt.Printf("(assv %v %v)\n", v, s)
	if s.IsNil() {
		return nil, false
	}
	pair := s.Car()
	if pair.IsAssosiation() {
		left := pair.List.Car()
		fmt.Printf("left %v\n", left)
		if left.IsVariable() {
			if v.Equal(left.Atom.Var) {
				fmt.Printf("got pair %v\n", pair)
				return pair, true
			}
		}
	}
	tail := s.Cdr()
	if tail.List == nil {
		return nil, false
	}
	return assv(v, tail.List)
}
