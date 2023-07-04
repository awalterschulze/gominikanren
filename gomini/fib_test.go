package gomini_test

import (
	"context"
	"testing"

	. "github.com/awalterschulze/gominikanren/gomini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

var result []*ast.SExpr

func benchFib(b *testing.B, num int, f func(...Goal) Goal) {
	var r []*ast.SExpr
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = runFib(num, f)
	}
	result = r
}

func BenchmarkFib10Co(b *testing.B) {
	benchFib(b, 10, Conjs)
}

func BenchmarkFib15Co(b *testing.B) {
	benchFib(b, 15, Conjs)
}

// peano numbers and extralogical convenience functions
var zero = ast.NewInt(0)
var one = makenat(1)

func succ(prev, next *ast.SExpr) Goal {
	return EqualO(next, ast.Cons(prev, nil))
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

func natplus(x, y, z *ast.SExpr) Goal {
	return Disjs(
		Conjs(EqualO(x, zero), EqualO(y, z)),
		Conjs(
			Exists(func(a *ast.SExpr) Goal {
				return Exists(func(b *ast.SExpr) Goal {
					return Conjs(
						succ(a, x),
						succ(b, z),
						natplus(a, y, b),
					)
				})
			}),
		),
	)
}

func fib(conj func(...Goal) Goal, x, y *ast.SExpr) Goal {
	return Disjs(
		Conjs(EqualO(x, zero), EqualO(y, zero)),
		Conjs(EqualO(x, one), EqualO(y, one)),
		Conjs(
			Exists(func(n1 *ast.SExpr) Goal {
				return Exists(func(n2 *ast.SExpr) Goal {
					return Exists(func(f1 *ast.SExpr) Goal {
						return Exists(func(f2 *ast.SExpr) Goal {
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
		),
	)
}

func runFib(n int, f func(...Goal) Goal) []*ast.SExpr {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewEmptyState(ast.CreateVar)
	return fmap(Run(ctx, -1, s, func(q *ast.SExpr) Goal {
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
