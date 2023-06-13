package comicro

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
func ReplaceAll(v any, s Substitutions) any {
	if vvar, ok := GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return vvar.SExpr()
		}
	}
	return Map(v, func(a any) any {
		return ReplaceAll(a, s)
	})
}
