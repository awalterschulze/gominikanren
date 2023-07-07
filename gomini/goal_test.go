package gomini

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// NeverO is a Goal that returns a never ending stream of suspensions.
func NeverO(ctx context.Context, s *State, ss Stream) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
		}
	}
}

// AlwaysO is a goal that returns a never ending stream of success.
func AlwaysO(ctx context.Context, s *State, ss Stream) {
	for {
		if ok := ss.Write(ctx, s); !ok {
			return
		}
	}
}

func TestEqualO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tests := []func() (string, string, string){
		tuple3("#f", "#f", "((() . 0))"),
		tuple3("#f", "#t", "()"),
	}
	for _, test := range tests {
		u, v, want := test()
		t.Run("((== "+u+" "+v+") empty-s)", func(t *testing.T) {
			uexpr, err := sexpr.Parse(u)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			ss := NewStreamForGoal(ctx, EqualO(uexpr, vexpr), NewState())
			got := ss.String()
			if got != want {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestFailureO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, FailureO, NewState())
	if got, want := ss.String(), "()"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestSuccessO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, SuccessO, NewState())
	if got, want := ss.String(), "((() . 0))"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestDisj1(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	initial := NewState()
	var x *ast.SExpr
	initial, x = newVarWithName(initial, "x", &ast.SExpr{})
	ss := NewStreamForGoal(ctx,
		DisjO(
			EqualO(
				ast.NewSymbol("olive"),
				x,
			),
			NeverO,
		),
		initial,
	)
	var s *State
	for s = range ss {
		if s != nil {
			break
		}
	}
	got := s.String()
	// reifying y; we assigned it a random uint64 and lost track of it
	want := "({,x: olive} . 1)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	select {
	case _, ok := <-ss:
		if !ok {
			t.Fatalf("expected never ending")
		}
	default:
	}
}

func TestDisj2(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	initial := NewState()
	var x *ast.SExpr
	initial, x = newVarWithName(initial, "x", &ast.SExpr{})
	ss := NewStreamForGoal(ctx,
		DisjO(
			NeverO,
			EqualO(
				ast.NewSymbol("olive"),
				x,
			),
		),
		initial,
	)
	for s := range ss {
		if s != nil {
			got := s.String()
			// reifying y; we assigned it a random uint64 and lost track of it
			want := "({,x: olive} . 1)"
			if got != want {
				t.Fatalf("got %s != want %s", got, want)
			}
			break
		}
	}
	select {
	case _, ok := <-ss:
		if !ok {
			t.Fatalf("expected never ending")
		}
	default:
	}
}

func TestAlwaysO(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, AlwaysO, NewState())
	var s *State
	for s = range ss {
		if s != nil {
			break
		}
	}
	got := s.String()
	want := "(() . 0)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
	_, ok := <-ss
	if !ok {
		t.Fatalf("expected never ending")
	}
}

func runGoal(ctx context.Context, n int, s *State, g Goal) []*State {
	ss := NewStreamForGoal(ctx, g, s)
	return Take(ctx, n, ss)
}

func TestRunGoalAlways3(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := runGoal(ctx, 3, NewState(), AlwaysO)
	if len(ss) != 3 {
		t.Fatalf("expected 3 got %d", len(ss))
	}
	sss := fmap(ss, func(s *State) string {
		return s.String()
	})
	got := "(" + strings.Join(sss, " ") + ")"
	want := "((() . 0) (() . 0) (() . 0))"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestRunGoalDisj2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState()
	var x *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		x,
	)
	g := DisjO(e1, e2)
	stream := NewStreamForGoal(ctx, g, s)
	ss := Take(ctx, 5, stream)
	if len(ss) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2NoResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState()
	var x *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("oil"),
		x,
	)
	g := ConjO(e1, e2)
	stream := NewStreamForGoal(ctx, g, s)
	ss := Take(ctx, 5, stream)
	if len(ss) != 0 {
		t.Fatalf("expected 0, got %d: %v", len(ss), ss)
	}
}

func TestRunGoalConj2OneResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	initial := NewState()
	var x *ast.SExpr
	initial, x = newVarWithName(initial, "x", &ast.SExpr{})
	e1 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	e2 := EqualO(
		ast.NewSymbol("olive"),
		x,
	)
	g := ConjO(e1, e2)
	stream := NewStreamForGoal(ctx, g, initial)
	ss := Take(ctx, 5, stream)
	if len(ss) != 1 {
		t.Fatalf("expected 1, got %d: %v", len(ss), ss)
	}
	got := ss[0].String()
	// reifying y; we assigned it a random uint64 and lost track of it
	want := "({,x: olive} . 1)"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

/*
scheme code:

	(run-goal 1
		(call/fresh kiwi
			(lambda (fruit)
				(== plum fruit)
			)
		)
	)
*/
func TestExistsKiwi(t *testing.T) {
	s := NewState(ast.CreateVar)
	ss := runGoal(context.Background(), 1, s,
		ExistO(func(fruit *ast.SExpr) Goal {
			return EqualO(
				ast.NewSymbol("plum"),
				fruit,
			)
		},
		),
	)
	if len(ss) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(ss))
	}
	want := "({,v0: plum} . 1)"
	got := ss[0].String()
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestPineapple(t *testing.T) {
	pineapple := "pineapple"
	pointy := &pineapple
	reify := func(varTyp any, name string) (any, bool) {
		return &name, true
	}
	s := NewState(reify)
	ss := Run(context.Background(), 1, s,
		func(fruit *string) Goal {
			return EqualO(
				pointy,
				fruit,
			)
		},
	)
	if len(ss) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(ss))
	}
	got := *(ss[0]).(*string)
	want := pineapple
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestEqualOSlice(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, EqualO([]int{1, 2}, []int{1, 2}), NewState())
	got := ss.String()
	want := "((() . 0))"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestNotEqualOSlice(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ss := NewStreamForGoal(ctx, EqualO([]int{1, 2}, []int{1, 3}), NewState())
	got := ss.String()
	want := "()"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestEqualOVarSliceOfPointers(t *testing.T) {
	happy := "happy"
	pineapple := "pineapple"
	pineapples := []*string{&happy, &pineapple}
	reify := func(varTyp any, name string) (any, bool) {
		switch varTyp.(type) {
		case *string:
			return &name, true
		case []*string:
			return []*string{&name}, true
		}
		return nil, false
	}
	s := NewState(reify)
	ss := Run(context.Background(), 1, s,
		func(fruits []*string) Goal {
			return EqualO(
				pineapples,
				fruits,
			)
		},
	)
	if len(ss) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(ss))
	}
	got := (ss[0]).([]*string)
	want := pineapples
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestEqualOVarSliceOfStrings(t *testing.T) {
	happy := "happy"
	pineapple := "pineapple"
	pineapples := []string{happy, pineapple}
	reify := func(varTyp any, name string) (any, bool) {
		switch varTyp.(type) {
		case string:
			return name, true
		case []string:
			return []string{name}, true
		}
		return nil, false
	}
	s := NewState(reify)
	ss := Run(context.Background(), 1, s,
		func(fruits []string) Goal {
			return EqualO(
				pineapples,
				fruits,
			)
		},
	)
	if len(ss) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(ss))
	}
	got := (ss[0]).([]string)
	want := pineapples
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v != want %v", got, want)
	}
}

func TestEqualOSliceOfVarPointer(t *testing.T) {
	happy := "happy"
	pineapple := "pineapple"
	pineapples := []*string{&happy, &pineapple}
	reify := func(varTyp any, name string) (any, bool) {
		switch varTyp.(type) {
		case *string:
			return &name, true
		case []*string:
			return []*string{&name}, true
		}
		return nil, false
	}
	s := NewState(reify)
	ss := Run(context.Background(), 1, s,
		func(fruit *string) Goal {
			return EqualO(
				pineapples,
				[]*string{&happy, fruit},
			)
		},
	)
	if len(ss) != 1 {
		t.Fatalf("expected %d, but got %d results", 1, len(ss))
	}
	got := *(ss[0]).(*string)
	want := pineapple
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v != want %v", got, want)
	}
}
