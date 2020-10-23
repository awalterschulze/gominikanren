package mini

import (
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestMemberO(t *testing.T) {
	list := ast.Cons(ast.NewInt(0), ast.Cons(ast.NewInt(1), ast.Cons(ast.NewInt(2), nil)))
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return MemberO(
			q,
			list,
		)
	})
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
	list := ast.Cons(ast.NewInt(0), ast.Cons(ast.NewInt(1), ast.Cons(ast.NewInt(2), nil)))
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return MapO(
			micro.EqualO,
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
