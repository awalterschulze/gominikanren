package gomini_test

import (
	"context"
	"testing"

	. "github.com/awalterschulze/gominikanren/gomini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

var result []*ast.SExpr

func benchFib(b *testing.B, num int) {
	var r []*ast.SExpr
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = runFib(num)
	}
	result = r
}

func BenchmarkFib10Co(b *testing.B) {
	benchFib(b, 10)
}

func BenchmarkFib15Co(b *testing.B) {
	benchFib(b, 15)
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
	return DisjO(
		ConjO(EqualO(x, zero), EqualO(y, z)),
		ConjO(
			ExistO(func(a *ast.SExpr) Goal {
				return ExistO(func(b *ast.SExpr) Goal {
					return ConjO(
						succ(a, x),
						succ(b, z),
						natplus(a, y, b),
					)
				})
			}),
		),
	)
}

func fib(x, y *ast.SExpr) Goal {
	return DisjO(
		ConjO(EqualO(x, zero), EqualO(y, zero)),
		ConjO(EqualO(x, one), EqualO(y, one)),
		ConjO(
			ExistO(func(n1 *ast.SExpr) Goal {
				return ExistO(func(n2 *ast.SExpr) Goal {
					return ExistO(func(f1 *ast.SExpr) Goal {
						return ExistO(func(f2 *ast.SExpr) Goal {
							return ConjO(
								succ(n1, x),
								succ(n2, n1),
								fib(n1, f1),
								fib(n2, f2),
								natplus(f1, f2, y),
							)
						})
					})
				})
			}),
		),
	)
}

func runFib(n int) []*ast.SExpr {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(ast.CreateVar)
	return fmap(RunTake(ctx, -1, s, func(q *ast.SExpr) Goal {
		return fib(makenat(n), q)
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
