package mini_test

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
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

func BenchmarkFib10Seq(b *testing.B) {
	benchFib(b, 10)
}

func BenchmarkFib15Seq(b *testing.B) {
	benchFib(b, 15)
}

// peano numbers and extralogical convenience functions
var zero = ast.NewInt(0)
var one = makenat(1)

func succ(prev, next *ast.SExpr) micro.Goal {
	return micro.EqualO(next, ast.Cons(prev, nil))
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

func natplus(x, y, z *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, zero), micro.EqualO(y, z)},
		[]micro.Goal{
			micro.CallFresh(func(a *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(b *ast.SExpr) micro.Goal {
					return mini.ConjPlus(
						succ(a, x),
						succ(b, z),
						natplus(a, y, b),
					)
				})
			}),
		},
	)
}

func fib(x, y *ast.SExpr) micro.Goal {
	return mini.Conde(
		[]micro.Goal{micro.EqualO(x, zero), micro.EqualO(y, zero)},
		[]micro.Goal{micro.EqualO(x, one), micro.EqualO(y, one)},
		[]micro.Goal{
			micro.CallFresh(func(n1 *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(n2 *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(f1 *ast.SExpr) micro.Goal {
						return micro.CallFresh(func(f2 *ast.SExpr) micro.Goal {
							return mini.ConjPlus(
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
		},
	)
}

func runFib(n int) []*ast.SExpr {
	return micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return fib(makenat(n), q)
	})
}
