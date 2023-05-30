package comini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
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
func veryRecursiveO(s *comicro.State) comicro.StreamOfStates {
	return Conde(
		[]comicro.Goal{comicro.NeverO},
		[]comicro.Goal{veryRecursiveO},
		[]comicro.Goal{comicro.AlwaysO},
		[]comicro.Goal{veryRecursiveO},
		[]comicro.Goal{comicro.NeverO},
	)(s)
}

func TestRecursive(t *testing.T) {
	ss := comicro.RunGoal(
		5,
		comicro.Zzz(veryRecursiveO),
	)
	if len(ss) != 5 {
		t.Fatalf("expected %d, but got %d results", 5, len(ss))
	}
	sexprs := comicro.Reify("q", ss)
	for _, sexpr := range sexprs {
		if sexpr.String() != "_0" {
			t.Fatalf("expected %s, but got %s instead", "_0", sexpr.String())
		}
	}
}
