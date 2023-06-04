package mini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
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
func veryRecursiveO(s *micro.State) *micro.StreamOfStates {
	return Conde(
		[]micro.Goal{micro.NeverO},
		[]micro.Goal{veryRecursiveO},
		[]micro.Goal{micro.AlwaysO},
		[]micro.Goal{veryRecursiveO},
		[]micro.Goal{micro.NeverO},
	)(s)
}

func TestRecursive(t *testing.T) {
	num := 5
	ss := micro.RunGoal(
		num,
		micro.Zzz(veryRecursiveO),
	)
	if len(ss) != num {
		t.Fatalf("expected %d, but got %d results", 5, len(ss))
	}
	sexprs := micro.Reify("q", ss)
	for _, sexpr := range sexprs {
		if sexpr.String() != "_0" {
			t.Fatalf("expected %s, but got %s instead", "_0", sexpr.String())
		}
	}
}
