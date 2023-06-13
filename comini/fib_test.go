package comini_test

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/comini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

var result []*ast.SExpr

func benchFib(b *testing.B, num int, f func(...comicro.Goal) comicro.Goal) {
	var r []*ast.SExpr
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = runFib(num, f)
	}
	result = r
}

func BenchmarkFib10Co(b *testing.B) {
	benchFib(b, 10, comini.Conjs)
}

func BenchmarkFib15Co(b *testing.B) {
	benchFib(b, 15, comini.Conjs)
}

// peano numbers and extralogical convenience functions
var zero = ast.NewInt(0)
var one = makenat(1)

func succ(prev, next *ast.SExpr) comicro.Goal {
	return comicro.EqualO(next, ast.Cons(prev, nil))
}

func makenat(n int) *ast.SExpr {
	if n == 0 {
		return zero
	}
	return ast.Cons(makenat(n-1), nil)
}

func parsenat(x *ast.SExpr) int {
	if x == zero {
		return 0
	}
	return parsenat(x.Car()) + 1
}

func toList(list []int) *ast.SExpr {
	var out *ast.SExpr = nil
	for i := 0; i < len(list); i++ {
		out = ast.Cons(makenat(list[len(list)-i-1]), out)
	}
	return out
}

func natplus(x, y, z *ast.SExpr) comicro.Goal {
	return comini.Conde(
		[]comicro.Goal{comicro.EqualO(x, zero), comicro.EqualO(y, z)},
		[]comicro.Goal{
			comicro.CallFresh(func(a comicro.Var) comicro.Goal {
				return comicro.CallFresh(func(b comicro.Var) comicro.Goal {
					return comini.Conjs(
						succ(a.SExpr(), x),
						succ(b.SExpr(), z),
						natplus(a.SExpr(), y, b.SExpr()),
					)
				})
			}),
		},
	)
}

func fib(conj func(...comicro.Goal) comicro.Goal, x, y *ast.SExpr) comicro.Goal {
	return comini.Conde(
		[]comicro.Goal{comicro.EqualO(x, zero), comicro.EqualO(y, zero)},
		[]comicro.Goal{comicro.EqualO(x, one), comicro.EqualO(y, one)},
		[]comicro.Goal{
			comicro.Fresh(4, func(vars ...comicro.Var) comicro.Goal {
				n1, n2, f1, f2 := vars[0], vars[1], vars[2], vars[3]
				return conj(
					succ(n1.SExpr(), x),
					succ(n2.SExpr(), n1.SExpr()),
					fib(conj, n1.SExpr(), f1.SExpr()),
					fib(conj, n2.SExpr(), f2.SExpr()),
					natplus(f1.SExpr(), f2.SExpr(), y),
				)
			}),
		},
	)
}

func runFib(n int, f func(...comicro.Goal) comicro.Goal) []*ast.SExpr {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return fmap(comicro.Run(ctx, -1, func(q comicro.Var) comicro.Goal {
		return fib(f, makenat(n), q.SExpr())
	}), func(x any) *ast.SExpr {
		return x.(*ast.SExpr)
	})
}

func fmap[A, B any](list []A, f func(A) B) []B {
	out := make([]B, len(list))
	for i, elem := range list {
		out[i] = f(elem)
	}
	return out
}
