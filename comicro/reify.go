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
	if vvar, ok := GetVar(v); ok {
		vv = walk(vvar, s)
	}
	if vvar, ok := GetVar(vv); ok {
		n := reifyName(len(s))
		return s.AddPair(vvar, n)
	}
	if vv.IsPair() {
		return Fold(vv.Pair, s, func(subs Substitutions, a any) Substitutions {
			sexpr := a.(*ast.SExpr)
			return reifys(sexpr, subs)
		})
	}
	return s
}

// ReifyIntVarFromState is a curried function that reifies the input variable for the given input state.
func ReifyIntVarFromState(v Var) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(v.SExpr(), s.Substitutions)
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
