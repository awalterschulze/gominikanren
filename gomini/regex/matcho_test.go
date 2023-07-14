package regex

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/awalterschulze/gominikanren/gomini"
)

func TestMatchOEmptySetNil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(EmptySet(), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOEmptySetChar(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(EmptySet(), NewString("a"), q)
		},
		EmptySet(),
	)
}

func TestMatchOEmptyStrNil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(EmptyStr(), nil, q)
		},
		EmptyStr(),
	)
}

func TestMatchOEmptyStrChar(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(EmptyStr(), NewString("a"), q)
		},
		EmptySet(),
	)
}

func TestMatchOCharNil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Char('a'), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOCharAChar(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Char('a'), NewString("a"), q)
		},
		EmptyStr(),
	)
}

func TestMatchOCharBChar(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Char('a'), NewString("b"), q)
		},
		EmptySet(),
	)
}

func TestMatchOCharStr(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Char('a'), NewString("aa"), q)
		},
		EmptySet(),
	)
}

func TestMatchOOrNil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Or(Char('a'), Char('b')), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOOrCharA(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Or(Char('a'), Char('b')), NewString("a"), q)
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharB(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Or(Char('a'), Char('b')), NewString("b"), q)
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharAAB(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Or(Char('a'), Char('a')), NewString("b"), q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatNil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharA(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), NewString("a"), q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharAB(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), NewString("ab"), q)
		},
		EmptyStr(),
	)
}

func TestMatchOConcatCharAAB(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('a')), NewString("ab"), q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharABC(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), NewString("abc"), q)
		},
		EmptySet(),
	)
}

func TestMatchOStarANil(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Star(Char('a')), nil, q)
		},
		EmptyStr(),
	)
}

func TestMatchOStarABCharA(t *testing.T) {
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), NewString("a"), q)
		},
		EmptySet(),
	)
}

func TestMatchOStarABCharAB(t *testing.T) {
	if testing.Short() {
		return
	}
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), NewString("ab"), q)
		},
		EmptyStr(),
	)
}

func TestGenSDerivOs(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = SetMaxRoutines(ctx, 100)
	defer cancel()
	g := func(q *Regex) Goal {
		return SDerivOs(q, NewString("a"), EmptyStr())
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
		if count > 0 {
			return
		}
	}
}

func TestGenStrIsMatchOStarAOrB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := time.Now()
	s := NewState(CreateVarRegex)
	ctx = SetMaxRoutines(ctx, 100)
	g := func(q *String) Goal {
		return IsMatchO(Star(Or(Char('a'), Char('b'))), q)
	}
	ss := Run(ctx, s, g)
	count := 0
	for {
		res, ok := <-ss
		if !ok {
			return
		}
		count++
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("Generated, %v, %d, %s\n", elapsed, count, res.(stringer).String())
		if count >= 10 {
			return
		}
	}
}

func TestStrIsMatchO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = SetMaxRoutines(ctx, 100)
	defer cancel()
	g := func(q *String) Goal {
		return IsMatchO(Star(Char('a')), NewString("ab11122233334b"))
	}
	s := NewState(CreateVarRegex)
	ss := Run(ctx, s, g)
	if res, ok := <-ss; ok {
		t.Fatalf("did expect result %s\n", res.(stringer).String())
	}
}

func TestGenIsMatchO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = SetMaxRoutines(ctx, 100)
	defer cancel()
	g := func(q *Regex) Goal {
		return IsMatchO(q, NewString("a"))
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
		if count > 0 {
			return
		}
	}
}

func TestIsMatchOStarABCharAB(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVarRegex)
	g := IsMatchO(Star(Concat(Char('a'), Char('b'))), NewString("ab"))
	ss := Run(ctx, s, func(q *Regex) Goal {
		return g
	})
	_, ok := <-ss
	if !ok {
		t.Errorf("IsMatchO failed")
	}
}

func TestIsMatchOStarABCharBA(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVarRegex)
	g := IsMatchO(Star(Concat(Char('a'), Char('b'))), NewString("ba"))
	ss := Run(ctx, s, func(q *Regex) Goal {
		return g
	})
	_, ok := <-ss
	if ok {
		t.Errorf("IsMatchO should have failed")
	}
}
