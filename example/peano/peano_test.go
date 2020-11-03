package peano

import (
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestNatPlus(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return Natplus(Makenat(3), Makenat(6), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := Parsenat(sexprs[0]), 9
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
}

func TestLeq(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return Leq(q, Makenat(6))
	})
	if len(sexprs) != 7 {
		t.Fatalf("expected len %d, but got len %d instead", 7, len(sexprs))
	}
	// answer we're looking for is the list of 0..6 incl
	for i, sexpr := range sexprs {
		got, want := Parsenat(sexpr), i
		if got != want {
			t.Fatalf("expected %d, but got %d instead", want, got)
		}
	}
	// we expect q to be any successor of 4
	sexprs = micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return Leq(Makenat(4), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	// representing succ as wrapping with a list, we expect to find a
	// var (reified to symbol _0) after taking car of answer 4 times
	got := sexprs[0]
	for i := 0; i < 4; i++ {
		got = got.Car()
	}
	want := ast.NewSymbol("_0")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, but got %v instead", want, got)
	}
}

func TestHalf(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return Half(Makenat(7), q)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	got, want := Parsenat(sexprs[0]), 3
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
	sexprs = micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return Half(q, Makenat(5))
	})
	if len(sexprs) != 2 {
		t.Fatalf("expected len %d, but got len %d instead", 2, len(sexprs))
	}
	got, want = Parsenat(sexprs[0]), 10
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
	got, want = Parsenat(sexprs[1]), 11
	if got != want {
		t.Fatalf("expected %d, but got %d instead", want, got)
	}
}
