package mini

import (
	"fmt"
	"testing"
	"time"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// remember, concurrency != parallelism

// setting up exclusive-or for a very simple or-parallelism demo
func disjExcl(g1, g2 micro.Goal) micro.Goal {
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			g1s := g1()(s)
			if g1s != nil {
				return g1s
			}
			g2s := g2()(s)
			return micro.Mplus(g1s, g2s)
		}
	}
}

func TestDisjExcl(t *testing.T) {
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return disjExcl(
			micro.EqualO(q, ast.NewInt(1)),
			micro.EqualO(q, ast.NewInt(2)),
		)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	for _, sexpr := range sexprs {
		if sexpr.String() != "1" {
			t.Fatalf("expected %s, but got %s instead", "1", sexpr.String())
		}
	}
}

// note: this is a naive implementation with runaway goroutines;
// in a proper implementation the goroutines are taken from a pool
// and properly cancelled when returning early
func concurrentDisjExcl(g1, g2 micro.Goal) micro.Goal {
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			c1 := make(chan micro.StreamOfStates, 1)
			c2 := make(chan micro.StreamOfStates, 1)
			go func() {
				c1 <- g1()(s)
			}()
			go func() {
				c2 <- g2()(s)
			}()
			var g1s, g2s micro.StreamOfStates
			for i := 0; i < 2; i++ {
				select {
				case cs := <-c1:
					if cs != nil {
						return cs
					}
					g1s = cs
				case cs := <-c2:
					if cs != nil {
						return cs
					}
					g2s = cs
				}
			}
			return micro.Mplus(g1s, g2s)
		}
	}
}

// illustrating or-parallelism for exclusive-or, where the first answer we find
// is the one we return.
func TestOrParallelismExcl(t *testing.T) {
	// use Zzz to delay evaluation of time.Sleep
	heavy := func(x *ast.SExpr) micro.Goal {
		return micro.Zzz(func() micro.GoalFn {
			time.Sleep(5 * time.Second)
			return micro.EqualO(
				x,
				ast.NewInt(2),
			)()
		})
	}

	sexprs := micro.Run(1, func(x *ast.SExpr) micro.Goal {
		return concurrentDisjExcl(
			heavy(x),
			micro.EqualO(
				x,
				ast.NewInt(1),
			),
		)
	})
	if len(sexprs) != 1 {
		t.Fatalf("expected len %d, but got len %d instead", 1, len(sexprs))
	}
	for _, sexpr := range sexprs {
		if sexpr.String() != "1" {
			t.Fatalf("expected %s, but got %s instead", "1", sexpr.String())
		}
	}
}

// executes both goals in parallel and awaits both results
// merging (using mplus) after. Also naive in terms of resource management
func concurrentDisjPlus(gs ...micro.Goal) micro.Goal {
	if len(gs) == 0 {
		return micro.FailureO
	}
	if len(gs) == 1 {
		return micro.Zzz(gs[0])
	}
	g1 := micro.Zzz(gs[0])
	g2 := concurrentDisjPlus(gs[1:]...)

	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			c1 := make(chan micro.StreamOfStates, 1)
			go func() {
				c1 <- g1()(s)
			}()
			g2s := g2()(s)
			var g1s micro.StreamOfStates
			select {
			case cs := <-c1:
				g1s = cs
			}
			return micro.Mplus(g1s, g2s)
		}
	}
}

func concurrentConde(gs ...[]micro.Goal) micro.Goal {
	conj := make([]micro.Goal, len(gs))
	for i, v := range gs {
		conj[i] = ConjPlus(v...)
	}
	return concurrentDisjPlus(conj...)
}

// this one is more general and useable, but hard to find a good example of it
// providing value in a simple testcase
// test with:
// go test ./... -v -run TestOrParallelism -count=1
// in order to see the time taken by both
func TestOrParallelism(t *testing.T) {
	//numGoals := 50000
	numGoals := 10000
	lotsOfGoals := []func(*ast.SExpr) micro.Goal{}
	for i := 0; i < numGoals; i++ {
		n := int64(i)
		f := func(x *ast.SExpr) micro.Goal {
			return micro.EqualO(x, ast.NewInt(n))
		}
		lotsOfGoals = append(lotsOfGoals, f)
	}
	generator := func(x *ast.SExpr) [][]micro.Goal {
		goals := make([][]micro.Goal, len(lotsOfGoals))
		for i, f := range lotsOfGoals {
			goals[i] = []micro.Goal{f(x)}
		}
		return goals
	}
	start := time.Now()
	sexprs := micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return concurrentConde(
			generator(q)...,
		)
	})
	took := time.Now().Sub(start)
	fmt.Println("concurrent", took)
	if len(sexprs) != numGoals {
		t.Fatalf("expected len %d, but got len %d instead", numGoals, len(sexprs))
	}
	for i, sexpr := range sexprs {
		if sexpr.String() != fmt.Sprint(i) {
			t.Fatalf("expected %s, but got %s instead", fmt.Sprint(i), sexpr.String())
		}
	}
	start = time.Now()
	sexprs = micro.Run(-1, func(q *ast.SExpr) micro.Goal {
		return Conde(
			generator(q)...,
		)
	})
	took = time.Now().Sub(start)
	fmt.Println("sequential", took)
	if len(sexprs) != numGoals {
		t.Fatalf("expected len %d, but got len %d instead", numGoals, len(sexprs))
	}
	for i, sexpr := range sexprs {
		if sexpr.String() != fmt.Sprint(i) {
			t.Fatalf("expected %s, but got %s instead", fmt.Sprint(i), sexpr.String())
		}
	}
}

// Sequential conj uses bind, which applies the second goal to all the states
// in the stream that results from applying the first goal to s
// this is inherently sequential. In order to make a concurrent conj,
// we need to come up with a version of bind that can merge two result streams
// NOTE: this is the only way i see and-parallelism working at the moment
// but it seems very counterintuitive to how miniKanren is set up
func concurrentConj(g1, g2 micro.Goal) micro.Goal {
	return func() micro.GoalFn {
		return func(s *micro.State) micro.StreamOfStates {
			c1 := make(chan micro.StreamOfStates, 1)
			go func() {
				c1 <- g1()(s)
			}()
			g2s := g2()(s)
			var g1s micro.StreamOfStates
			select {
			case cs := <-c1:
				g1s = cs
			}
			return mergeStreams(g1s, g2s)
		}
	}
}

// mergeStreams takes two streams, evaluates the cartesian product of states,
// and returns a stream that consists of the succesful union of pairs
// not only needs to merge states but also streams!
// since that is _super hard_ when dealing with infinite streams,
// I think the best way is to wrap each pair of the cartesian product in a function
// that does this unification when evaluated
// so the result is a stream of states (func() (state, streamofstates))
// which takes as a closure stream s1, s2 and the merge func??
func mergeStreams(s1, s2 micro.StreamOfStates) micro.StreamOfStates {
	// TODO discuss and implement
	return nil
}
