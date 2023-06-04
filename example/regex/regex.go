package regex

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func EmptySet() *ast.SExpr {
	return ast.NewSymbol("∅")
}

func EmptyStr() *ast.SExpr {
	return ast.NewSymbol("ε")
}

func Char(c rune) *ast.SExpr {
	return ast.NewList(ast.NewSymbol("char"), ast.NewSymbol(string([]rune{c})))
}

func CharSymbol(c rune) *ast.SExpr {
	return ast.NewSymbol(string([]rune{c}))
}

func CharFromSExpr(c *ast.SExpr) *ast.SExpr {
	return ast.NewList(ast.NewSymbol("char"), c)
}

func Or(a, b *ast.SExpr) *ast.SExpr {
	return ast.NewList(ast.NewSymbol("or"), a, b)
}

func Concat(a, b *ast.SExpr) *ast.SExpr {
	return ast.NewList(ast.NewSymbol("concat"), a, b)
}

func Star(a *ast.SExpr) *ast.SExpr {
	return ast.NewList(ast.NewSymbol("star"), a)
}
