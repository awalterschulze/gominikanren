package micro

import (
    "reflect"
    "testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func fives(x *ast.SExpr) GoalFn {
	return Zzz(DisjointO(
		EqualO(
			ast.NewVariable("x"),
			ast.NewInt(5),
		),
		func() GoalFn {return fives(x)},
	))()
}

func sixes(x *ast.SExpr) GoalFn {
	return Zzz(DisjointO(
		EqualO(
			ast.NewVariable("x"),
			ast.NewInt(6),
		),
		func() GoalFn {return sixes(x)},
	))()
}

func TestFivesAndSixes(t *testing.T) {
    ast5 := ast.NewInt(5)
    ast6 := ast.NewInt(6)

	// ((call/fresh (λ (q) (≡ q 5))) empty-state)
	states := RunGoal(
		1,
		CallFresh(func(q *ast.SExpr) Goal {
			return EqualO(
				ast.NewVariable("q"),
				ast.NewInt(5),
			)
		}),
	)
	got := Reify("q", states)
    want := []*ast.SExpr{ast5}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("expected %#v but got %#v", got, want)
    }

	// (define (fives x) (disj (≡ x 5) (fives x)))
	// ((call/fresh fives) empty-state)
	states = RunGoal(
		2,
		CallFresh(func(x *ast.SExpr) Goal {
			return func() GoalFn {return fives(x)}
		}),
	)
	got = Reify("x", states)
    want = []*ast.SExpr{ast5, ast5}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("expected %#v but got %#v", got, want)
    }

	states = RunGoal(
		10,
		CallFresh(func(x *ast.SExpr) Goal {
			return DisjointO(
				func() GoalFn {return fives(x)},
				func() GoalFn {return sixes(x)},
			)
		}),
	)
	got = Reify("x", states)
    want = []*ast.SExpr{ast5, ast6, ast5, ast6, ast5, ast6, ast5, ast6, ast5, ast6}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("expected %#v but got %#v", got, want)
    }
}
