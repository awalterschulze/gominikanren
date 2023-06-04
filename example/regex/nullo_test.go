package regex

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func nullo(t *testing.T, f func(q *ast.SExpr) comicro.Goal, want *ast.SExpr) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sexprs := comicro.Run(ctx, 1, f)
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d result, but got %d instead", 1, len(sexprs))
	}
	got := sexprs[0]
	if !got.Equal(want) {
		t.Fatalf("expected %s, but got %s instead", want, got)
	}
}

func TestNullOEmptySet(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(EmptySet(), q)
		},
		EmptySet(),
	)
}

func TestNullOEmptyStr(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(EmptyStr(), q)
		},
		EmptyStr(),
	)
}

func TestNullOChar(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Char('a'), q)
		},
		EmptySet(),
	)
}

func TestNullOOrNilNil(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(EmptyStr(), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrNilA(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(EmptyStr(), Char('a')), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrANil(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(Char('a'), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOOrAB(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Or(Char('a'), Char('b')), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatNilNil(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(EmptyStr(), EmptyStr()), q)
		},
		EmptyStr(),
	)
}

func TestNullOConcatNilA(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(EmptyStr(), Char('a')), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatANil(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(Char('a'), EmptyStr()), q)
		},
		EmptySet(),
	)
}

func TestNullOConcatAB(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Concat(Char('a'), Char('b')), q)
		},
		EmptySet(),
	)
}

func TestNullOStar(t *testing.T) {
	nullo(
		t,
		func(q *ast.SExpr) comicro.Goal {
			return NullO(Star(Char('a')), q)
		},
		EmptyStr(),
	)
}

// func TestNullOGenEmptyStr(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	g := comicro.CallFresh(func(q *ast.SExpr) comicro.Goal {
// 		return NullO(q, EmptyStr())
// 	})
// 	ss := g(ctx, &comicro.State{Substitutions: nil, Counter: 1})
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

// func TestNullOGenEmptySet(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	g := comicro.CallFresh(func(q *ast.SExpr) comicro.Goal {
// 		return NullO(q, EmptySet())
// 	})
// 	ss := g(ctx, &comicro.State{Substitutions: nil, Counter: 1})
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
