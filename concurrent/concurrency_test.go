package concurrent

import (
    "reflect"
    "testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
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

func memberOfLargeList(n int) func(*ast.SExpr) micro.Goal {
    ints := make([]*ast.SExpr, n)
    for i := 0; i < n; i++ {
        ints[i] = ast.NewInt(int64(i))
    }
    membero := memberOUnrolled(ast.NewList(ints...))
    return func(q *ast.SExpr) micro.Goal {
        return membero(q)
    }
}

func TestHeavilyConcurrentDisj(t *testing.T) {
    disjunction = DisjPlus
    n := 10000
    membero := memberOfLargeList(n)
    sexprs := micro.Run(-1, membero)
    if len(sexprs) != n {
        t.Fatalf("expected %d results but got %d", n, len(sexprs))
    }
}

func BenchmarkMembero10000(b *testing.B) {
	disjunction = disjPlus
    membero := memberOfLargeList(10000)
	var r []*ast.SExpr
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkConcMembero10000(b *testing.B) {
	disjunction = DisjPlus
    membero := memberOfLargeList(10000)
	var r []*ast.SExpr
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMembero50000(b *testing.B) {
	disjunction = disjPlus
    membero := memberOfLargeList(50000)
	var r []*ast.SExpr
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkConcMembero50000(b *testing.B) {
	disjunction = DisjPlus
    membero := memberOfLargeList(50000)
	var r []*ast.SExpr
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}
