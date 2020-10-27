package mini

import (
	"testing"
	"time"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
Let us assume that there are five houses of different colors next to each other on the same road.
In each house lives a man of a different nationality.
Every man has his favorite drink, his favorite brand of cigarettes, and keeps pets of a particular kind.
The Englishman lives in the red house.
The Swede keeps dogs.
The Dane drinks tea.
The green house is just to the left of the white one.
The owner of the green house drinks coffee.
The Pall Mall smoker keeps birds.
The owner of the yellow house smokes Dunhills.
The man in the center house drinks milk.
The Norwegian lives in the first house.
The Blend smoker has a neighbor who keeps cats.
The man who smokes Blue Masters drinks bier.
The man who keeps horses lives next to the Dunhill smoker.
The German smokes Prince.
The Norwegian lives next to the blue house.
The Blend smoker has a neighbor who drinks water.
The question to be answered is: Who keeps fish?
Answer adapted from https://gist.github.com/jgilchrist/1b85b04f4395057972dd
*/

// takes almost a minute (52s) to run atm on my machine with MemberO
// takes about 5s to run with MemberO unrolled, _without any concurrency support_
// the equivalent prolog and minikanren program runs virtually instantaneously
func TestEinstein(t *testing.T) {
	// not implementing counting right now, so numbers are symbols :)
	one := ast.NewSymbol("one")
	two := ast.NewSymbol("two")
	three := ast.NewSymbol("three")
	four := ast.NewSymbol("four")
	five := ast.NewSymbol("five")

	rightOf := func(x, y *ast.SExpr) micro.Goal {
		return DisjPlus(
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

	// NOTE: these should be anonymous vars, but we cannot make those atm without using call/fresh
	// I'm creating vars with unique random names (from unix time) instead, for now
	// somehow ast.NewVariable can introduce vars with scope: if you change one of the var names in
	// street to duplicate another, you will see them substituted in variable bindings in state
	anon := func() *ast.SExpr {
		return ast.NewVariable(time.Now().String())
	}

	streetDef := ast.NewList(
		// each list represents a house on the street
		// anonymous vars are nationality, color, pet, drink, smoke, in that order
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

	mUnrolled := MemberOUnrolled(streetDef)

	// abusing what I noted above here; using named vars to introduce them as fresh
	// because I don't want to write 10x nested callfresh and the fresh macro is eluding me
	solution := func(street, fishowner *ast.SExpr) micro.Goal {
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
												return ConjPlus(
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

	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return micro.CallFresh(func(street *ast.SExpr) micro.Goal {
			return micro.CallFresh(func(fishowner *ast.SExpr) micro.Goal {
				return micro.ConjunctionO(
					micro.EqualO(
						q,
						ast.Cons(street, ast.Cons(fishowner, nil)),
					),
					solution(street, fishowner),
				)
			})
		})
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(sexprs))
	}
	streetSolved := sexprs[0].Car()
	fishOwner := sexprs[0].Cdr().Car()
	if fishOwner.String() != "german" {
		t.Fatalf("expected %s, but got %s instead", "german", fishOwner.String())
	}
	wantFullSolution := "((one norwegian yellow cats water dunhill) (two dane blue horse tea blend) (three brit red birds milk pall-mall) (four german green fish coffee prince) (five swede white dog beer bluemaster))"
	if streetSolved.String() != wantFullSolution {
		t.Fatalf("expected %s, but got %s instead", wantFullSolution, streetSolved.String())
	}
}
