package regex

import (
	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func SDerivO(r, char, out *ast.SExpr) comicro.Goal {
	return comini.Disjs(
		comicro.Conj(
			comicro.EqualO(r, EmptySet()),
			comicro.EqualO(out, EmptySet()),
		),
		comicro.Conj(
			comicro.EqualO(r, EmptyStr()),
			comicro.EqualO(out, EmptySet()),
		),
		DeriveCharO(r, char, out),
		comicro.Fresh(4, func(vars ...comicro.Var) comicro.Goal {
			a, da, b, db := vars[0].SExpr(), vars[1].SExpr(), vars[2].SExpr(), vars[3].SExpr()
			return comini.Conjs(
				comicro.EqualO(r, Or(a, b)),
				SDerivO(a, char, da),
				SDerivO(b, char, db),
				SimpleOrO(da, db, out),
			)
		}),
		comicro.Fresh(7, func(vars ...comicro.Var) comicro.Goal {
			a, da, na, b, db, ca, cb := vars[0].SExpr(), vars[1].SExpr(), vars[2].SExpr(), vars[3].SExpr(), vars[4].SExpr(), vars[5].SExpr(), vars[6].SExpr()
			return comini.Conjs(
				comicro.EqualO(r, Concat(a, b)),
				SDerivO(a, char, da),
				SDerivO(b, char, db),
				NullO(a, na),
				SimpleConcatO(da, b, ca),
				SimpleConcatO(na, db, cb),
				SimpleOrO(ca, cb, out),
			)
		}),
		comicro.Fresh(2, func(vars ...comicro.Var) comicro.Goal {
			a, da := vars[0].SExpr(), vars[1].SExpr()
			return comini.Conjs(
				comicro.EqualO(r, Star(a)),
				SDerivO(a, char, da),
				SimpleConcatO(da, r, out),
			)
		}),
	)
}
