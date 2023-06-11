package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// unify returns either (ok = false) or the substitution s extended with zero or more associations,
// where cycles in substitutions can lead to (ok = false)
func unify(u, v *ast.SExpr, s Substitutions) (Substitutions, bool) {
	uu := u
	if IsVar(u) {
		uu = walk(NewVar(u.Atom.Var.Index), s)
	}
	vv := v
	if IsVar(v) {
		vv = walk(NewVar(v.Atom.Var.Index), s)
	}
	if IsVar(uu) && IsVar(vv) && uu.Atom.Var.Equal(uu.Atom.Var) {
		return s, true
	}
	if IsVar(uu) {
		return exts(NewVar(uu.Atom.Var.Index), vv, s)
	}
	if IsVar(vv) {
		return exts(NewVar(vv.Atom.Var.Index), uu, s)
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
