package regex

import (
	"context"
	"fmt"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
)

func TestSimplOEmptySet(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(EmptySet(), q.SExpr())
		},
		EmptySet(),
	)
}

func TestSimplOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(EmptyStr(), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestSimplOChar(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(Char('a'), q.SExpr())
		},
		Char('a'),
	)
}

func TestSimplOOrAEmptySet(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(Or(EmptySet(), Char('a')), q.SExpr())
		},
		Char('a'),
	)
}

func TestSimplOOrEmptySetA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(Or(Char('a'), EmptySet()), q.SExpr())
		},
		Char('a'),
	)
}

func TestSimplOConcatAEmptySet(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(Concat(EmptySet(), Char('a')), q.SExpr())
		},
		EmptySet(),
	)
}

func TestSimplOConcatEmptySetA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(Concat(Char('a'), EmptySet()), q.SExpr())
		},
		EmptySet(),
	)
}

func TestSimplOConcatAEmptyStr(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(Concat(EmptyStr(), Char('a')), q.SExpr())
		},
		Char('a'),
	)
}

func TestSimplOConcatEmptyStrA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return SimplO(Concat(Char('a'), EmptyStr()), q.SExpr())
		},
		Char('a'),
	)
}

func TestGenSimplOA(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q comicro.Var) comicro.Goal {
		return SimplO(q.SExpr(), Char('a'))
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
