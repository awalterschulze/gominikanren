package gomini

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// OnceO is a goal that returns one successful state.
func OnceO(g Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		gs := NewStreamForGoal(ctx, g, s)
		state, ok := gs.Read(ctx)
		if !ok {
			return
		}
		ss.Write(ctx, state)
	}
}

func TestOnce(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var x, y *ast.SExpr
	state := NewState()
	state, x = newVarWithName(state, "x", &ast.SExpr{})
	state, y = newVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		OnceO(Disj(
			EqualO(ast.NewSymbol("#t"), x),
			EqualO(ast.NewSymbol("#f"), x),
		)),
		EqualO(ast.NewSymbol("#f"), y),
		EqualO(ast.NewSymbol("#t"), y),
	)
	ss := NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want1 := "(({,x: #t}, {,y: #f} . 2))"
	want2 := "(({,x: #f}, {,y: #f} . 2))"
	if got != want1 && got != want2 {
		t.Fatalf("got %v != want %v or %v", got, want1, want2)
	}
}
