package comicro

import (
	"reflect"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(u, v any, s Substitutions) (Substitutions, bool) {
	var uu any = u
	var vv any = v
	if uvar, ok := GetVar(u); ok {
		uu = Lookup(uvar, s)
	}
	if vvar, ok := GetVar(v); ok {
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

func Unify(s *State, x, y any) *State {
	// make sure untyped nil is converted to typed nil
	if !reflect.ValueOf(y).IsValid() {
		y = reflect.Zero(reflect.TypeOf(x)).Interface()
	}
	substitutions, ok := unify(x, y, s.Substitutions)
	if !ok {
		return nil
	}
	return &State{Substitutions: substitutions, Counter: s.Counter}
}
