package comicro

import (
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
scheme code:

	(define (reify-name n)
		(string->symbol
			(string-append "_" "." (number->string n))
		)
	)
*/
func reifyName(n int) *ast.SExpr {
	return ast.NewSymbol("_" + strconv.Itoa(n))
}

// reifys expects a value and initially an empty substitution.
func reifyS(v *ast.SExpr) Substitutions {
	return reifys(v, nil)
}

/*
scheme code:

	(define (reify-s v s)
		(let
			((v (walk v s)))
			(cond
				(
					(var? v)
					(let
						((n (reify-name (length s))))
						(cons `(,v . ,n) s)
					)
				)
				(
					(pair? v)
					(reify-s (cdr v) (reify-s (car v) s))
				)
				(else s)
			)
		)
	)
*/
func reifys(v *ast.SExpr, s Substitutions) Substitutions {
	vv := v
	if v.IsVariable() {
		vv = walk(v.Atom.Var, s)
	}
	if vv.IsVariable() {
		n := reifyName(len(s))
		return s.AddPair(vv.Atom.Var.Index, n)
	}
	if vv.IsPair() {
		car := vv.Car()
		cars := reifys(car, s)
		cdr := vv.Cdr()
		return reifys(cdr, cars)
	}
	return s
}

// ReifyVarFromState is a curried function that reifies the input variable for the given input state.
func ReifyVarFromState(v string) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(&ast.SExpr{Atom: &ast.Atom{Var: &ast.Variable{Name: v}}}, s.Substitutions)
		r := reifyS(vv)
		return walkStar(vv, r)
	}
}

// ReifyIntVarFromState is a curried function that reifies the input variable for the given input state.
func ReifyIntVarFromState(v uint64) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(&ast.SExpr{Atom: &ast.Atom{Var: &ast.Variable{Index: v}}}, s.Substitutions)
		r := reifyS(vv)
		return walkStar(vv, r)
	}
}

// Reify reifies the input variable for the given input states.
// NOTE: this is not the same reify as in the paper, mKreify
func Reify(v string, ss []*State) []*ast.SExpr {
	return fmap(ReifyVarFromState(v), ss)
}

// MKReifys finds reifications for the first introduced var
// NOTE: the way we've set this up now, vX is a reserved keyword
func MKReifys(ss []*State) []*ast.SExpr {
	return fmap(ReifyIntVarFromState(0), ss)
}

// MKReifys finds reifications for the first introduced var
// NOTE: the way we've set this up now, vX is a reserved keyword
func MKReify(s *State) *ast.SExpr {
	return ReifyIntVarFromState(0)(s)
}
