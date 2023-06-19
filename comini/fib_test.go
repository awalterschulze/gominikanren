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
			comicro.Exists(func(a *ast.SExpr) comicro.Goal {
				return comicro.Exists(func(b *ast.SExpr) comicro.Goal {
					return comini.Conjs(
						succ(a, x),
						succ(b, z),
						natplus(a, y, b),
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
			comicro.Exists(func(n1 *ast.SExpr) comicro.Goal {
				return comicro.Exists(func(n2 *ast.SExpr) comicro.Goal {
					return comicro.Exists(func(f1 *ast.SExpr) comicro.Goal {
						return comicro.Exists(func(f2 *ast.SExpr) comicro.Goal {
							return conj(
								succ(n1, x),
								succ(n2, n1),
								fib(conj, n1, f1),
								fib(conj, n2, f2),
								natplus(f1, f2, y),
							)
						})
					})
				})
			}),
		},
	)
}

func runFib(n int, f func(...comicro.Goal) comicro.Goal) []*ast.SExpr {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := comicro.NewEmptyState().WithVarCreators(ast.CreateVar)
	return fmap(comicro.Run(ctx, -1, s, func(q *ast.SExpr) comicro.Goal {
		return fib(f, makenat(n), q)
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
