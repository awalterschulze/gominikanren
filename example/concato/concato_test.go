package concato

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"unsafe"

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

type Pair struct {
	Left  *Node
	Right *Node
}

func toStrings(xs []any) []string {
	res := fmap(xs, func(a any) string {
		switch a := a.(type) {
		case *string:
			return *a
		case *Node:
			return a.String()
		case *Pair:
			return `"` + a.Left.String() + `","` + a.Right.String() + `"`
		default:
			panic(fmt.Sprintf("unknown type %T", a))
		}
	})
	sort.Strings(res)
	return res
}

func TestSlicePointer(t *testing.T) {
	s := make([]int, 1)
	v := reflect.ValueOf(s)
	p := v.Pointer()
	fmt.Printf("%v\n", p)
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("%v\n", h.Data)
}

func TestPrependO(t *testing.T) {
	a := "a"
	b := "b"
	bs := &Node{&b, nil}
	as := &Node{&a, nil}
	abs := &Node{&a, &Node{&b, nil}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := comicro.NewEmptyState().WithVarCreators(CreateVar)
	gots := toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return comicro.EqualO(
			abs,
			q,
		)
	}))
	wants := toStrings([]any{abs})
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected list: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return PrependO(
			&a,
			bs,
			q,
		)
	}))
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected consed list: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *string) comicro.Goal {
		return PrependO(
			q,
			bs,
			abs,
		)
	}))
	wants = toStrings([]any{as})
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected consed element: got %v, want %v", gots, wants)
	}

	gots = toStrings(comicro.Run(ctx, -1, s, func(q *Node) comicro.Goal {
		return PrependO(
			&a,
			q,
			abs,
		)
	}))
	wants = toStrings([]any{bs})
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("expected suffix list: got %v, want %v", gots, wants)
	}
}

func TestConcatOSingleList(t *testing.T) {
	a := "a"
	b := "b"
	bs := &Node{&b, nil}
	as := &Node{&a, nil}
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
	wants := []string{"ab"}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("got %s != want %s", gots, wants)
	}
}

func TestConcatOAllCombinations(t *testing.T) {
	cake := "cake"
	and := "&"
	ice := "ice"
	d := "d"
	tea := "t"
	cakeandicedtea := &Node{&cake, &Node{&and, &Node{&ice, &Node{&d, &Node{&tea, nil}}}}}
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
		`"","cake&icedt"`,
		`"cake","&icedt"`,
		`"cake&","icedt"`,
		`"cake&ice","dt"`,
		`"cake&iced","t"`,
		`"cake&icedt",""`,
	}
	if !reflect.DeepEqual(gots, wants) {
		t.Fatalf("got %s != want %s", gots, wants)
	}
}
