package regex

import (
	"context"
	"fmt"
	"testing"

	. "github.com/awalterschulze/gominikanren/gomini"
)

func TestNullOEmptySet(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(EmptySet(), q)
		},
		EmptySet(),
	)
}

func TestNullOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(EmptyStr(), q)
		},
		EmptyStr(),
	)
}

func TestNullOChar(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Char('a'), q)
		},
		EmptySet(),
	)
}

func TestNullOOrNilNil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Or(EmptyStr(), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrNilA(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Or(EmptyStr(), Char('a')), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrANil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Or(Char('a'), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrAB(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Or(Char('a'), Char('b')), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatNilNil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Concat(EmptyStr(), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOConcatNilA(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Concat(EmptyStr(), Char('a')), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatANil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Concat(Char('a'), EmptyStr()), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatAB(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Concat(Char('a'), Char('b')), q)
		},
		EmptySet(),
	)
}

func TestNullOStar(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return NullO(Star(Char('a')), q)
		},
		EmptyStr(),
	)
}

func TestGenNullO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = SetMaxRoutines(ctx, 100)
	defer cancel()
	g := func(q *Regex) Goal {
		return NullO(q, EmptyStr())
	}
	s := NewState(CreateVarRegex)
	ss := Run(ctx, s, g)
	count := 0
	for {
		res, ok := <-ss
		if !ok {
			return
		}
		count++
		fmt.Printf("%s\n", res.(stringer).String())
		if count > 10 {
			return
		}
	}
}

func TestIsNullOEmptyStr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVarRegex)
	g := IsNullO(EmptyStr())
	ss := Run(ctx, s, func(q *Regex) Goal {
		return g
	})
	_, ok := <-ss
	if !ok {
		t.Errorf("IsNullO failed")
	}
}

func TestIsNullOEmptySet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVarRegex)
	g := IsNullO(EmptySet())
	ss := Run(ctx, s, func(q *Regex) Goal {
		return g
	})
	_, ok := <-ss
	if ok {
		t.Errorf("expected IsNullO to fail")
	}
}

func TestIsNullOOrNilA(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVarRegex)
	g := IsNullO(Or(EmptyStr(), Char('a')))
	ss := Run(ctx, s, func(q *Regex) Goal {
		return g
	})
	_, ok := <-ss
	if !ok {
		t.Errorf("IsNullO failed")
	}
}
