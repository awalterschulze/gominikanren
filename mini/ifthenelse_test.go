package mini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestIfThenElseSuccess(t *testing.T) {
	ifte := IfThenElseO(
		micro.SuccessO,
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	want := "(((,y . #f) . 0))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseFailure(t *testing.T) {
	ifte := IfThenElseO(
		micro.FailureO,
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	want := "(((,y . #t) . 0))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseXIsTrue(t *testing.T) {
	ifte := IfThenElseO(
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("x")),
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	want := "(((,y . #f) (,x . #t) . 0))"
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
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	want := "(((,y . #f) (,x . #t) . 0) ((,y . #f) (,x . #f) . 0))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
