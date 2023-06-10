package regex

import (
	"context"
	"fmt"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestDerivOEmptySet(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(EmptySet(), CharSymbol('a'), q)
		},
		EmptySet(),
	)
}

func TestDerivOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(EmptyStr(), CharSymbol('a'), q)
		},
		EmptySet(),
	)
}

func TestDerivOCharA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(Char('a'), CharSymbol('a'), q)
		},
		EmptyStr(),
	)
}

func TestDerivOCharB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(Char('a'), CharSymbol('b'), q)
		},
		EmptySet(),
	)
}

func TestDerivOOrAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(Or(Char('a'), Char('b')), CharSymbol('a'), q)
		},
		Or(EmptyStr(), EmptySet()),
	)
}

func TestDerivOOrNilA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(Or(EmptyStr(), Char('a')), CharSymbol('a'), q)
		},
		Or(EmptySet(), EmptyStr()),
	)
}

func TestDerivOConcatAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(Concat(Char('a'), Char('b')), CharSymbol('a'), q)
		},
		Or(Concat(EmptyStr(), Char('b')), Concat(EmptySet(), EmptySet())),
	)
}

func TestDerivOStarA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(Star(Char('a')), CharSymbol('a'), q)
		},
		Concat(EmptyStr(), Star(Char('a'))),
	)
}

func TestDerivOStarAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return DerivO(Star(Concat(Char('a'), Char('b'))), CharSymbol('a'), q)
		},
		Concat(
			Or(Concat(EmptyStr(), Char('b')), Concat(EmptySet(), EmptySet())),
			Star(Concat(Char('a'), Char('b'))),
		),
	)
}

func TestGenDeriveOA(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q *ast.SExpr) comicro.Goal {
		return DerivO(Char('a'), q, EmptyStr())
	}
	ss := comicro.RunStream(ctx, g)
	for {
		s, ok := comicro.ReadNonNull(ctx, ss)
		if !ok {
			return
		}
		fmt.Printf("%s\n", s.String())
	}
}

func TestGenDeriveOB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q *ast.SExpr) comicro.Goal {
		return DerivO(Char('a'), q, EmptySet())
	}
	ss := comicro.RunStream(ctx, g)
	for {
		s, ok := comicro.ReadNonNull(ctx, ss)
		if !ok {
			return
		}
		fmt.Printf("%s\n", s.String())
	}
}

func TestGenDeriveOAOrB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q *ast.SExpr) comicro.Goal {
		return DerivO(Or(Char('a'), Char('b')), q, Or(EmptySet(), EmptyStr()))
	}
	ss := comicro.RunStream(ctx, g)
	for {
		s, ok := comicro.ReadNonNull(ctx, ss)
		if !ok {
			return
		}
		fmt.Printf("%s\n", s.String())
	}
}
