package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// walk has been modified to assume it is getting a variable, but here is the original scheme code:
func walk(v *ast.Variable, s Substitutions) *ast.SExpr {
	value, ok := assv(v, s)
	if !ok {
		return &ast.SExpr{Atom: &ast.Atom{Var: v}}
	}
	if !value.IsVariable() {
		return value
	}
	return walk(value.Atom.Var, s)
}

// assv either produces the first association in s that has v as its car using eqv,
// or produces ok = false if l has no such association.
func assv(v *ast.Variable, ss Substitutions) (*ast.SExpr, bool) {
	if ss == nil {
		return nil, false
	}
	return ss.Get(NewVar(v.Index))
}

func walkStar(v *ast.SExpr, s Substitutions) *ast.SExpr {
	vv := v
	if v.IsVariable() {
		vv = walk(v.Atom.Var, s)
	}
	if vv.IsVariable() {
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
