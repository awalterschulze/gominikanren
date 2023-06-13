package regex

import (
	"context"
	"fmt"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
)

func TestDerivOEmptySet(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(EmptySet(), CharSymbol('a'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestDerivOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(EmptyStr(), CharSymbol('a'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestDerivOCharA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(Char('a'), CharSymbol('a'), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestDerivOCharB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(Char('a'), CharSymbol('b'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestDerivOOrAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(Or(Char('a'), Char('b')), CharSymbol('a'), q.SExpr())
		},
		Or(EmptyStr(), EmptySet()),
	)
}

func TestDerivOOrNilA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(Or(EmptyStr(), Char('a')), CharSymbol('a'), q.SExpr())
		},
		Or(EmptySet(), EmptyStr()),
	)
}

func TestDerivOConcatAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(Concat(Char('a'), Char('b')), CharSymbol('a'), q.SExpr())
		},
		Or(Concat(EmptyStr(), Char('b')), Concat(EmptySet(), EmptySet())),
	)
}

func TestDerivOStarA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(Star(Char('a')), CharSymbol('a'), q.SExpr())
		},
		Concat(EmptyStr(), Star(Char('a'))),
	)
}

func TestDerivOStarAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return DerivO(Star(Concat(Char('a'), Char('b'))), CharSymbol('a'), q.SExpr())
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
	g := func(q comicro.Var) comicro.Goal {
		return DerivO(Char('a'), q.SExpr(), EmptyStr())
	}
	ss := comicro.RunStream(ctx, g)
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
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
	g := func(q comicro.Var) comicro.Goal {
		return DerivO(Char('a'), q.SExpr(), EmptySet())
	}
	ss := comicro.RunStream(ctx, g)
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
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
	g := func(q comicro.Var) comicro.Goal {
		return DerivO(Or(Char('a'), Char('b')), q.SExpr(), Or(EmptySet(), EmptyStr()))
	}
	ss := comicro.RunStream(ctx, g)
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			return
		}
		fmt.Printf("%s\n", s.String())
	}
}
