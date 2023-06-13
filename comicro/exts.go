package comicro

// exts either extends a substitution s with an association between the variable x and the value v ,
// or it produces (ok = false)
// if extending the substitution with the pair `(,x . ,v) would create a cycle.
func exts(x Var, v any, s Substitutions) (Substitutions, bool) {
	if occurs(x, v, s) {
		return nil, false
	}
	return s.AddKeyValue(x, v), true
}

func occurs(x Var, v any, s Substitutions) bool {
	if vvar, ok := GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return vvar == x
		}
	}
	return Any(v, func(a any) bool {
		return occurs(x, a, s)
	})
}
