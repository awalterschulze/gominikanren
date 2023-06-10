package comicro

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func fives(x *ast.SExpr) Goal {
	return Disj(
		EqualO(
			x,
			ast.NewInt(5),
		),
		func(ctx context.Context, s *State, ss StreamOfStates) { fives(x)(ctx, s, ss) },
	)
}

func sixes(x *ast.SExpr) Goal {
	return Disj(
		EqualO(
			x,
			ast.NewInt(6),
		),
		func(ctx context.Context, s *State, ss StreamOfStates) { sixes(x)(ctx, s, ss) },
	)
}

func TestFivesAndSixes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ast5 := ast.NewInt(5)

	// ((call/fresh (λ (q) (≡ q 5))) empty-state)
	states := RunGoal(ctx,
		1,
		CallFresh(func(q Var) Goal {
			return EqualO(
				q.SExpr(),
				ast5,
			)
		}))
	got := MKReifys(states)
	want := []*ast.SExpr{ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", want, got)
	}

	// (define (fives x) (disj (≡ x 5) (fives x)))
	// ((call/fresh fives) empty-state)
	states = RunGoal(ctx,
		2,
		CallFresh(func(x Var) Goal {
			return fives(x.SExpr())
		}),
	)
	got = MKReifys(states)
	want = []*ast.SExpr{ast5, ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", want, got)
	}

	stream := RunStream(ctx,
		func(x Var) Goal {
			return Disj(
				fives(x.SExpr()),
				sixes(x.SExpr()),
			)
		},
	)
	has5 := false
	has6 := false
	for s := range stream {
		if !has5 && strings.Contains(s.String(), "5") {
			has5 = true
		}
		if !has6 && strings.Contains(s.String(), "6") {
			has6 = true
		}
		if has5 && has6 {
			break
		}
	}
	if !has5 {
		t.Fatalf("expected to find least one 5")
	}
	if !has6 {
		t.Fatalf("expected to find least one 6")
	}
}
