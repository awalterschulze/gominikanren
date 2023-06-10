package regex

import (
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestNullOEmptySet(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(EmptySet(), q)
		},
		EmptySet(),
	)
}

func TestNullOEmptyStr(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(EmptyStr(), q)
		},
		EmptyStr(),
	)
}

func TestNullOChar(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Char('a'), q)
		},
		EmptySet(),
	)
}

func TestNullOOrNilNil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(EmptyStr(), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrNilA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(EmptyStr(), Char('a')), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrANil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(Char('a'), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(Char('a'), Char('b')), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatNilNil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(EmptyStr(), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOConcatNilA(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(EmptyStr(), Char('a')), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatANil(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(Char('a'), EmptyStr()), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatAB(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(Char('a'), Char('b')), q)
		},
		EmptySet(),
	)
}

func TestNullOStar(t *testing.T) {
	testo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Star(Char('a')), q)
		},
		EmptyStr(),
	)
}

// func TestGenNullO(t *testing.T) {
// if testing.Short() {
// 	return
// }
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	g := func(q *ast.SExpr) comicro.Goal {
// 		return NullO(q, EmptyStr())
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
