package regex

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func testo(t *testing.T, f func(q *ast.SExpr) comicro.Goal, want *ast.SExpr) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sexprs := comicro.Run(ctx, 1, f)
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d result, but got %d instead", 1, len(sexprs))
	}
	got := sexprs[0]
	if !got.Equal(want) {
		t.Fatalf("expected %s, but got %s instead", want, got)
	}
}
