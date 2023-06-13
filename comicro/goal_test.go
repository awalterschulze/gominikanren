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
			ss := NewStreamForGoal(ctx, EqualO(uexpr, vexpr), EmptyState())
			got := ss.String()
			if got != want {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestFailureO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, FailureO, EmptyState())
	if got, want := ss.String(), "()"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestSuccessO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, SuccessO, EmptyState())
	if got, want := ss.String(), "((() . 0))"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNeverO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, NeverO, EmptyState())
	s, ok := <-ss
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
	ss := NewStreamForGoal(ctx,
		Disj(
			EqualO(
				ast.NewSymbol("olive"),
				x,
			),
			NeverO,
		),
		EmptyState(),
	)
	var s *State
	for s = range ss {
		if s != nil {
			break
		}
	}
	got := s.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", indexOf(x)), "x", -1)
	want := "({,x: olive} . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	_, ok := <-ss
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
	ss := NewStreamForGoal(ctx,
		Disj(
			NeverO,
			EqualO(
				ast.NewSymbol("olive"),
				x,
			),
		),
		EmptyState(),
	)
	for s := range ss {
		if s != nil {
			got := s.String()
			// reifying y; we assigned it a random uint64 and lost track of it
			got = strings.Replace(got, fmt.Sprintf("v%d", indexOf(x)), "x", -1)
			want := "({,x: olive} . 0)"
			if got != want {
				t.Fatalf("got %s != want %s", got, want)
			}
			break
		}
	}
	_, ok := <-ss
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
	ss := NewStreamForGoal(ctx, AlwaysO, EmptyState())
	var s *State
	for s = range ss {
		if s != nil {
			break
		}
	}
	got := s.String()
	want := "(() . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	_, ok := <-ss
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
	sss := fmap(ss, func(s *State) string {
		return s.String()
	})
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
	want := "({,x: olive} . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
