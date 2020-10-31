package concurrent

import (
    "reflect"
    "testing"

	"github.com/awalterschulze/gominikanren/mini"
)

func TestConcurrentDisjPlus(t *testing.T) {
    disjunction = mini.DisjPlus
    goal := einstein()
    sexprs := runEinstein(goal)
    disjunction = DisjPlus
    goal = einstein()
    concsexprs := runEinstein(goal)
    if !reflect.DeepEqual(sexprs, concsexprs) {
        t.Fatalf("expected equal outcomes but got %#v %#v", sexprs, concsexprs)
    }
}
