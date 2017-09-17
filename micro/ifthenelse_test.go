package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestIfThenElseSuccess(t *testing.T) {
	ifte := IfThenElseO(
		SuccessO(),
		EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseFailure(t *testing.T) {
	ifte := IfThenElseO(
		FailureO(),
		EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #t)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseXIsTrue(t *testing.T) {
	ifte := IfThenElseO(
		EqualO(ast.NewSymbol("#t"), ast.NewVariable("x")),
		EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f) (,x . #t)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseDisjoint(t *testing.T) {
	ifte := IfThenElseO(
		DisjointO(
			EqualO(ast.NewSymbol("#t"), ast.NewVariable("x")),
			EqualO(ast.NewSymbol("#f"), ast.NewVariable("x")),
		),
		EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f) (,x . #t)) `((,y . #f) (,x . #f)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
