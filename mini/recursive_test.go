package mini

import (
    "testing"

    "github.com/awalterschulze/gominikanren/micro"
    "github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestRecursive(t *testing.T) {
	ss := micro.RunGoal(
		5,
		micro.CallFresh(func(q *ast.SExpr) micro.Goal {
			return VeryRecursiveO()
		}),
	)
    if len(ss) != 5 {
        t.Fatalf("expected %d, but got %d results", 5, len(ss))
    }
    sexprs := micro.Reify("q", ss)
    for _, sexpr := range sexprs {
        if sexpr.String() != "_0" {
            t.Fatalf("expected %s, but got %s instead", "_0", sexpr.String())
        }
    }
}
