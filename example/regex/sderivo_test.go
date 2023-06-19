package regex

import (
	"context"
	"fmt"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
)

func TestSDerivOEmptySet(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(EmptySet(), &a, q)
		},
		EmptySet(),
	)
}

func TestSDerivOEmptyStr(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(EmptyStr(), &a, q)
		},
		EmptySet(),
	)
}

func TestSDerivOCharA(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(Char('a'), &a, q)
		},
		EmptyStr(),
	)
}

func TestSDerivOCharB(t *testing.T) {
	b := 'b'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(Char('a'), &b, q)
		},
		EmptySet(),
	)
}

func TestSDerivOOrAB(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(Or(Char('a'), Char('b')), &a, q)
		},
		EmptyStr(),
	)
}

func TestSDerivOOrNilA(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(Or(EmptyStr(), Char('a')), &a, q)
		},
		EmptyStr(),
	)
}

func TestSDerivOConcatAB(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(Concat(Char('a'), Char('b')), &a, q)
		},
		Char('b'),
	)
}

func TestSDerivOStarA(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(Star(Char('a')), &a, q)
		},
		Star(Char('a')),
	)
}

func TestSDerivOStarAB(t *testing.T) {
	a := 'a'
	testo(
		t,
		func(q *Regex) comicro.Goal {
			return SDerivO(Star(Concat(Char('a'), Char('b'))), &a, q)
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
	g := func(q *rune) comicro.Goal {
		return SDerivO(Char('a'), q, EmptyStr())
	}
	ss := comicro.RunStream(ctx, VarCreator, g)
	count := 0
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			break
		}
		count++
		r := s.(*rune)
		fmt.Printf("%c\n", *r)
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
	g := func(q *rune) comicro.Goal {
		return SDerivO(Char('a'), q, EmptySet())
	}
	ss := comicro.RunStream(ctx, VarCreator, g)
	count := 0
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			break
		}
		count++
		r := s.(*rune)
		fmt.Printf("%c\n", *r)
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
	g := func(q *rune) comicro.Goal {
		return SDerivO(Or(Char('a'), Char('b')), q, EmptyStr())
	}
	ss := comicro.RunStream(ctx, VarCreator, g)
	count := 0
	for {
		s, ok := comicro.ReadNonNilFromStream(ctx, ss)
		if !ok {
			break
		}
		count++
		r := s.(*rune)
		fmt.Printf("%c\n", *r)
	}
	if count < 2 {
		t.Fatalf("expected at least 2 solutions, got %d", count)
	}
}
