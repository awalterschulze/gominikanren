package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// exts either extends a substitution s with an association between the variable x and the value v ,
// or it produces (ok = false)
// if extending the substitution with the pair `(,x . ,v) would create a cycle.
func exts(x Var, v *ast.SExpr, s Substitutions) (Substitutions, bool) {
	if occurs(x, v, s) {
		return nil, false
	}
	return s.AddPair(x, v), true
}

func occurs(x Var, v *ast.SExpr, s Substitutions) bool {
	vv := v
	if vvar, ok := GetVar(v); ok {
		vv = walk(vvar, s)
	}
	if vvar, ok := GetVar(vv); ok {
		return vvar == x
	}
	return Any(vv, func(a any) bool {
		sexpr, ok := a.(*ast.SExpr)
		if !ok {
			return false
		}
		return occurs(x, sexpr, s)
	})
}
