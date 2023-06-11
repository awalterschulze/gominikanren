package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(u, v *ast.SExpr, s Substitutions) (Substitutions, bool) {
	uu := u
	if uvar, ok := GetVar(u); ok {
		uu = walk(uvar, s)
	}
	vv := v
	if vvar, ok := GetVar(v); ok {
		vv = walk(vvar, s)
	}
	uuvar, uok := GetVar(uu)
	vvvar, vok := GetVar(vv)
	if uok && vok && uuvar == vvvar {
		return s, true
	}
	if uuvar, ok := GetVar(uu); ok {
		return exts(uuvar, vv, s)
	}
	if vvar, ok := GetVar(vv); ok {
		return exts(vvar, uu, s)
	}
	if uu.IsPair() && vv.IsPair() {
		scar, sok := unify(uu.Car(), vv.Car(), s)
		if !sok {
			return nil, false
		}
		return unify(uu.Cdr(), vv.Cdr(), scar)
	}
	if uu.Equal(vv) {
		return s, true
	}
	return nil, false
}

func Unify(s *State, u, v *ast.SExpr) *State {
	substitutions, ok := unify(u, v, s.Substitutions)
	if !ok {
		return nil
	}
	return &State{Substitutions: substitutions, Counter: s.Counter}
}
