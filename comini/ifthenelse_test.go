package comini

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// helper func used in all tests that use Substitution of vars
func indexOf(x *ast.SExpr) uint64 {
	return x.Atom.Var.Index
}

func TestIfThenElseSuccess(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var y *ast.SExpr
	state := comicro.NewEmptyState()
	state, y = comicro.NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		comicro.SuccessO,
		comicro.EqualO(ast.NewSymbol("#f"), y),
		comicro.EqualO(ast.NewSymbol("#t"), y),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want := "(({,y: #f} . 1))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseFailure(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var y *ast.SExpr
	state := comicro.NewEmptyState()
	state, y = comicro.NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		comicro.FailureO,
		comicro.EqualO(ast.NewSymbol("#f"), y),
		comicro.EqualO(ast.NewSymbol("#t"), y),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want := "(({,y: #t} . 1))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseXIsTrue(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var x, y *ast.SExpr
	state := comicro.NewEmptyState()
	state, x = comicro.NewVarWithName(state, "x", &ast.SExpr{})
	state, y = comicro.NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		comicro.EqualO(ast.NewSymbol("#t"), x),
		comicro.EqualO(ast.NewSymbol("#f"), y),
		comicro.EqualO(ast.NewSymbol("#t"), y),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want := "(({,x: #t}, {,y: #f} . 2))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestIfThenElseDisjoint(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var x, y *ast.SExpr
	state := comicro.NewEmptyState()
	state, x = comicro.NewVarWithName(state, "x", &ast.SExpr{})
	state, y = comicro.NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		comicro.Disj(
			comicro.EqualO(ast.NewSymbol("#t"), x),
			comicro.EqualO(ast.NewSymbol("#f"), x),
		),
		comicro.EqualO(ast.NewSymbol("#f"), y),
		comicro.EqualO(ast.NewSymbol("#t"), y),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want := "(({,x: #f}, {,y: #f} . 2) ({,x: #t}, {,y: #f} . 2))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
