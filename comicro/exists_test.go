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
	s := NewEmptyState().WithVarCreators(ast.CreateVar)
	ss := RunGoal(context.Background(), 1, s,
		Exists(func(fruit *ast.SExpr) Goal {
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

func TestPineapple(t *testing.T) {
	pineapple := "pineapple"
	pointy := &pineapple
	reify := func(varTyp any, name string) (any, bool) {
		return &name, true
	}
	s := NewEmptyState().WithVarCreators(reify)
	ss := Run(context.Background(), 1, s,
		func(fruit *string) Goal {
			return EqualO(
				pointy,
				fruit,
			)
		},
	)
	if len(ss) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(ss))
	}
	got := *(ss[0]).(*string)
	want := pineapple
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
