package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(u, v *ast.SExpr, s Substitutions) (Substitutions, bool) {
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
		case *ast.SExpr:
			return exts(uu, vv, s)
		}
	case *ast.SExpr:
		switch vv := vv.(type) {
		case Var:
			return exts(vv, uu, s)
		case *ast.SExpr:
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
		}
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
