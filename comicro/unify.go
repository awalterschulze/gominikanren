package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(u, v *ast.SExpr, s Substitutions) (Substitutions, bool) {
	uu := u
	if u.IsVariable() {
		uu = walk(u.Atom.Var, s)
	}
	vv := v
	if v.IsVariable() {
		vv = walk(v.Atom.Var, s)
	}
	if uu.IsVariable() && vv.IsVariable() && uu.Atom.Var.Equal(uu.Atom.Var) {
		return s, true
	}
	if uu.IsVariable() {
		return exts(uu.Atom.Var, vv, s)
	}
	if vv.IsVariable() {
		return exts(vv.Atom.Var, uu, s)
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
