package concurrent

import (
	"reflect"
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/mini"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

var (
	result      []*ast.SExpr
	disjunction func(...micro.Goal) micro.Goal
	conjunction func(...micro.Goal) micro.Goal
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

func TestConcurrentConjPlus(t *testing.T) {
	conjunction = mini.ConjPlus
	goal := equalLargeList(10)
	sexprs := micro.Run(-1, goal)
	conjunction = ConjPlus
	goal = equalLargeList(10)
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

func memberOfLargeList(n int) func(*ast.SExpr) micro.Goal {
	membero := memberOUnrolled(largeList(n))
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

// before calling bench, remember to set disjunction
func benchMembero(b *testing.B, listsize int) {
	membero := memberOfLargeList(listsize)
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMembero10000(b *testing.B) {
	disjunction = disjPlus
	benchMembero(b, 10000)
}

func BenchmarkConcMembero10000(b *testing.B) {
	disjunction = DisjPlus
	benchMembero(b, 10000)
}

func BenchmarkMembero50000(b *testing.B) {
	disjunction = disjPlus
	benchMembero(b, 50000)
}

func BenchmarkConcMembero50000(b *testing.B) {
	disjunction = DisjPlus
	benchMembero(b, 50000)
}

func mapODoubleUnrolled(f func(*ast.SExpr, *ast.SExpr) micro.Goal, list *ast.SExpr) func(*ast.SExpr) micro.Goal {
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
		return conjunction(gs...)
	}
}

func equalLargeList(n int) func(*ast.SExpr) micro.Goal {
	mapo := mapODoubleUnrolled(micro.EqualO, largeList(n))
	return func(q *ast.SExpr) micro.Goal {
		return mapo(q)
	}
}

// before calling bench, remember to set conjunction
func benchMapo(b *testing.B, listsize int) {
	membero := equalLargeList(listsize)
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func conjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.SuccessO
	}
	if len(gs) == 1 {
		return gs[0]
	}
	g1 := gs[0]
	g2 := conjPlus(gs[1:]...)
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			g1s := g1()(s)
			return micro.Bind(g1s, g2)
		}
	}
}

func BenchmarkMapo10000(b *testing.B) {
	conjunction = conjPlus
	benchMapo(b, 10000)
}

func BenchmarkConcMapo10000(b *testing.B) {
	conjunction = ConjPlus
	benchMapo(b, 10000)
}

func largeListHeadFail(n int) *ast.SExpr {
	ints := make([]*ast.SExpr, n)
	for i := 0; i < n-1; i++ {
		ints[i] = ast.NewInt(int64(i))
	}
	ints[0] = ast.NewInt(1)
	return ast.NewList(ints...)
}

func equalFailHeadLargeList(n int) func(*ast.SExpr) micro.Goal {
	list := largeList(n)
	mapo := mapODoubleUnrolled(micro.EqualO, list)
	list2 := largeListHeadFail(n)
	return func(*ast.SExpr) micro.Goal { return mapo(list2) }
}

// before calling bench, remember to set conjunction
func benchMapoFailHead(b *testing.B, listsize int) {
	membero := equalFailHeadLargeList(listsize)
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMapoFailHead10000(b *testing.B) {
	conjunction = conjPlus
	benchMapoFailHead(b, 10000)
}

func BenchmarkConcMapoFailHead10000(b *testing.B) {
	conjunction = ConjPlus
	benchMapoFailHead(b, 10000)
}

func largeListTailFail(n int) *ast.SExpr {
	ints := make([]*ast.SExpr, n)
	for i := 0; i < n-1; i++ {
		ints[i] = ast.NewInt(int64(i))
	}
	ints[n-1] = ast.NewInt(0)
	return ast.NewList(ints...)
}

func equalFailTailLargeList(n int) func(*ast.SExpr) micro.Goal {
	list := largeList(n)
	mapo := mapODoubleUnrolled(micro.EqualO, list)
	list2 := largeListTailFail(n)
	return func(*ast.SExpr) micro.Goal { return mapo(list2) }
}

// before calling bench, remember to set conjunction
func benchMapoFailTail(b *testing.B, listsize int) {
	membero := equalFailTailLargeList(listsize)
	var r []*ast.SExpr
    b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = micro.Run(-1, membero)
	}
	result = r
}

func BenchmarkMapoFailTail10000(b *testing.B) {
	conjunction = conjPlus
	benchMapoFailTail(b, 10000)
}

func BenchmarkConcMapoFailTail10000(b *testing.B) {
	conjunction = ConjPlus
	benchMapoFailTail(b, 10000)
}
