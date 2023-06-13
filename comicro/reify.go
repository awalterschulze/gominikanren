package comicro

import (
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func reifyName(n int) *ast.SExpr {
	return ast.NewSymbol("_" + strconv.Itoa(n))
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
	return Fold(vv, s, func(subs Substitutions, a any) Substitutions {
		sexpr, ok := a.(*ast.SExpr)
		if !ok {
			return s
		}
		return reifys(sexpr, subs)
	})
}

// ReifyIntVarFromState is a curried function that reifies the input variable for the given input state.
func ReifyIntVarFromState(v Var) func(s *State) *ast.SExpr {
	return func(s *State) *ast.SExpr {
		vv := walkStar(v.SExpr(), s.Substitutions)
		r := reifys(vv, nil)
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
