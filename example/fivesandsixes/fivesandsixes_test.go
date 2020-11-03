package example

import (
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func fives(x *ast.SExpr) micro.GoalFn {
	return micro.Zzz(micro.DisjointO(
		micro.EqualO(
			x,
			ast.NewInt(5),
		),
		func() micro.GoalFn { return fives(x) },
	))()
}

func sixes(x *ast.SExpr) micro.GoalFn {
	return micro.Zzz(micro.DisjointO(
		micro.EqualO(
			x,
			ast.NewInt(6),
		),
		func() micro.GoalFn { return sixes(x) },
	))()
}

func TestFivesAndSixes(t *testing.T) {
	ast5 := ast.NewInt(5)
	ast6 := ast.NewInt(6)

	// ((call/fresh (λ (q) (≡ q 5))) empty-state)
	states := micro.RunGoal(
		1,
		micro.CallFresh(func(q *ast.SExpr) micro.Goal {
			return micro.EqualO(
				q,
				ast5,
			)
		}))
	got := micro.MKReify(states)
	want := []*ast.SExpr{ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", got, want)
	}

	// (define (fives x) (disj (≡ x 5) (fives x)))
	// ((call/fresh fives) empty-state)
	states = micro.RunGoal(
		2,
		micro.CallFresh(func(x *ast.SExpr) micro.Goal {
			return func() micro.GoalFn { return fives(x) }
		}),
	)
	got = micro.MKReify(states)
	want = []*ast.SExpr{ast5, ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", got, want)
	}

	states = micro.RunGoal(
		10,
		micro.CallFresh(func(x *ast.SExpr) micro.Goal {
			return micro.DisjointO(
				func() micro.GoalFn { return fives(x) },
				func() micro.GoalFn { return sixes(x) },
			)
		}),
	)
	got = micro.MKReify(states)
	want = []*ast.SExpr{ast5, ast6, ast5, ast6, ast5, ast6, ast5, ast6, ast5, ast6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", got, want)
	}
}
