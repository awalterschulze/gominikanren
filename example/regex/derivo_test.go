package regex

import (
	"context"
	"fmt"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
)

func TestDerivOEmptySet(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(EmptySet(), &a, q)
		},
		EmptySet(),
	)
}

func TestDerivOEmptyStr(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(EmptyStr(), &a, q)
		},
		EmptySet(),
	)
}

func TestDerivOCharA(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(Char('a'), &a, q)
		},
		EmptyStr(),
	)
}

func TestDerivOCharB(t *testing.T) {
	b := rune('b')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(Char('a'), &b, q)
		},
		EmptySet(),
	)
}

func TestDerivOOrAB(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(Or(Char('a'), Char('b')), &a, q)
		},
		Or(EmptyStr(), EmptySet()),
	)
}

func TestDerivOOrNilA(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(Or(EmptyStr(), Char('a')), &a, q)
		},
		Or(EmptySet(), EmptyStr()),
	)
}

func TestDerivOConcatAB(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(Concat(Char('a'), Char('b')), &a, q)
		},
		Or(Concat(EmptyStr(), Char('b')), Concat(EmptySet(), EmptySet())),
	)
}

func TestDerivOStarA(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(Star(Char('a')), &a, q)
		},
		Concat(EmptyStr(), Star(Char('a'))),
	)
}

func TestDerivOStarAB(t *testing.T) {
	a := rune('a')
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return DerivO(Star(Concat(Char('a'), Char('b'))), &a, q)
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
	g := func(q *rune) comicro.Goal {
		return DerivO(Char('a'), q, EmptyStr())
	}
	ss := comicro.RunStream(ctx, &runeType, VarCreator, g)
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			return
		}
		rptr := s.(*rune)
		fmt.Printf("%c\n", rune(*rptr))
	}
}

var runeType = rune(0)

func TestGenDeriveOB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q *rune) comicro.Goal {
		return DerivO(Char('a'), q, EmptySet())
	}
	ss := comicro.RunStream(ctx, &runeType, VarCreator, g)
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			return
		}
		rptr := s.(*rune)
		fmt.Printf("%c\n", rune(*rptr))
	}
}

func TestGenDeriveOAOrB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := func(q *rune) comicro.Goal {
		return DerivO(Or(Char('a'), Char('b')), q, Or(EmptySet(), EmptyStr()))
	}
	ss := comicro.RunStream(ctx, &runeType, VarCreator, g)
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			return
		}
		rptr := s.(*rune)
		fmt.Printf("%c\n", rune(*rptr))
	}
}
