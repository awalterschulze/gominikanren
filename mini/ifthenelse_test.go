package mini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestIfThenElseSuccess(t *testing.T) {
	ifte := IfThenElseO(
		micro.SuccessO(),
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(micro.EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseFailure(t *testing.T) {
	ifte := IfThenElseO(
		micro.FailureO(),
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(micro.EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #t)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseXIsTrue(t *testing.T) {
	ifte := IfThenElseO(
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("x")),
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(micro.EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f) (,x . #t)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseDisjoint(t *testing.T) {
	ifte := IfThenElseO(
		micro.DisjointO(
			micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("x")),
			micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("x")),
		),
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(micro.EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f) (,x . #t)) `((,y . #f) (,x . #f)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
