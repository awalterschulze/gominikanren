package concato

import (
	"context"
	"reflect"
	"sort"
	"testing"

	. "github.com/awalterschulze/gominikanren/gomini"
)

func CreateVar(varTyp any, name string) (any, bool) {
	switch varTyp.(type) {
	case *string:
		return &name, true
	case *Node:
		return &Node{&name, nil}, true
	}
	return nil, false
}

func newString(s string) *string {
	return &s
}

type Pair struct {
	Left  *Node
	Right *Node
}

func (p *Pair) String() string {
	return `{` + p.Left.String() + `, ` + p.Right.String() + `}`
}

func fmap[A, B any](list []A, f func(A) B) []B {
	out := make([]B, len(list))
	for i, elem := range list {
		out[i] = f(elem)
	}
	return out
}

type stringer interface {
	String() string
}

func toStrings(xs []any) []string {
	res := fmap(xs, func(a any) string {
		switch a := a.(type) {
		case *string:
			return *a
		default:
			return a.(stringer).String()
		}
	})
	sort.Strings(res)
	return res
}

func TestPrependO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVar)
	gots := toStrings(RunTake(ctx, -1, s, func(q *Node) Goal {
		return EqualO(
			NewNode("a", NewNode("b", NewNode("c", nil))),
			q,
		)
	}))
	wants := []string{"[a,b,c]"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected list: got %v, want %v", gots, wants)
	}

	gots = toStrings(RunTake(ctx, -1, s, func(q *Node) Goal {
		return PrependO(
			newString("a"),
			NewNode("b", NewNode("c", nil)),
			q,
		)
	}))
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected consed list: got %v, want %v", gots, wants)
	}

	gots = toStrings(RunTake(ctx, -1, s, func(q *string) Goal {
		return PrependO(
			q,
			NewNode("b", NewNode("c", nil)),
			NewNode("a", NewNode("b", NewNode("c", nil))),
		)
	}))
	wants = []string{"a"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected consed element: got %v, want %v", gots, wants)
	}

	gots = toStrings(RunTake(ctx, -1, s, func(q *Node) Goal {
		return PrependO(
			newString("a"),
			q,
			NewNode("a", NewNode("b", NewNode("c", nil))),
		)
	}))
	wants = []string{"[b,c]"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected suffix list: got %v, want %v", gots, wants)
	}
}

func TestConcatOSingleList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVar)
	gots := toStrings(RunTake(ctx, -1, s, func(q *Node) Goal {
		return ConcatO(
			NewNode("a", nil),
			NewNode("b", nil),
			q,
		)
	}))
	wants := []string{"[a,b]"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("got %s != want %s", gots, wants)
	}
}

func TestConcatOOut(t *testing.T) {
	if testing.Short() {
		return
	}

	gots := toStrings(RunTake(context.Background(), -1, NewState(), func(q *Node) Goal {
		return ConcatO(
			NewNode("a", NewNode("b", nil)),
			NewNode("c", NewNode("d", nil)),
			q,
		)
	}))
	wants := []string{
		`[a,b,c,d]`,
	}

	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("got %s != want %s", gots, wants)
	}
}

func TestConcatOAllCombinations(t *testing.T) {
	if testing.Short() {
		return
	}

	gots := toStrings(RunTake(context.Background(), -1, NewState(), func(q *Pair) Goal {
		return ExistO(func(x *Node) Goal {
			return ExistO(func(y *Node) Goal {
				return ConjO(
					EqualO(&Pair{x, y}, q),
					ConcatO(
						x,
						y,
						NewNode("a", NewNode("b", NewNode("c", NewNode("d", nil)))),
					),
				)
			})
		})
	}))
	wants := []string{
		`{[], [a,b,c,d]}`,
		`{[a,b,c,d], []}`,
		`{[a,b,c], [d]}`,
		`{[a,b], [c,d]}`,
		`{[a], [b,c,d]}`,
	}

	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("got %s != want %s", gots, wants)
	}
}

func TestConcatOAllCombinationsCakeAndIcedTea(t *testing.T) {
	cakeandicedtea := NewNode("cake", NewNode("&", NewNode("ice", NewNode("d", NewNode("t", nil)))))
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVar)
	gots := toStrings(RunTake(ctx, -1, s, func(q *Pair) Goal {
		return ExistO(func(x *Node) Goal {
			return ExistO(func(y *Node) Goal {
				return ConjO(
					EqualO(
						&Pair{x, y},
						q,
					),
					ConcatO(
						x,
						y,
						cakeandicedtea,
					),
				)
			})
		})
	}))
	wants := []string{
		`{[], [cake,&,ice,d,t]}`,
		`{[cake,&,ice,d,t], []}`,
		`{[cake,&,ice,d], [t]}`,
		`{[cake,&,ice], [d,t]}`,
		`{[cake,&], [ice,d,t]}`,
		`{[cake], [&,ice,d,t]}`,
	}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("got %s != want %s", gots, wants)
	}
}
