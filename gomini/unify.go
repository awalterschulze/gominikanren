package gomini

import (
	"reflect"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(x, y any, s *State) *State {
	var xx any = x
	if xvar, ok := s.CastVar(x); ok {
		xx = walk(xvar, s)
	}
	var yy any = y
	if yyar, ok := s.CastVar(y); ok {
		yy = walk(yyar, s)
	}
	if reflect.DeepEqual(xx, yy) {
		return s
	}
	if xvar, ok := xx.(Var); ok {
		return exts(xvar, yy, s)
	}
	if yyar, ok := yy.(Var); ok {
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
func walk(key Var, s *State) any {
	value, ok := s.Get(key)
	if !ok {
		// There are no more substitutions to be made, this variable is the final value.
		return key
	}
	if newvar, ok := s.CastVar(value); ok {
		// If this value is a variable, then we need to Lookup it.
		return walk(newvar, s)
	}
	return value
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
	if vvar, ok := s.CastVar(v); ok {
		v = walk(vvar, s)
		if vvar, ok := v.(Var); ok {
			return reflect.DeepEqual(x, vvar)
		}
	}
	return Any(v, func(a any) bool {
		return occurs(x, a, s)
	})
}

// rewrite replaces all variables with their values, if it finds any in the substitutions map.
// It only replaces the variables that it finds on it's recursive walk, starting at the input variable.
func rewrite(v any, s *State) any {
	if vvar, ok := s.CastVar(v); ok {
		v = walk(vvar, s)
		if vvar, ok := v.(Var); ok {
			return s.GetPlaceHolder(vvar)
		}
	}
	return Map(v, func(a any) any {
		return rewrite(a, s)
	})
}
