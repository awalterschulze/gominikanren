package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// Lookup returns the value of the variable,
// if it has any substitutions it recurses,
// otherwise it returns the variable with no association or the final value.
func Lookup(key Var, s Substitutions) any {
	value, ok := assv(key, s)
	if !ok {
		// There are no more substitutions to be made, this variable is the final value.
		return key
	}
	if newvar, ok := GetVar(value); ok {
		// If this value is a variable, then we need to Lookup it.
		return Lookup(newvar, s)
	}
	return value
}

// assv either produces the first association in s that has v as its car using eqv,
// or produces ok = false if l has no such association.
func assv(v Var, ss Substitutions) (any, bool) {
	return ss.Get(v)
}

// Replace replaces all variables with their values, if it finds any in the substitutions.
func ReplaceAll(v *ast.SExpr, s Substitutions) *ast.SExpr {
	if vvar, ok := GetVar(v); ok {
		vvalue := Lookup(vvar, s)
		if vvar, ok := vvalue.(Var); ok {
			return vvar.SExpr()
		}
		v = vvalue.(*ast.SExpr)
	}
	vva := Map(v, func(a any) any {
		sexpr, ok := a.(*ast.SExpr)
		if !ok {
			return a
		}
		return ReplaceAll(sexpr, s)
	})
	return vva.(*ast.SExpr)
}
