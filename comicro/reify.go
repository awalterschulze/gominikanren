package comicro

import (
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func reifyName(n int) *ast.SExpr {
	return ast.NewSymbol("_" + strconv.Itoa(n))
}

// reifys expects a value and initially an empty substitution.
func reifyS(v *ast.SExpr) Substitutions {
	return reifys(v, nil)
}

func reifys(v *ast.SExpr, s Substitutions) Substitutions {
	vv := v
	if v.IsVariable() {
		vv = walk(NewVar(v.Atom.Var.Index), s)
	}
	if vv.IsVariable() {
		n := reifyName(len(s))
		return s.AddPair(NewVar(vv.Atom.Var.Index), n)
	}
	if vv.IsPair() {
		car := vv.Car()
		cars := reifys(car, s)
		cdr := vv.Cdr()
		return reifys(cdr, cars)
	}
	return s
}

// ReifyIntVarFromState is a curried function that reifies the input variable for the given input state.
func ReifyIntVarFromState(v Var) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(&ast.SExpr{Atom: &ast.Atom{Var: &ast.Variable{Index: uint64(v)}}}, s.Substitutions)
		r := reifyS(vv)
		return walkStar(vv, r)
	}
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
