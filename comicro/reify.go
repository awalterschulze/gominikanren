package comicro

import (
	"strconv"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func reifys(v any, s Substitutions) Substitutions {
	if vvar, ok := GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			n := ast.NewSymbol("_" + strconv.Itoa(len(s)))
			return s.AddKeyValue(vvar, n)
		}
	}
	return Fold(v, s, func(subs Substitutions, a any) Substitutions {
		return reifys(a, subs)
	})
}

// reifyFromState is a curried function that reifies the input variable for the given input state.
func reifyFromState(v Var, s *State) *ast.SExpr {
	vv := ReplaceAll(v, s.Substitutions)
	r := reifys(vv, nil)
	return ReplaceAll(vv, r).(*ast.SExpr)
}

// Reify finds reifications for the first introduced var
// NOTE: the way we've set this up now, vX is a reserved keyword
func Reify(s *State) *ast.SExpr {
	return reifyFromState(0, s)
}
