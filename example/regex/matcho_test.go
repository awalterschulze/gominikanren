package regex

import (
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestMatchOEmptySetNil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(EmptySet(), nil, q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOEmptySetChar(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(EmptySet(), ast.Cons(CharSymbol('a'), nil), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOEmptyStrNil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(EmptyStr(), nil, q.SExpr())
		},
		EmptyStr(),
	)
}

func TestMatchOEmptyStrChar(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(EmptyStr(), ast.Cons(CharSymbol('a'), nil), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOCharNil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Char('a'), nil, q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOCharAChar(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Char('a'), ast.Cons(CharSymbol('a'), nil), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestMatchOCharBChar(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Char('a'), ast.Cons(CharSymbol('b'), nil), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOCharStr(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Char('a'), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('a'), nil)), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOOrNil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), nil, q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOOrCharA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), CharSymbol('a'), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), CharSymbol('b'), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharC(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), CharSymbol('c'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOConcatNil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), nil, q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), CharSymbol('a'), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharAB(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('b'), nil)), q.SExpr())
		},
		EmptyStr(),
	)
}

func TestMatchOConcatCharAC(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('c'), nil)), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharABC(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')),
				ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('b'), ast.Cons(CharSymbol('c'), nil))), q.SExpr())
		},
		EmptySet(),
	)
}

func TestMatchOStarANil(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Star(Char('a')), nil, q.SExpr())
		},
		EmptyStr(),
	)
}

func TestMatchOStarABCharA(t *testing.T) {
	testo(
		t,
		func(q comicro.Var) comicro.Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), CharSymbol('a'), q.SExpr())
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
		func(q comicro.Var) comicro.Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('b'), nil)), q.SExpr())
		},
		EmptyStr(),
	)
}

// func TestGenSDerivOs(t *testing.T) {
// 	if testing.Short() {
// 		return
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	ctx = comicro.SetMaxRoutines(ctx, 10000)
// 	defer cancel()
// 	g := func(q comicro.Var) comicro.Goal {
// 		return SDerivOs(q, ast.Cons(CharSymbol('a'), nil), EmptyStr())
// 	}
// 	ss := comicro.RunStream(ctx, g)
// 	count := 0
// 	for {
// 		s, ok := comicro.ReadNonNull(ctx, ss)
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
// 	ctx = comicro.SetMaxRoutines(ctx, 100)
// 	defer cancel()
// 	g := func(q comicro.Var) comicro.Goal {
// 		return MatchO(q, ast.Cons(CharSymbol('a'), nil), EmptyStr())
// 	}
// 	ss := comicro.RunStream(ctx, g)
// 	count := 0
// 	for {
// 		s, ok := comicro.ReadNonNull(ctx, ss)
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
