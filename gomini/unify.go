package gomini

import (
	"reflect"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(x, y any, s *State) (*State, bool) {
	var xx any = x
	if xvar, ok := s.CastVar(x); ok {
		xx = Lookup(xvar, s)
	}
	var yy any = y
	if yyar, ok := s.CastVar(y); ok {
		yy = Lookup(yyar, s)
	}
	if reflect.DeepEqual(xx, yy) {
		return s, true
	}
	if xvar, ok := xx.(Var); ok {
		return exts(xvar, yy, s)
	}
	if yyar, ok := yy.(Var); ok {
		return exts(yyar, xx, s)
	}
	ss, sok := ZipReduce(xx, yy, s, unify)
	if !sok {
		return nil, false
	}
	return ss, true
}

// exts either extends a substitution s with an association between the variable x and the value v ,
// or it produces (ok = false)
// if extending the substitution with the pair `(,x . ,v) would create a cycle.
func exts(x Var, v any, s *State) (*State, bool) {
	if occurs(x, v, s) {
		return nil, false
	}
	return s.AddKeyValue(x, v), true
}

func occurs(x Var, v any, s *State) bool {
	if vvar, ok := s.CastVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return reflect.DeepEqual(x, vvar)
		}
	}
	return Any(v, func(a any) bool {
		return occurs(x, a, s)
	})
}
