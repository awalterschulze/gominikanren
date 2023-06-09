package comini

import (
	"context"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOnce(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// assigning var index by hand since we want x < y
	// otherwise test can fail because of int order in substitution map
	x := ast.NewVar("x", 10001)
	y := ast.NewVar("y", 10002)
	ifte := IfThenElseO(
		OnceO(comicro.Disj(
			comicro.EqualO(ast.NewSymbol("#t"), x),
			comicro.EqualO(ast.NewSymbol("#f"), x),
		)),
		comicro.EqualO(ast.NewSymbol("#f"), y),
		comicro.EqualO(ast.NewSymbol("#t"), y),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, comicro.EmptyState())
	got := ss.String()
	// reifying x and y; we assigned them a random uint64 and lost track of it
	got = strings.Replace(got, "v10001", "x", -1)
	got = strings.Replace(got, "v10002", "y", -1)
	want1 := "(((,y . #f) (,x . #t) . 0))"
	want2 := "(((,y . #f) (,x . #f) . 0))"
	if got != want1 && got != want2 {
		t.Fatalf("got %v != want %v or %v", got, want1, want2)
	}
}
