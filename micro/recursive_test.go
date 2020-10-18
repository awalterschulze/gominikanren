package micro

import (
    "testing"

    "github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(defrel (very-recursiveo)
  (conde
    ((nevero))
    ((very-recursiveo))
    ((alwayso))
    ((very-recursiveo))
    ((nevero))))
*/
func veryRecursiveO() Goal {
	return Conde(
		[]Goal{NeverO()},
		[]Goal{Zzz2(veryRecursiveO)},
		[]Goal{AlwaysO()},
		[]Goal{Zzz2(veryRecursiveO)},
		[]Goal{NeverO()},
	)
}

func TestRecursive(t *testing.T) {
	ss := RunGoal(
		5,
		CallFresh(func(q *ast.SExpr) Goal {
			return veryRecursiveO()
		}),
	)
    if len(ss) != 5 {
        t.Fatalf("expected %d, but got %d results", 5, len(ss))
    }
    sexprs := Reify("q", ss)
    for _, sexpr := range sexprs {
        if sexpr.String() != "_0" {
            t.Fatalf("expected %s, but got %s instead", "_0", sexpr.String())
        }
    }
}
