package regex

import (
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestSDerivOEmptySet(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(EmptySet(), CharSymbol('a'), q)
		},
		EmptySet(),
	)
}

func TestSDerivOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(EmptyStr(), CharSymbol('a'), q)
		},
		EmptySet(),
	)
}

func TestSDerivOCharA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(Char('a'), CharSymbol('a'), q)
		},
		EmptyStr(),
	)
}

func TestSDerivOCharB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(Char('a'), CharSymbol('b'), q)
		},
		EmptySet(),
	)
}

func TestSDerivOOrAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(Or(Char('a'), Char('b')), CharSymbol('a'), q)
		},
		EmptyStr(),
	)
}

func TestSDerivOOrNilA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(Or(EmptyStr(), Char('a')), CharSymbol('a'), q)
		},
		EmptyStr(),
	)
}

func TestSDerivOConcatAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(Concat(Char('a'), Char('b')), CharSymbol('a'), q)
		},
		Char('b'),
	)
}

func TestSDerivOStarA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(Star(Char('a')), CharSymbol('a'), q)
		},
		Star(Char('a')),
	)
}

func TestSDerivOStarAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SDerivO(Star(Concat(Char('a'), Char('b'))), CharSymbol('a'), q)
		},
		Concat(
			Char('b'),
			Star(Concat(Char('a'), Char('b'))),
		),
	)
}
