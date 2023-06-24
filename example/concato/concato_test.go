package concato

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
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

func toStrings(xs []any) []string {
	res := fmap(xs, func(a any) string {
		switch a := a.(type) {
		case *string:
			return *a
		default:
			return a.(Stringer).String()
		}
	})
	sort.Strings(res)
	return res
}

func TestPrependO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := comicro.NewEmptyState().WithVarCreators(CreateVar)
	gots := toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return comicro.EqualO(
			NewNode("a", NewNode("b", NewNode("c", nil))),
			q,
		)
	}))
	wants := []string{"[a,b,c]"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected list: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return PrependO(
			newString("a"),
			NewNode("b", NewNode("c", nil)),
			q,
		)
	}))
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected consed list: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *string) comicro.Goal {
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

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
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
	s := comicro.NewEmptyState().WithVarCreators(CreateVar)
	gots := toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
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

func TestConcatOAllCombinations(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := comicro.NewEmptyState().WithVarCreators(ast.CreateVar)
	gots := toStrings(comicro.Run(ctx, -1, s, func(q *Pair) comicro.Goal {
		return comicro.Exists(func(x *Node) comicro.Goal {
			return comicro.Exists(func(y *Node) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(
						&Pair{x, y},
						q,
					),
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
	s := comicro.NewEmptyState().WithVarCreators(ast.CreateVar)
	gots := toStrings(comicro.Run(ctx, -1, s, func(q *Pair) comicro.Goal {
		return comicro.Exists(func(x *Node) comicro.Goal {
			return comicro.Exists(func(y *Node) comicro.Goal {
				return comicro.Conj(
					comicro.EqualO(
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
