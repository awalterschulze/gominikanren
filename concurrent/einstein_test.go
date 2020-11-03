package concurrent

import (
	"testing"
	"time"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func BenchmarkEinsteinSeqZzz(b *testing.B) {
	disjunction = mini.DisjPlus
	goal := einstein()
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = runEinstein(goal)
	}
	result = r
}

func BenchmarkEinsteinSeq(b *testing.B) {
	disjunction = disjPlus
	goal := einstein()
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = runEinstein(goal)
	}
	result = r
}

func BenchmarkEinsteinConcZzz(b *testing.B) {
	disjunction = concDisjPlusZzz
	goal := einstein()
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = runEinstein(goal)
	}
	result = r
}

func BenchmarkEinsteinConc(b *testing.B) {
	disjunction = DisjPlus
	goal := einstein()
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = runEinstein(goal)
	}
	result = r
}

func disjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := disjPlus(gs[1:]...)
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			g1s := g1()(s)
			g2s := g2()(s)
			return micro.Mplus(g1s, g2s)
		}
	}
}

func concDisjPlusZzz(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return micro.Zzz(gs[0])
	}
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			list := make([]micro.StreamOfStates, len(gs))
			ch := make(chan answer)
			for i, g := range gs {
				go func(index int, goal micro.Goal) {
					ss := micro.Zzz(goal)()(s)
					ch <- answer{i: index, s: ss}
				}(i, g)
			}
			for range gs {
				ans := <-ch
				list[ans.i] = ans.s
			}
			stream := list[len(gs)-1]
			for i := len(gs) - 2; i >= 0; i-- {
				stream = micro.Mplus(list[i], stream)
			}
			return stream
		}
	}
}

func memberOUnrolled(y *ast.SExpr) func(*ast.SExpr) micro.Goal {
	goals := []func(*ast.SExpr) micro.Goal{}
	for {
		if y == nil {
			break
		}
		car, cdr := y.Car(), y.Cdr()
		goals = append(goals, func(x *ast.SExpr) micro.Goal {
			return micro.EqualO(x, car)
		})
		y = cdr
	}
	n := len(goals)
	return func(x *ast.SExpr) micro.Goal {
		gs := make([]micro.Goal, n)
		for i, g := range goals {
			gs[i] = g(x)
		}
		return disjunction(gs...)

	}
}

func einstein() func(*ast.SExpr, *ast.SExpr) micro.Goal {
	one := ast.NewSymbol("one")
	two := ast.NewSymbol("two")
	three := ast.NewSymbol("three")
	four := ast.NewSymbol("four")
	five := ast.NewSymbol("five")

	rightOf := func(x, y *ast.SExpr) micro.Goal {
		return disjunction(
			micro.ConjunctionO(
				micro.EqualO(y, one),
				micro.EqualO(x, two),
			),
			micro.ConjunctionO(
				micro.EqualO(y, two),
				micro.EqualO(x, three),
			),
			micro.ConjunctionO(
				micro.EqualO(y, three),
				micro.EqualO(x, four),
			),
			micro.ConjunctionO(
				micro.EqualO(y, four),
				micro.EqualO(x, five),
			),
		)
	}

	leftOf := func(x, y *ast.SExpr) micro.Goal {
		return rightOf(y, x)
	}

	nextTo := func(x, y *ast.SExpr) micro.Goal {
		return micro.DisjointO(
			leftOf(x, y),
			rightOf(x, y),
		)
	}

	anon := func() *ast.SExpr {
		return ast.NewVariable(time.Now().String())
	}

	streetDef := ast.NewList(
		ast.NewList(one, anon(), anon(), anon(), anon(), anon()),
		ast.NewList(two, anon(), anon(), anon(), anon(), anon()),
		ast.NewList(three, anon(), anon(), anon(), anon(), anon()),
		ast.NewList(four, anon(), anon(), anon(), anon(), anon()),
		ast.NewList(five, anon(), anon(), anon(), anon(), anon()),
	)

	streetEq := func(q *ast.SExpr) micro.Goal {
		return micro.EqualO(
			q,
			streetDef,
		)
	}

	mUnrolled := memberOUnrolled(streetDef)

	solution := func(street, fishowner, a, b, c, d, e, f, g, h, i, j *ast.SExpr) micro.Goal {
		return mini.ConjPlus(
			streetEq(street),
			mUnrolled(ast.NewList(anon(), ast.NewSymbol("brit"), ast.NewSymbol("red"), anon(), anon(), anon())),
			mUnrolled(ast.NewList(anon(), ast.NewSymbol("swede"), anon(), ast.NewSymbol("dog"), anon(), anon())),
			mUnrolled(ast.NewList(anon(), ast.NewSymbol("dane"), anon(), anon(), ast.NewSymbol("tea"), anon())),
			mUnrolled(ast.NewList(a, anon(), ast.NewSymbol("green"), anon(), anon(), anon())),
			mUnrolled(ast.NewList(b, anon(), ast.NewSymbol("white"), anon(), anon(), anon())),
			leftOf(a, b),
			mUnrolled(ast.NewList(anon(), anon(), ast.NewSymbol("green"), anon(), ast.NewSymbol("coffee"), anon())),
			mUnrolled(ast.NewList(anon(), anon(), anon(), ast.NewSymbol("birds"), anon(), ast.NewSymbol("pall-mall"))),
			mUnrolled(ast.NewList(anon(), anon(), ast.NewSymbol("yellow"), anon(), anon(), ast.NewSymbol("dunhill"))),
			mUnrolled(ast.NewList(three, anon(), anon(), anon(), ast.NewSymbol("milk"), anon())),
			mUnrolled(ast.NewList(one, ast.NewSymbol("norwegian"), anon(), anon(), anon(), anon())),
			mUnrolled(ast.NewList(c, anon(), anon(), anon(), anon(), ast.NewSymbol("blend"))),
			mUnrolled(ast.NewList(d, anon(), anon(), ast.NewSymbol("cats"), anon(), anon())),
			nextTo(c, d),
			mUnrolled(ast.NewList(e, anon(), anon(), ast.NewSymbol("horse"), anon(), anon())),
			mUnrolled(ast.NewList(f, anon(), anon(), anon(), anon(), ast.NewSymbol("dunhill"))),
			nextTo(e, f),
			mUnrolled(ast.NewList(anon(), anon(), anon(), anon(), ast.NewSymbol("beer"), ast.NewSymbol("bluemaster"))),
			mUnrolled(ast.NewList(anon(), ast.NewSymbol("german"), anon(), anon(), anon(), ast.NewSymbol("prince"))),
			mUnrolled(ast.NewList(g, ast.NewSymbol("norwegian"), anon(), anon(), anon(), anon())),
			mUnrolled(ast.NewList(h, anon(), ast.NewSymbol("blue"), anon(), anon(), anon())),
			nextTo(g, h),
			mUnrolled(ast.NewList(i, anon(), anon(), anon(), anon(), ast.NewSymbol("blend"))),
			mUnrolled(ast.NewList(j, anon(), anon(), anon(), ast.NewSymbol("water"), anon())),
			nextTo(i, j),
			mUnrolled(ast.NewList(anon(), fishowner, anon(), ast.NewSymbol("fish"), anon(), anon())),
		)
	}

	return func(street, fishowner *ast.SExpr) micro.Goal {
		return micro.CallFresh(func(a *ast.SExpr) micro.Goal {
			return micro.CallFresh(func(b *ast.SExpr) micro.Goal {
				return micro.CallFresh(func(c *ast.SExpr) micro.Goal {
					return micro.CallFresh(func(d *ast.SExpr) micro.Goal {
						return micro.CallFresh(func(e *ast.SExpr) micro.Goal {
							return micro.CallFresh(func(f *ast.SExpr) micro.Goal {
								return micro.CallFresh(func(g *ast.SExpr) micro.Goal {
									return micro.CallFresh(func(h *ast.SExpr) micro.Goal {
										return micro.CallFresh(func(i *ast.SExpr) micro.Goal {
											return micro.CallFresh(func(j *ast.SExpr) micro.Goal {
												return solution(street, fishowner, a, b, c, d, e, f, g, h, i, j)
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
	}
}

func runEinstein(goal func(*ast.SExpr, *ast.SExpr) micro.Goal) []*ast.SExpr {
	return micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return micro.CallFresh(func(street *ast.SExpr) micro.Goal {
			return micro.CallFresh(func(fishowner *ast.SExpr) micro.Goal {
				return micro.ConjunctionO(
					micro.EqualO(
						q,
						ast.Cons(street, ast.Cons(fishowner, nil)),
					),
					goal(street, fishowner),
				)
			})
		})
	})
}
