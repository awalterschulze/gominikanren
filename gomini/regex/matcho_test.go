package regex

import (
	"testing"

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
	a := "a"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(EmptySet(), &a, q)
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
	a := "a"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(EmptyStr(), &a, q)
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
	a := "a"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Char('a'), &a, q)
		},
		EmptyStr(),
	)
}

func TestMatchOCharBChar(t *testing.T) {
	b := "b"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Char('a'), &b, q)
		},
		EmptySet(),
	)
}

func TestMatchOCharStr(t *testing.T) {
	aa := "aa"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Char('a'), &aa, q)
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
	a := "a"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Or(Char('a'), Char('b')), &a, q)
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharB(t *testing.T) {
	b := "b"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Or(Char('a'), Char('b')), &b, q)
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharC(t *testing.T) {
	c := "c"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Or(Char('a'), Char('b')), &c, q)
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
	a := "a"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), &a, q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharAB(t *testing.T) {
	ab := "ab"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), &ab, q)
		},
		EmptyStr(),
	)
}

func TestMatchOConcatCharAC(t *testing.T) {
	ac := "ac"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), &ac, q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharABC(t *testing.T) {
	abc := "abc"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Concat(Char('a'), Char('b')), &abc, q)
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
	a := "a"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), &a, q)
		},
		EmptySet(),
	)
}

func TestMatchOStarABCharAB(t *testing.T) {
	if testing.Short() {
		return
	}
	ab := "ab"
	testo(
		t,
		func(q *Regex) Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), &ab, q)
		},
		EmptyStr(),
	)
}

// func TestGenSDerivOs(t *testing.T) {
// 	if testing.Short() {
// 		return
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	ctx = SetMaxRoutines(ctx, 10000)
// 	defer cancel()
// 	g := func(q *Regex) Goal {
// 		return SDerivOs(q, ast.Cons(CharSymbol('a'), nil), EmptyStr())
// 	}
// 	ss := RunStream(ctx, g)
// 	count := 0
// 	for {
// 		s, ok := ReadNonNull(ctx, ss)
// 		if !ok {
// 			return
// 		}
// 		count++
// 		fmt.Printf("%s\n", s.String())
// 		if count > 10 {
// 			return
// 		}
// 	}
// }

// func TestGenMatchO(t *testing.T) {
// 	if testing.Short() {
// 		return
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	ctx = SetMaxRoutines(ctx, 100)
// 	defer cancel()
// 	g := func(q *Regex) Goal {
// 		return MatchO(q, ast.Cons(CharSymbol('a'), nil), EmptyStr())
// 	}
// 	ss := RunStream(ctx, g)
// 	count := 0
// 	for {
// 		s, ok := ReadNonNull(ctx, ss)
// 		if !ok {
// 			return
// 		}
// 		count++
// 		fmt.Printf("%s\n", s.String())
// 		if count > 10 {
// 			return
// 		}
// 	}
// }