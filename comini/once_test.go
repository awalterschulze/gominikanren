package comini

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOnce(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var x, y *ast.SExpr
	state := comicro.NewEmptyState()
	state, x = comicro.NewVarWithName(state, "x", &ast.SExpr{})
	state, y = comicro.NewVarWithName(state, "y", &ast.SExpr{})
	ifte := IfThenElseO(
		OnceO(comicro.Disj(
			comicro.EqualO(ast.NewSymbol("#t"), x),
			comicro.EqualO(ast.NewSymbol("#f"), x),
		)),
		comicro.EqualO(ast.NewSymbol("#f"), y),
		comicro.EqualO(ast.NewSymbol("#t"), y),
	)
	ss := comicro.NewStreamForGoal(ctx, ifte, state)
	got := ss.String()
	want1 := "(({,x: #t}, {,y: #f} . 2))"
	want2 := "(({,x: #f}, {,y: #f} . 2))"
	if got != want1 && got != want2 {
		t.Fatalf("got %v != want %v or %v", got, want1, want2)
	}
}
