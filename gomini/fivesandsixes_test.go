package gomini

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func fives(x *ast.SExpr) Goal {
	return DisjO(
		EqualO(
			x,
			ast.NewInt(5),
		),
		func(ctx context.Context, s *State, ss Stream) { fives(x)(ctx, s, ss) },
	)
}

func sixes(x *ast.SExpr) Goal {
	return DisjO(
		EqualO(
			x,
			ast.NewInt(6),
		),
		func(ctx context.Context, s *State, ss Stream) { sixes(x)(ctx, s, ss) },
	)
}

func TestFivesAndSixesSExpr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ast5 := ast.NewInt(5)
	s := NewState(ast.CreateVar)

	// ((call/fresh (λ (q) (≡ q 5))) empty-state)
	gots := RunTake(ctx,
		1,
		s,
		func(q *ast.SExpr) Goal {
			return EqualO(
				q,
				ast5,
			)
		})
	got := fmap(gots, func(a any) *ast.SExpr {
		return a.(*ast.SExpr)
	})
	want := []*ast.SExpr{ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v but got %v", want, got)
	}

	// (define (fives x) (disj (≡ x 5) (fives x)))
	// ((call/fresh fives) empty-state)
	gots = RunTake(ctx,
		2,
		s,
		func(x *ast.SExpr) Goal {
			return fives(x)
		},
	)
	got = fmap(gots, func(a any) *ast.SExpr {
		return a.(*ast.SExpr)
	})
	want = []*ast.SExpr{ast5, ast5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v but got %#v", want, got)
	}

	stream := Run(ctx,
		s,
		func(x *ast.SExpr) Goal {
			return DisjO(
				fives(x),
				sixes(x),
			)
		},
	)
	has5 := false
	has6 := false
	for s := range stream {
		if !has5 && strings.Contains(s.(stringer).String(), "5") {
			has5 = true
		}
		if !has6 && strings.Contains(s.(stringer).String(), "6") {
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

func TestFivesAsInt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	five := int64(5)
	pointToFive := &five
	s := NewState()
	gots := RunTake(ctx,
		1,
		s,
		func(q *int64) Goal {
			return EqualO(
				q,
				pointToFive,
			)
		})
	got := fmap(gots, func(a any) *int64 {
		return a.(*int64)
	})
	want := []*int64{pointToFive}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v but got %v", want, got)
	}
}
