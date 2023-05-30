package comicro

import (
	"reflect"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func fives(x *ast.SExpr) Goal {
	return Zzz(Disj(
		EqualO(
			x,
			ast.NewInt(5),
		),
		func(s *State) StreamOfStates { return fives(x)(s) },
	))
}

func sixes(x *ast.SExpr) Goal {
	return Zzz(Disj(
		EqualO(
			x,
			ast.NewInt(6),
		),
		func(s *State) StreamOfStates { return sixes(x)(s) },
	))
}

func TestFivesAndSixes(t *testing.T) {
	ast5 := ast.NewInt(5)

	// ((call/fresh (λ (q) (≡ q 5))) empty-state)
	states := RunGoal(
		1,
		CallFresh(func(q *ast.SExpr) Goal {
			return EqualO(
				q,
				ast5,
			)
		}))
	got := MKReify(states)
	want := []*ast.SExpr{ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", got, want)
	}

	// (define (fives x) (disj (≡ x 5) (fives x)))
	// ((call/fresh fives) empty-state)
	states = RunGoal(
		2,
		CallFresh(func(x *ast.SExpr) Goal {
			return fives(x)
		}),
	)
	got = MKReify(states)
	want = []*ast.SExpr{ast5, ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", got, want)
	}

	states = RunGoal(
		10,
		CallFresh(func(x *ast.SExpr) Goal {
			return Disj(
				fives(x),
				sixes(x),
			)
		}),
	)
	got = MKReify(states)
	ss := fmap(func(s *ast.SExpr) string { return s.String() }, got)
	s := strings.Join(ss, ",	")
	if !strings.Contains(s, "5") {
		t.Fatalf("expected %s to contain at least one 5", s)
	}
	if !strings.Contains(s, "6") {
		t.Fatalf("expected %s to contain at least one 6", s)
	}

}
