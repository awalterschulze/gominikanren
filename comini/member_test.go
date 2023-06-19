package comini

import (
	"context"
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestMemberO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := comicro.NewEmptyState().WithReifyNames(ast.ReifyName)
	list := ast.Cons(ast.NewInt(0), ast.Cons(ast.NewInt(1), ast.Cons(ast.NewInt(2), nil)))
	sanys := comicro.Run(ctx, -1, s, func(q *ast.SExpr) comicro.Goal {
		return MemberO(
			q,
			list,
		)
	})
	sexprs := ast.Sort(fmap(sanys, func(a any) *ast.SExpr {
		return a.(*ast.SExpr)
	}))
	if len(sexprs) != 3 {
		t.Fatalf("expected len %d, but got len %d instead", 3, len(sexprs))
	}
	for i, sexpr := range sexprs {
		if *sexpr.Atom.Int != int64(i) {
			t.Fatalf("expected %d, but got %d instead", i, *sexpr.Atom.Int)
		}
	}
}

func TestMapO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := comicro.NewEmptyState().WithReifyNames(ast.ReifyName)
	list := ast.Cons(ast.NewInt(0), ast.Cons(ast.NewInt(1), ast.Cons(ast.NewInt(2), nil)))
	sexprs := comicro.Run(ctx, -1, s, func(q *ast.SExpr) comicro.Goal {
		return MapO(
			func(x, y *ast.SExpr) comicro.Goal { return comicro.EqualO(x, y) },
			q,
			list,
		)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	if !reflect.DeepEqual(sexprs[0], list) {
		t.Fatalf("expected output matches input list, but got %#v", sexprs[0])
	}
}
