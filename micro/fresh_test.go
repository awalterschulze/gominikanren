package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(run-goal 1
	(call/fresh kiwi
		(lambda (fruit)
			(== plum fruit)
		)
	)
)
*/
func TestFreshKiwi(t *testing.T) {
	ss := RunGoal(1,
		CallFresh(func(fruit *ast.SExpr) Goal {
			return EqualO(
				ast.NewSymbol("plum"),
				fruit,
			)
		},
		),
	)
	if len(ss) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(ss))
	}
	want := "((,v0 . plum) . 1)"
	got := ss[0].String()
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
