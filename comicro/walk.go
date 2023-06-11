package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func walk(v Var, s Substitutions) *ast.SExpr {
	value, ok := assv(v, s)
	if !ok {
		return v.SExpr()
	}
	if varvalue, ok := GetVar(value); ok {
		return walk(varvalue, s)
	}
	return value
}

// assv either produces the first association in s that has v as its car using eqv,
// or produces ok = false if l has no such association.
func assv(v Var, ss Substitutions) (*ast.SExpr, bool) {
	return ss.Get(v)
}

func walkStar(v *ast.SExpr, s Substitutions) *ast.SExpr {
	vv := v
	if vvar, ok := GetVar(v); ok {
		vv = walk(vvar, s)
	}
	if vvar, ok := GetVar(vv); ok {
		return vvar.SExpr()
	}
	if vv.IsPair() {
		vva := Apply(vv, func(a any) any {
			sexpr := a.(*ast.SExpr)
			return walkStar(sexpr, s)
		})
		return vva.(*ast.SExpr)
	}
	return vv
}
