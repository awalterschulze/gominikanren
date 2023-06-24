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

func ref[A any](s A) *A {
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
			&Node{ref("a"), &Node{ref("b"), nil}},
			q,
		)
	}))
	wants := []string{"[a,b]"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected list: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return PrependO(
			ref("a"),
			&Node{ref("b"), nil},
			q,
		)
	}))
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected consed list: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *string) comicro.Goal {
		return PrependO(
			q,
			&Node{ref("b"), nil},
			&Node{ref("a"), &Node{ref("b"), nil}},
		)
	}))
	wants = []string{"a"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected consed element: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return PrependO(
			ref("a"),
			q,
			&Node{ref("a"), &Node{ref("b"), nil}},
		)
	}))
	wants = []string{"[b]"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected suffix list: got %v, want %v", gots, wants)
	}
}

func TestConcatOSingleList(t *testing.T) {
	bs := &Node{ref("b"), nil}
	as := &Node{ref("a"), nil}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := comicro.NewEmptyState().WithVarCreators(CreateVar)
	gots := toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return ConcatO(
			as,
			bs,
			q,
		)
	}))
	wants := []string{"[a,b]"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("got %s != want %s", gots, wants)
	}
}

func TestConcatOAllCombinations(t *testing.T) {
	cakeandicedtea := &Node{ref("cake"), &Node{ref("&"), &Node{ref("ice"), &Node{ref("d"), &Node{ref("t"), nil}}}}}
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
