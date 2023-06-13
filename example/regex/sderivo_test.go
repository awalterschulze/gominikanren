package regex

import (
	"context"
	"fmt"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
)

func TestSDerivOEmptySet(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(EmptySet(), CharSymbol('a'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestSDerivOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(EmptyStr(), CharSymbol('a'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestSDerivOCharA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(Char('a'), CharSymbol('a'), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestSDerivOCharB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(Char('a'), CharSymbol('b'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestSDerivOOrAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(Or(Char('a'), Char('b')), CharSymbol('a'), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestSDerivOOrNilA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(Or(EmptyStr(), Char('a')), CharSymbol('a'), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestSDerivOConcatAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(Concat(Char('a'), Char('b')), CharSymbol('a'), q.SExpr())
		},
		Char('b'),
	)
}

func TestSDerivOStarA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(Star(Char('a')), CharSymbol('a'), q.SExpr())
		},
		Star(Char('a')),
	)
}

func TestSDerivOStarAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SDerivO(Star(Concat(Char('a'), Char('b'))), CharSymbol('a'), q.SExpr())
		},
		Concat(
			Char('b'),
			Star(Concat(Char('a'), Char('b'))),
		),
	)
}

func TestGenSDerivO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q comicro.Var) comicro.Goal {
		return SDerivO(Char('a'), q.SExpr(), EmptyStr())
	}
	ss := comicro.RunStream(ctx, g)
	count := 0
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			break
		}
		count++
		fmt.Printf("%s\n", s.String())
	}
	if count < 1 {
		t.Fatalf("expected at least 1 solution, got %d", count)
	}
}

func TestGenSDerivOB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q comicro.Var) comicro.Goal {
		return SDerivO(Char('a'), q.SExpr(), EmptySet())
	}
	ss := comicro.RunStream(ctx, g)
	count := 0
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			break
		}
		count++
		fmt.Printf("%s\n", s.String())
	}
	if count < 2 {
		t.Fatalf("expected at least 2 solutions, got %d", count)
	}
}

func TestGenSDerivOAOrB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q comicro.Var) comicro.Goal {
		return SDerivO(Or(Char('a'), Char('b')), q.SExpr(), EmptyStr())
	}
	ss := comicro.RunStream(ctx, g)
	count := 0
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			break
		}
		count++
		fmt.Printf("%s\n", s.String())
	}
	if count < 2 {
		t.Fatalf("expected at least 2 solutions, got %d", count)
	}
}
