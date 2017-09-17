package micro

import (
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestEqualO(t *testing.T) {
	tests := []func() (string, string, string){
		deriveTuple3("#f", "#f", "(`())"),
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
			stream := EqualO(uexpr, vexpr)(EmptySubstitution())
			got := stream.String()
			if got != want {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestFailureO(t *testing.T) {
	if got, want := FailureO()(EmptySubstitution()).String(), "()"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestSuccessO(t *testing.T) {
	if got, want := SuccessO()(EmptySubstitution()).String(), "(`())"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestNeverO(t *testing.T) {
	n := NeverO()(EmptySubstitution())
	s, sok := n()
	if s != nil {
		t.Fatalf("expected suspension")
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestDisjointO1(t *testing.T) {
	d := DisjointO(
		EqualO(
			ast.NewSymbol("olive"),
			ast.NewVariable("x"),
		),
		NeverO(),
	)(EmptySubstitution())
	s, sok := d()
	got := s.String()
	want := "`((,x . olive))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestDisjointO2(t *testing.T) {
	d := DisjointO(
		NeverO(),
		EqualO(
			ast.NewSymbol("olive"),
			ast.NewVariable("x"),
		),
	)(EmptySubstitution())
	s, sok := d()
	if s != nil {
		t.Fatalf("expected suspension")
	}
	if sok == nil {
		t.Fatalf("expected more")
	}
	s, sok = sok()
	got := s.String()
	want := "`((,x . olive))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestAlwaysO(t *testing.T) {
	a := AlwaysO()(EmptySubstitution())
	s, sok := a()
	if s != nil {
		t.Fatal("expected suspension")
	}
	s, sok = sok()
	got := s.String()
	want := "`()"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	if sok == nil {
		t.Fatalf("expected never ending")
	}
}

func TestRunGoalAlways3(t *testing.T) {
	ss := RunGoal(3, AlwaysO())
	if len(ss) != 3 {
		t.Fatalf("expected 3 got %d", len(ss))
	}
	sss := deriveFmaps(func(s Substitution) string {
		return s.String()
	}, ss)
	got := "(" + strings.Join(sss, " ") + ")"
	want := "(`() `() `())"
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
	g := DisjointO(e1, e2)
	ss := RunGoal(5, g)
	if len(ss) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2NoResults(t *testing.T) {
	e1 := EqualO(
		ast.NewSymbol("olive"),
		ast.NewVariable("x"),
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		ast.NewVariable("x"),
	)
	g := ConjunctionO(e1, e2)
	ss := RunGoal(5, g)
	if len(ss) != 0 {
		t.Fatalf("expected 0, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2OneResults(t *testing.T) {
	e1 := EqualO(
		ast.NewSymbol("olive"),
		ast.NewVariable("x"),
	)
	e2 := EqualO(
		ast.NewSymbol("olive"),
		ast.NewVariable("x"),
	)
	g := ConjunctionO(e1, e2)
	ss := RunGoal(5, g)
	if len(ss) != 1 {
		t.Fatalf("expected 1, got %d: %v", len(ss), ss)
	}
	got := ss[0].String()
	want := "`((,x . olive))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
