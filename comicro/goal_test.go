package comicro

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestEqualO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tests := []func() (string, string, string){
		tuple3("#f", "#f", "((() . 0))"),
		tuple3("#f", "#t", "()"),
	}
	for _, test := range tests {
		u, v, want := test()
		t.Run("((== "+u+" "+v+") empty-s)", func(t *testing.T) {
			uexpr, err := sexpr.Parse(u)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			stream := EqualO(uexpr, vexpr)(ctx, EmptyState())
			got := stream.String()
			if got != want {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestFailureO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if got, want := FailureO(ctx, EmptyState()).String(), "()"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestSuccessO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if got, want := SuccessO(ctx, EmptyState()).String(), "((() . 0))"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNeverO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	n := NeverO(ctx, EmptyState())
	s, ok := <-n
	if s != nil {
		t.Fatalf("expected suspension")
	}
	if !ok {
		t.Fatalf("expected never ending")
	}
}

func TestDisj1(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	x := ast.NewVariable("x")
	d := Disj(
		EqualO(
			ast.NewSymbol("olive"),
			x,
		),
		NeverO,
	)(ctx, EmptyState())
	var s *State
	for s = range d {
		if s != nil {
			break
		}
	}
	got := s.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", indexOf(x)), "x", -1)
	want := "((,x . olive) . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	_, ok := <-d
	if !ok {
		t.Fatalf("expected never ending")
	}
}

func TestDisj2(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	x := ast.NewVariable("x")
	d := Disj(
		NeverO,
		EqualO(
			ast.NewSymbol("olive"),
			x,
		),
	)(ctx, EmptyState())
	for s := range d {
		if s != nil {
			got := s.String()
			// reifying y; we assigned it a random uint64 and lost track of it
			got = strings.Replace(got, fmt.Sprintf("v%d", indexOf(x)), "x", -1)
			want := "((,x . olive) . 0)"
			if got != want {
				t.Fatalf("got %s != want %s", got, want)
			}
			break
		}
	}
	_, ok := <-d
	if !ok {
		t.Fatalf("expected never ending")
	}

}

func TestAlwaysO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := AlwaysO(ctx, EmptyState())
	var s *State
	for s = range a {
		if s != nil {
			break
		}
	}
	got := s.String()
	want := "(() . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	_, ok := <-a
	if !ok {
		t.Fatalf("expected never ending")
	}
}

func TestRunGoalAlways3(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := RunGoal(ctx, 3, AlwaysO)
	if len(ss) != 3 {
		t.Fatalf("expected 3 got %d", len(ss))
	}
	sss := fmap(func(s *State) string {
		return s.String()
	}, ss)
	got := "(" + strings.Join(sss, " ") + ")"
	want := "((() . 0) (() . 0) (() . 0))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestRunGoalDisj2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e1 := EqualO(
		ast.NewSymbol("olive"),
		ast.NewVariable("x"),
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		ast.NewVariable("x"),
	)
	g := Disj(e1, e2)
	ss := RunGoal(ctx, 5, g)
	if len(ss) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2NoResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	x := ast.NewVariable("x")
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		x,
	)
	g := Conj(e1, e2)
	ss := RunGoal(ctx, 5, g)
	if len(ss) != 0 {
		t.Fatalf("expected 0, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2OneResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	x := ast.NewVariable("x")
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	g := Conj(e1, e2)
	ss := RunGoal(ctx, 5, g)
	if len(ss) != 1 {
		t.Fatalf("expected 1, got %d: %v", len(ss), ss)
	}
	got := ss[0].String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", indexOf(x)), "x", -1)
	want := "((,x . olive) . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
