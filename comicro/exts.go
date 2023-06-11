package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// exts either extends a substitution s with an association between the variable x and the value v ,
// or it produces (ok = false)
// if extending the substitution with the pair `(,x . ,v) would create a cycle.
func exts(x Var, v *ast.SExpr, s Substitutions) (Substitutions, bool) {
	if occurs(x, v, s) {
		return nil, false
	}
	return s.AddPair(x, v), true
}

func occurs(x Var, v *ast.SExpr, s Substitutions) bool {
	vv := v
	if IsVar(v) {
		vv = walk(NewVar(v.Atom.Var.Index), s)
	}
	if IsVar(vv) {
		return NewVar(vv.Atom.Var.Index) == x
	}
	if vv.IsPair() {
		return occurs(x, vv.Car(), s) || occurs(x, vv.Cdr(), s)
	}
	return false
}
