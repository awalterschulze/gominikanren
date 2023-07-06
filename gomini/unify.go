package gomini

import (
	"reflect"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(x, y any, s *State) *State {
	var xx any = x
	xx = walk(xx, s)
	var yy any = y
	yy = walk(yy, s)
	if reflect.DeepEqual(xx, yy) {
		return s
	}
	if xvar, ok := s.CastVar(xx); ok {
		return exts(xvar, yy, s)
	}
	if yyar, ok := s.CastVar(yy); ok {
		return exts(yyar, xx, s)
	}
	ss := ZipReduce(xx, yy, s, unify)
	if ss == nil {
		return nil
	}
	return ss
}

// walk returns the value of the variable in the substitutions map.
// If the value is a variable, it walks, until it finds a final value or variable without a substitution.
func walk(key any, s *State) any {
	kvar, ok := s.CastVar(key)
	if !ok {
		return key
	}
	value, ok := s.Get(kvar)
	if !ok {
		// There are no more substitutions to be made, this variable is the final value.
		return key
	}
	return walk(value, s)
}

// exts either extends a substitution s with an association between the variable x and the value v ,
// or it produces (ok = false)
// if extending the substitution with the pair `(,x . ,v) would create a cycle.
func exts(x Var, v any, s *State) *State {
	if occurs(x, v, s) {
		return nil
	}
	return s.Set(x, v)
}

func occurs(x Var, v any, s *State) bool {
	v = walk(v, s)
	if vvar, ok := s.CastVar(v); ok {
		return reflect.DeepEqual(x, vvar)
	}
	return Any(v, func(a any) bool {
		return occurs(x, a, s)
	})
}

// rewrite replaces all variables with their values, if it finds any in the substitutions map.
// It only replaces the variables that it finds on it's recursive walk, starting at the input variable.
func rewrite(v any, s *State) any {
	v = walk(v, s)
	if _, ok := s.CastVar(v); ok {
		return v
	}
	return Map(v, func(a any) any {
		return rewrite(a, s)
	})
}
