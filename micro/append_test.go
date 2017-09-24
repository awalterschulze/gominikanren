package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(let
	(
		(q (var 'q))
	)
	(map
		(reify q)
		(run-goal #f
			(call/fresh 'x
				(lambda (x)
					(call/fresh 'y
						(lambda (y)
							(conj
								(== `(,x ,y) q)
								(appendo x y `(cake & ice d t))
							)
						)
					)
				)
			)
		)
	)
)
*/
// results in all the combinations of two lists that when appended will result in (cake & ice d t)
func TestAppendOAllCombinations(t *testing.T) {
	subs := RunGoal(
		-1,
		CallFresh("x", func(x *ast.SExpr) Goal {
			return CallFresh("y", func(y *ast.SExpr) Goal {
				return ConjunctionO(
					EqualO(
						ast.NewList(true, x, y),
						ast.NewVariable("q"),
					),
					AppendO(
						x,
						y,
						ast.NewList(true,
							ast.NewSymbol("cake"),
							ast.NewSymbol("&"),
							ast.NewSymbol("ice"),
							ast.NewSymbol("d"),
							ast.NewSymbol("t"),
						),
					),
				)
			})
		}),
	)
	r := reify(ast.NewVariable("q"))
	ss := deriveFmapR(r, subs)
	got := ast.NewList(false, ss...).String()
	want := ""
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestAppendOSingleList(t *testing.T) {
	subs := RunGoal(
		-1,
		AppendO(
			ast.NewList(false,
				ast.NewSymbol("a"),
			),
			ast.NewList(false,
				ast.NewSymbol("b"),
			),
			ast.NewVariable("q"),
		),
	)
	r := reify(ast.NewVariable("q"))
	ss := deriveFmapR(r, subs)
	got := ast.NewList(false, ss...).String()
	want := "(`(a b))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestAppendOSingleAtom(t *testing.T) {
	subs := RunGoal(
		-1,
		AppendO(
			ast.NewList(false,
				ast.NewSymbol("a"),
			),
			ast.NewSymbol("b"),
			ast.NewVariable("q"),
		),
	)
	r := reify(ast.NewVariable("q"))
	ss := deriveFmapR(r, subs)
	got := ast.NewList(false, ss...).String()
	want := "(`(a b))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestCarO(t *testing.T) {
	ifte := IfThenElseO(
		CarO(
			ast.NewList(false,
				ast.NewSymbol("a"),
				ast.NewSymbol("c"),
				ast.NewSymbol("o"),
				ast.NewSymbol("r"),
				ast.NewSymbol("n"),
			),
			ast.NewSymbol("a"),
		),
		EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
		EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
	)
	ss := ifte(EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #t) (,d . (c o r n))))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
