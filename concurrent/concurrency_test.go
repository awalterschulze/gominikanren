package concurrent

import (
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

var result []*ast.SExpr

func TestConcurrentDisjPlus(t *testing.T) {
	goal := einstein(mini.DisjPlus)
	sexprs := runEinstein(goal)
	goal = einstein(DisjPlus)
	concsexprs := runEinstein(goal)
	if !reflect.DeepEqual(sexprs, concsexprs) {
		t.Fatalf("expected equal outcomes but got %#v %#v", sexprs, concsexprs)
	}
}

func TestConcurrentConjPlus(t *testing.T) {
	goal := equalLargeList(10, mini.ConjPlus)
	sexprs := micro.Run(-1, goal)
	goal = equalLargeList(10, ConjPlus)
	concsexprs := micro.Run(-1, goal)
	if !reflect.DeepEqual(sexprs, concsexprs) {
		t.Fatalf("expected equal outcomes but got %#v %#v", sexprs, concsexprs)
	}
}

func largeList(n int) *ast.SExpr {
	ints := make([]*ast.SExpr, n)
	for i := 0; i < n; i++ {
		ints[i] = ast.NewInt(int64(i))
	}
	return ast.NewList(ints...)
}

func memberOfLargeList(f func(...micro.Goal) micro.Goal, n int) func(*ast.SExpr) micro.Goal {
	membero := memberOUnrolled(f, largeList(n))
	return func(q *ast.SExpr) micro.Goal {
		return membero(q)
	}
}

func TestHeavilyConcurrentDisj(t *testing.T) {
	if testing.Short() {
		return
	}
	n := 10000
	membero := memberOfLargeList(DisjPlus, n)
	sexprs := micro.Run(-1, membero)
	if len(sexprs) != n {
		t.Fatalf("expected %d results but got %d", n, len(sexprs))
	}
}

// before calling bench, remember to set disjunction
func benchMembero(b *testing.B, listsize int, f func(...micro.Goal) micro.Goal) {
	membero := memberOfLargeList(f, listsize)
	var r []*ast.SExpr
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMembero10000(b *testing.B) {
	benchMembero(b, 10000, mini.DisjPlusNoZzz)
}

func BenchmarkConcMembero10000(b *testing.B) {
	benchMembero(b, 10000, DisjPlus)
}

func BenchmarkMembero50000(b *testing.B) {
	benchMembero(b, 50000, mini.DisjPlusNoZzz)
}

func BenchmarkConcMembero50000(b *testing.B) {
	benchMembero(b, 50000, DisjPlus)
}

func mapODoubleUnrolled(f func(*ast.SExpr, *ast.SExpr) micro.Goal, list *ast.SExpr, conj func(...micro.Goal) micro.Goal) func(*ast.SExpr) micro.Goal {
	goals := []func(*ast.SExpr) micro.Goal{}
	for {
		if list == nil {
			break
		}
		car, cdr := list.Car(), list.Cdr()
		goals = append(goals, func(x *ast.SExpr) micro.Goal {
			return f(x, car)
		})
		list = cdr
	}
	n := len(goals)
	return func(x *ast.SExpr) micro.Goal {
		gs := make([]micro.Goal, n+1)
		vars := make([]*ast.SExpr, n)
		var cons *ast.SExpr = nil
		for i := 0; i < n; i++ {
			v := ast.NewVariable("")
			vars[i] = v
			gs[n-i] = goals[n-i-1](v)
			cons = ast.Cons(v, cons)
		}
		gs[0] = micro.EqualO(x, cons)
		return conj(gs...)
	}
}

func equalLargeList(n int, f func(...micro.Goal) micro.Goal) func(*ast.SExpr) micro.Goal {
	mapo := mapODoubleUnrolled(micro.EqualO, largeList(n), f)
	return func(q *ast.SExpr) micro.Goal {
		return mapo(q)
	}
}

// before calling bench, remember to set conjunction
func benchMapo(b *testing.B, listsize int, f func(...micro.Goal) micro.Goal) {
	membero := equalLargeList(listsize, f)
	var r []*ast.SExpr
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMapo10000(b *testing.B) {
	benchMapo(b, 10000, mini.ConjPlusNoZzz)
}

func BenchmarkConcMapo10000(b *testing.B) {
	benchMapo(b, 10000, ConjPlus)
}

func largeListHeadFail(n int) *ast.SExpr {
	ints := make([]*ast.SExpr, n)
	for i := 0; i < n-1; i++ {
		ints[i] = ast.NewInt(int64(i))
	}
	ints[0] = ast.NewInt(1)
	return ast.NewList(ints...)
}

func equalFailHeadLargeList(n int, f func(...micro.Goal) micro.Goal) func(*ast.SExpr) micro.Goal {
	list := largeList(n)
	mapo := mapODoubleUnrolled(micro.EqualO, list, f)
	list2 := largeListHeadFail(n)
	return func(*ast.SExpr) micro.Goal { return mapo(list2) }
}

// before calling bench, remember to set conjunction
func benchMapoFailHead(b *testing.B, listsize int, f func(...micro.Goal) micro.Goal) {
	membero := equalFailHeadLargeList(listsize, f)
	var r []*ast.SExpr
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMapoFailHead10000(b *testing.B) {
	benchMapoFailHead(b, 10000, mini.ConjPlusNoZzz)
}

func BenchmarkConcMapoFailHead10000(b *testing.B) {
	benchMapoFailHead(b, 10000, ConjPlus)
}

func largeListTailFail(n int) *ast.SExpr {
	ints := make([]*ast.SExpr, n)
	for i := 0; i < n-1; i++ {
		ints[i] = ast.NewInt(int64(i))
	}
	ints[n-1] = ast.NewInt(0)
	return ast.NewList(ints...)
}

func equalFailTailLargeList(n int, f func(...micro.Goal) micro.Goal) func(*ast.SExpr) micro.Goal {
	list := largeList(n)
	mapo := mapODoubleUnrolled(micro.EqualO, list, f)
	list2 := largeListTailFail(n)
	return func(*ast.SExpr) micro.Goal { return mapo(list2) }
}

// before calling bench, remember to set conjunction
func benchMapoFailTail(b *testing.B, listsize int, f func(...micro.Goal) micro.Goal) {
	membero := equalFailTailLargeList(listsize, f)
	var r []*ast.SExpr
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMapoFailTail10000(b *testing.B) {
	benchMapoFailTail(b, 10000, mini.ConjPlusNoZzz)
}

func BenchmarkConcMapoFailTail10000(b *testing.B) {
	benchMapoFailTail(b, 10000, ConjPlus)
}
