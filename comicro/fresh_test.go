package comicro

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
scheme code:

	(run-goal 1
		(call/fresh kiwi
			(lambda (fruit)
				(== plum fruit)
			)
		)
	)
*/
func TestFreshKiwi(t *testing.T) {
	ss := RunGoal(context.Background(), 1,
		CallFresh(&ast.SExpr{}, func(fruit *ast.SExpr) Goal {
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
	want := "({,v0: plum} . 1)"
	got := ss[0].String()
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
