package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func walk(v Var, s Substitutions) *ast.SExpr {
	value, ok := assv(v, s)
	if !ok {
		return v.SExpr()
	}
	if !IsVar(value) {
		return value
	}
	return walk(NewVar(value.Atom.Var.Index), s)
}

// assv either produces the first association in s that has v as its car using eqv,
// or produces ok = false if l has no such association.
func assv(v Var, ss Substitutions) (*ast.SExpr, bool) {
	return ss.Get(v)
}

func walkStar(v *ast.SExpr, s Substitutions) *ast.SExpr {
	vv := v
	if IsVar(v) {
		vv = walk(NewVar(v.Atom.Var.Index), s)
	}
	if IsVar(vv) {
		return vv
	}
	if vv.IsPair() {
		carv := vv.Pair.Car
		cdrv := vv.Pair.Cdr
		wcar := walkStar(carv, s)
		wcdr := walkStar(cdrv, s)
		w := ast.Cons(wcar, wcdr)
		return w
	}
	return vv
}
