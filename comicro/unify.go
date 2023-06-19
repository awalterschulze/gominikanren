package comicro

import (
	"reflect"
)

func Unify(s *State, x, y any) *State {
	s1, ok := unify(x, y, s)
	if !ok {
		return nil
	}
	return s1
}

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(u, v any, s *State) (*State, bool) {
	var uu any = u
	var vv any = v
	if uvar, ok := s.GetVar(u); ok {
		uu = Lookup(uvar, s)
	}
	if vvar, ok := s.GetVar(v); ok {
		vv = Lookup(vvar, s)
	}
	if reflect.DeepEqual(uu, vv) {
		return s, true
	}
	if uvar, ok := uu.(Var); ok {
		return exts(uvar, vv, s)
	}
	if vvar, ok := vv.(Var); ok {
		return exts(vvar, uu, s)
	}
	ss, sok := ZipFold(uu, vv, s, unify)
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
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return s.SameVar(vvar, x)
		}
	}
	return Any(v, func(a any) bool {
		return occurs(x, a, s)
	})
}
