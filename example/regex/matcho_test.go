package regex

import (
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestMatchOEmptySetNil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(EmptySet(), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOEmptySetChar(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(EmptySet(), ast.Cons(CharSymbol('a'), nil), q)
		},
		EmptySet(),
	)
}

func TestMatchOEmptyStrNil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(EmptyStr(), nil, q)
		},
		EmptyStr(),
	)
}

func TestMatchOEmptyStrChar(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(EmptyStr(), ast.Cons(CharSymbol('a'), nil), q)
		},
		EmptySet(),
	)
}

func TestMatchOCharNil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Char('a'), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOCharAChar(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Char('a'), ast.Cons(CharSymbol('a'), nil), q)
		},
		EmptyStr(),
	)
}

func TestMatchOCharBChar(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Char('a'), ast.Cons(CharSymbol('b'), nil), q)
		},
		EmptySet(),
	)
}

func TestMatchOCharStr(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Char('a'), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('a'), nil)), q)
		},
		EmptySet(),
	)
}

func TestMatchOOrNil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOOrCharA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), CharSymbol('a'), q)
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), CharSymbol('b'), q)
		},
		EmptyStr(),
	)
}

func TestMatchOOrCharC(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Or(Char('a'), Char('b')), CharSymbol('c'), q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatNil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), nil, q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), CharSymbol('a'), q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('b'), nil)), q)
		},
		EmptyStr(),
	)
}

func TestMatchOConcatCharAC(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('c'), nil)), q)
		},
		EmptySet(),
	)
}

func TestMatchOConcatCharABC(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Concat(Char('a'), Char('b')),
				ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('b'), ast.Cons(CharSymbol('c'), nil))), q)
		},
		EmptySet(),
	)
}

func TestMatchOStarANil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Star(Char('a')), nil, q)
		},
		EmptyStr(),
	)
}

func TestMatchOStarABCharA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), CharSymbol('a'), q)
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
		func(q *ast.SExpr) comicro.Goal {
			return MatchO(Star(Concat(Char('a'), Char('b'))), ast.Cons(CharSymbol('a'), ast.Cons(CharSymbol('b'), nil)), q)
		},
		EmptyStr(),
	)
}

// func TestGenMatchO(t *testing.T) {
// if testing.Short() {
// 	return
// }
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	g := func(q *ast.SExpr) comicro.Goal {
// 		return MatchO(q, ast.Cons(CharSymbol('a'), nil), EmptyStr())
// 	}
// 	ss := comicro.RunStream(ctx, g)
// 	for {
// 		s, ok := <-ss
// 		if !ok {
// 			return
// 		}
// 		if s != nil {
// 			fmt.Printf("%s\n", s.String())
// 		}
// 	}
// }

// func TestGenSDerivOs(t *testing.T) {
// 	if testing.Short() {
// 		return
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	g := func(q *ast.SExpr) comicro.Goal {
// 		return SDerivOs(q, ast.Cons(CharSymbol('a'), nil), EmptyStr())
// 	}
// 	ss := comicro.RunStream(ctx, g)
// 	for {
// 		s, ok := <-ss
// 		if !ok {
// 			return
// 		}
// 		if s != nil {
// 			fmt.Printf("%s\n", s.String())
// 		}
// 	}
// }
