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
// the equivalent prolog program runs instantaneously
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
	// I don't think miniKanren _should_ work like that..
	a := func() *ast.SExpr {
		return ast.NewVariable(time.Now().String())
	}

	streetDef := ast.NewList(
		// each list represents a house on the street
		// anonymous vars are nationality, color, pet, drink, smoke, in that order
		ast.NewList(one, a(), a(), a(), a(), a()),
		ast.NewList(two, a(), a(), a(), a(), a()),
		ast.NewList(three, a(), a(), a(), a(), a()),
		ast.NewList(four, a(), a(), a(), a(), a()),
		ast.NewList(five, a(), a(), a(), a(), a()),
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
		return ConjPlus(
			streetEq(street),
			mUnrolled(ast.NewList(a(), ast.NewSymbol("brit"), ast.NewSymbol("red"), a(), a(), a())),
			mUnrolled(ast.NewList(a(), ast.NewSymbol("swede"), a(), ast.NewSymbol("dog"), a(), a())),
			mUnrolled(ast.NewList(a(), ast.NewSymbol("dane"), a(), a(), ast.NewSymbol("tea"), a())),
			mUnrolled(ast.NewList(ast.NewVariable("A"), a(), ast.NewSymbol("green"), a(), a(), a())),
			mUnrolled(ast.NewList(ast.NewVariable("B"), a(), ast.NewSymbol("white"), a(), a(), a())),
			leftOf(ast.NewVariable("A"), ast.NewVariable("B")),
			mUnrolled(ast.NewList(a(), a(), ast.NewSymbol("green"), a(), ast.NewSymbol("coffee"), a())),
			mUnrolled(ast.NewList(a(), a(), a(), ast.NewSymbol("birds"), a(), ast.NewSymbol("pall-mall"))),
			mUnrolled(ast.NewList(a(), a(), ast.NewSymbol("yellow"), a(), a(), ast.NewSymbol("dunhill"))),
			mUnrolled(ast.NewList(three, a(), a(), a(), ast.NewSymbol("milk"), a())),
			mUnrolled(ast.NewList(one, ast.NewSymbol("norwegian"), a(), a(), a(), a())),
			mUnrolled(ast.NewList(ast.NewVariable("C"), a(), a(), a(), a(), ast.NewSymbol("blend"))),
			mUnrolled(ast.NewList(ast.NewVariable("D"), a(), a(), ast.NewSymbol("cats"), a(), a())),
			nextTo(ast.NewVariable("C"), ast.NewVariable("D")),
			mUnrolled(ast.NewList(ast.NewVariable("E"), a(), a(), ast.NewSymbol("horse"), a(), a())),
			mUnrolled(ast.NewList(ast.NewVariable("F"), a(), a(), a(), a(), ast.NewSymbol("dunhill"))),
			nextTo(ast.NewVariable("E"), ast.NewVariable("F")),
			mUnrolled(ast.NewList(a(), a(), a(), a(), ast.NewSymbol("beer"), ast.NewSymbol("bluemaster"))),
			mUnrolled(ast.NewList(a(), ast.NewSymbol("german"), a(), a(), a(), ast.NewSymbol("prince"))),
			mUnrolled(ast.NewList(ast.NewVariable("G"), ast.NewSymbol("norwegian"), a(), a(), a(), a())),
			mUnrolled(ast.NewList(ast.NewVariable("H"), a(), ast.NewSymbol("blue"), a(), a(), a())),
			nextTo(ast.NewVariable("G"), ast.NewVariable("H")),
			mUnrolled(ast.NewList(ast.NewVariable("I"), a(), a(), a(), a(), ast.NewSymbol("blend"))),
			mUnrolled(ast.NewList(ast.NewVariable("J"), a(), a(), a(), ast.NewSymbol("water"), a())),
			nextTo(ast.NewVariable("I"), ast.NewVariable("J")),
			mUnrolled(ast.NewList(a(), fishowner, a(), ast.NewSymbol("fish"), a(), a())),
			/* old version uses membero directly, looping over 'street' over and over again
			   MemberO(ast.NewList(a(), ast.NewSymbol("brit"), ast.NewSymbol("red"), a(), a(), a()), street),
			   MemberO(ast.NewList(a(), ast.NewSymbol("swede"), a(), ast.NewSymbol("dog"), a(), a()), street),
			   MemberO(ast.NewList(a(), ast.NewSymbol("dane"), a(), a(), ast.NewSymbol("tea"), a()), street),
			   MemberO(ast.NewList(ast.NewVariable("A"), a(), ast.NewSymbol("green"), a(), a(), a()), street),
			   MemberO(ast.NewList(ast.NewVariable("B"), a(), ast.NewSymbol("white"), a(), a(), a()), street),
			   leftOf(ast.NewVariable("A"), ast.NewVariable("B")),
			   MemberO(ast.NewList(a(), a(), ast.NewSymbol("green"), a(), ast.NewSymbol("coffee"), a()), street),
			   MemberO(ast.NewList(a(), a(), a(), ast.NewSymbol("birds"), a(), ast.NewSymbol("pall-mall")), street),
			   MemberO(ast.NewList(a(), a(), ast.NewSymbol("yellow"), a(), a(), ast.NewSymbol("dunhill")), street),
			   MemberO(ast.NewList(three, a(), a(), a(), ast.NewSymbol("milk"), a()), street),
			   MemberO(ast.NewList(one, ast.NewSymbol("norwegian"), a(), a(), a(), a()), street),
			   MemberO(ast.NewList(ast.NewVariable("C"), a(), a(), a(), a(), ast.NewSymbol("blend")), street),
			   MemberO(ast.NewList(ast.NewVariable("D"), a(), a(), ast.NewSymbol("cats"), a(), a()), street),
			   nextTo(ast.NewVariable("C"), ast.NewVariable("D")),
			   MemberO(ast.NewList(ast.NewVariable("E"), a(), a(), ast.NewSymbol("horse"), a(), a()), street),
			   MemberO(ast.NewList(ast.NewVariable("F"), a(), a(), a(), a(), ast.NewSymbol("dunhill")), street),
			   nextTo(ast.NewVariable("E"), ast.NewVariable("F")),
			   MemberO(ast.NewList(a(), a(), a(), a(), ast.NewSymbol("beer"), ast.NewSymbol("bluemaster")), street),
			   MemberO(ast.NewList(a(), ast.NewSymbol("german"), a(), a(), a(), ast.NewSymbol("prince")), street),
			   MemberO(ast.NewList(ast.NewVariable("G"), ast.NewSymbol("norwegian"), a(), a(), a(), a()), street),
			   MemberO(ast.NewList(ast.NewVariable("H"), a(), ast.NewSymbol("blue"), a(), a(), a()), street),
			   nextTo(ast.NewVariable("G"), ast.NewVariable("H")),
			   MemberO(ast.NewList(ast.NewVariable("I"), a(), a(), a(), a(), ast.NewSymbol("blend")), street),
			   MemberO(ast.NewList(ast.NewVariable("J"), a(), a(), a(), ast.NewSymbol("water"), a()), street),
			   nextTo(ast.NewVariable("I"), ast.NewVariable("J")),
			   MemberO(ast.NewList(a(), fishowner, a(), ast.NewSymbol("fish"), a(), a()), street),
			*/
		)
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
