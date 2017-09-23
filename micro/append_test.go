package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"

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
func TestAppendO(t *testing.T) {
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
	want := "(`((,y . #t)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestCons(t *testing.T) {
	tests := []func() (string, string, string){
		deriveTuple3("a", "(c o r n)", "(a c o r n)"),
		deriveTuple3("a", "(n)", "(a n)"),
		deriveTuple3("a", "()", "(a)"),
	}
	for _, test := range tests {
		car, cdr, want := test()
		t.Run(want, func(t *testing.T) {
			care, err := sexpr.Parse(car)
			if err != nil {
				t.Fatal(err)
			}
			cdre, err := sexpr.Parse(cdr)
			if err != nil {
				t.Fatal(err)
			}
			l := cons(care, cdre)
			got := l.String()
			if got != want {
				t.Fatalf("got %v != want %v", got, want)
			}
		})
	}
}
