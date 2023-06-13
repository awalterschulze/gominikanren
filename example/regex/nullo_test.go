package regex

import (
	"context"
	"fmt"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
)

func TestNullOEmptySet(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(EmptySet(), q.SExpr())
		},
		EmptySet(),
	)
}

func TestNullOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(EmptyStr(), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestNullOChar(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Char('a'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestNullOOrNilNil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Or(EmptyStr(), EmptyStr()), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestNullOOrNilA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Or(EmptyStr(), Char('a')), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestNullOOrANil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Or(Char('a'), EmptyStr()), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestNullOOrAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Or(Char('a'), Char('b')), q.SExpr())
		},
		EmptySet(),
	)
}

func TestNullOConcatNilNil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Concat(EmptyStr(), EmptyStr()), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestNullOConcatNilA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Concat(EmptyStr(), Char('a')), q.SExpr())
		},
		EmptySet(),
	)
}

func TestNullOConcatANil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Concat(Char('a'), EmptyStr()), q.SExpr())
		},
		EmptySet(),
	)
}

func TestNullOConcatAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Concat(Char('a'), Char('b')), q.SExpr())
		},
		EmptySet(),
	)
}

func TestNullOStar(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return NullO(Star(Char('a')), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestGenNullO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = comicro.SetMaxRoutines(ctx, 100)
	defer cancel()
	g := func(q comicro.Var) comicro.Goal {
		return NullO(q.SExpr(), EmptyStr())
	}
	ss := comicro.RunStream(ctx, g)
	count := 0
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			return
		}
		count++
		fmt.Printf("%s\n", s.String())
		if count > 10 {
			return
		}
	}
}
