package regex

import (
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestSimplOEmptySet(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(EmptySet(), q)
		},
		EmptySet(),
	)
}

func TestSimplOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(EmptyStr(), q)
		},
		EmptyStr(),
	)
}

func TestSimplOChar(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(Char('a'), q)
		},
		Char('a'),
	)
}

func TestSimplOOrAEmptySet(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(Or(EmptySet(), Char('a')), q)
		},
		Char('a'),
	)
}

func TestSimplOOrEmptySetA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(Or(Char('a'), EmptySet()), q)
		},
		Char('a'),
	)
}

func TestSimplOConcatAEmptySet(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(Concat(EmptySet(), Char('a')), q)
		},
		EmptySet(),
	)
}

func TestSimplOConcatEmptySetA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(Concat(Char('a'), EmptySet()), q)
		},
		EmptySet(),
	)
}

func TestSimplOConcatAEmptyStr(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(Concat(EmptyStr(), Char('a')), q)
		},
		Char('a'),
	)
}

func TestSimplOConcatEmptyStrA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return SimplO(Concat(Char('a'), EmptyStr()), q)
		},
		Char('a'),
	)
}