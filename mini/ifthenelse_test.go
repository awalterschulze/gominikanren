package mini

import (
	"fmt"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestIfThenElseSuccess(t *testing.T) {
	y := ast.NewVariable("y")
	ifte := IfThenElseO(
		micro.SuccessO,
		micro.EqualO(ast.NewSymbol("#f"), y),
		micro.EqualO(ast.NewSymbol("#t"), y),
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", y.Atom.Var.Index), "y", -1)
	want := "(((,y . #f) . 0))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseFailure(t *testing.T) {
	y := ast.NewVariable("y")
	ifte := IfThenElseO(
		micro.FailureO,
		micro.EqualO(ast.NewSymbol("#f"), y),
		micro.EqualO(ast.NewSymbol("#t"), y),
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	got = strings.Replace(got, fmt.Sprintf("v%d", y.Atom.Var.Index), "y", -1)
	want := "(((,y . #t) . 0))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseXIsTrue(t *testing.T) {
	// assigning var index by hand since we want x < y
	// otherwise test can fail because of int order in substitution map
	x := ast.NewVar("x", 10001)
	y := ast.NewVar("y", 10002)
	ifte := IfThenElseO(
		micro.EqualO(ast.NewSymbol("#t"), x),
		micro.EqualO(ast.NewSymbol("#f"), y),
		micro.EqualO(ast.NewSymbol("#t"), y),
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	// reifying x and y; we assigned them a random uint64 and lost track of it
	got = strings.Replace(got, "v10001", "x", -1)
	got = strings.Replace(got, "v10002", "y", -1)
	want := "(((,y . #f) (,x . #t) . 0))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseDisjoint(t *testing.T) {
	// assigning var index by hand since we want x < y
	// otherwise test can fail because of int order in substitution map
	x := ast.NewVar("x", 10001)
	y := ast.NewVar("y", 10002)
	ifte := IfThenElseO(
		micro.DisjointO(
			micro.EqualO(ast.NewSymbol("#t"), x),
			micro.EqualO(ast.NewSymbol("#f"), x),
		),
		micro.EqualO(ast.NewSymbol("#f"), y),
		micro.EqualO(ast.NewSymbol("#t"), y),
	)()
	ss := ifte(micro.EmptyState())
	got := ss.String()
	// reifying x and y; we assigned them a random uint64 and lost track of it
	got = strings.Replace(got, "v10001", "x", -1)
	got = strings.Replace(got, "v10002", "y", -1)
	want := "(((,y . #f) (,x . #t) . 0) ((,y . #f) (,x . #f) . 0))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
