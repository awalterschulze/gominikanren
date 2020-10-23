package mini

import (
	"testing"
    "time"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
	"github.com/awalterschulze/gominikanren/micro"
)

// remember, concurrency != parallellism

// setting up exclusive-or for a very simple or-parallellism demo
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
            for i:=0; i<2; i++ {
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

// illustrating or-parallellism for exclusive-or, where the first answer we find
// is the one we return. More interesting would be to extend conde, but for that
// we would need to think a bit more on how to interleave stream results
func TestOrParallellism(t *testing.T) {
    // use Zzz to delay evaluation of time.Sleep
    heavy := func(x *ast.SExpr) micro.Goal {
        return micro.Zzz(func() micro.GoalFn {
            time.Sleep(5*time.Second)
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

// and-parallellism is breaking my brain, since in minikanren conj uses bind
// which is very explicit about sequential execution of underlying goals
// I don't think I understand bind/mplus well enough to play with that
