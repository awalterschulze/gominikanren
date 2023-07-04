package gomini

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestIfThenElseSuccess(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var y *ast.SExpr
	state := NewState()
	state, y = NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		SuccessO,
		EqualO(ast.NewSymbol("#f"), y),
		EqualO(ast.NewSymbol("#t"), y),
	)
	ss := NewStreamForGoal(ctx, ifte, state)
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
	state := NewState()
	state, y = NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		FailureO,
		EqualO(ast.NewSymbol("#f"), y),
		EqualO(ast.NewSymbol("#t"), y),
	)
	ss := NewStreamForGoal(ctx, ifte, state)
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
	state := NewState()
	state, x = NewVarWithName(state, "x", &ast.SExpr{})
	state, y = NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		EqualO(ast.NewSymbol("#t"), x),
		EqualO(ast.NewSymbol("#f"), y),
		EqualO(ast.NewSymbol("#t"), y),
	)
	ss := NewStreamForGoal(ctx, ifte, state)
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
	state := NewState()
	state, x = NewVarWithName(state, "x", &ast.SExpr{})
	state, y = NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		Disj(
			EqualO(ast.NewSymbol("#t"), x),
			EqualO(ast.NewSymbol("#f"), x),
		),
		EqualO(ast.NewSymbol("#f"), y),
		EqualO(ast.NewSymbol("#t"), y),
	)
	ss := NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want := "(({,x: #f}, {,y: #f} . 2) ({,x: #t}, {,y: #f} . 2))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
