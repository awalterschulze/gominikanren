package micro

import (
	"fmt"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestEqualO(t *testing.T) {
	tests := []func() (string, string, string){
		deriveTuple3("#f", "#f", "((() . 0))"),
		deriveTuple3("#f", "#t", "()"),
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
			stream := EqualO(uexpr, vexpr)(EmptyState())
			got := stream.String()
			if got != want {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestFailureO(t *testing.T) {
	if got, want := FailureO(EmptyState()).String(), "()"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestSuccessO(t *testing.T) {
	if got, want := SuccessO(EmptyState()).String(), "((() . 0))"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNeverO(t *testing.T) {
	n := NeverO(EmptyState())
	s, sok := n.CarCdr()
	if s != nil {
		t.Fatalf("expected suspension")
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestDisj1(t *testing.T) {
	x := ast.NewVariable("x")
	d := Disj(
		EqualO(
			ast.NewSymbol("olive"),
			x,
		),
		NeverO,
	)(EmptyState())
	s, sok := d.CarCdr()
	got := s.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", indexOf(x)), "x", -1)
	want := "((,x . olive) . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestDisj2(t *testing.T) {
	x := ast.NewVariable("x")
	d := Disj(
		NeverO,
		EqualO(
			ast.NewSymbol("olive"),
			x,
		),
	)(EmptyState())
	s, sok := d.CarCdr()
	if s != nil {
		t.Fatalf("expected suspension")
	}
	if sok == nil {
		t.Fatalf("expected more")
	}
	s, sok = sok.CarCdr()
	got := s.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", indexOf(x)), "x", -1)
	want := "((,x . olive) . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestAlwaysO(t *testing.T) {
	a := AlwaysO(EmptyState())
	s, sok := a.CarCdr()
	if s != nil {
		t.Fatal("expected suspension")
	}
	s, sok = sok.CarCdr()
	got := s.String()
	want := "(() . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestRunGoalAlways3(t *testing.T) {
	ss := RunGoal(3, AlwaysO)
	if len(ss) != 3 {
		t.Fatalf("expected 3 got %d", len(ss))
	}
	sss := deriveFmaps(func(s *State) string {
		return s.String()
	}, ss)
	got := "(" + strings.Join(sss, " ") + ")"
	want := "((() . 0) (() . 0) (() . 0))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestRunGoalDisj2(t *testing.T) {
	e1 := EqualO(
		ast.NewSymbol("olive"),
		ast.NewVariable("x"),
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		ast.NewVariable("x"),
	)
	g := Disj(e1, e2)
	ss := RunGoal(5, g)
	if len(ss) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2NoResults(t *testing.T) {
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
	ss := RunGoal(5, g)
	if len(ss) != 0 {
		t.Fatalf("expected 0, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2OneResults(t *testing.T) {
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
	ss := RunGoal(5, g)
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
