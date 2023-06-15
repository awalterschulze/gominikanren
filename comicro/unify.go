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
	switch uu := uu.(type) {
	case Var:
		switch vv := vv.(type) {
		case Var:
			if uu == vv {
				return s, true
			}
			return exts(uu, vv.SExpr(), s)
		default:
			return exts(uu, vv, s)
		}
	default:
		switch vv := vv.(type) {
		case Var:
			return exts(vv, uu, s)
		default:
			if IsContainer(uu) && IsContainer(vv) {
				ss, sok := ZipFold(uu, vv, s, unify)
				if !sok {
					return nil, false
				}
				return ss, true
			}
		}
	}
	if reflect.DeepEqual(uu, vv) {
		return s, true
	}
	return nil, false
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
