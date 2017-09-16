package gominikanren

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

type Substitution = *ast.List

type StreamOfSubstitutions func() (Substitution, StreamOfSubstitutions)

type Goal func(Substitution) StreamOfSubstitutions

func Success() Goal {
	return func(s Substitution) StreamOfSubstitutions {
		return NewSingletonStream(s)
	}
}

func NewSingletonStream(s Substitution) StreamOfSubstitutions {
	return func() (Substitution, StreamOfSubstitutions) {
		return s, nil
	}
}

func Failure() Goal {
	return func(s Substitution) StreamOfSubstitutions {
		return nil
	}
}

func Suspension(s StreamOfSubstitutions) StreamOfSubstitutions {
	return func() (Substitution, StreamOfSubstitutions) {
		return nil, s
	}
}

type Variable struct {
}

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
			if eqv(v, left.Atom.Var) {
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

// TODO compare IDs rather than names
func eqv(v, vv *ast.Variable) bool {
	fmt.Printf("(eqv? %v %v)\n", v, vv)
	return v.Name == vv.Name
}
