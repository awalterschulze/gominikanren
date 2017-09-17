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
		CallFresh("kiwi",
			func(fruit *ast.SExpr) Goal {
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
	want := "`((,kiwi . plum))"
	got := ss[0].String()
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
