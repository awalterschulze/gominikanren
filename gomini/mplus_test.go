package gomini

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func s_x1() *State {
	s := NewState()
	var x *ast.SExpr
	s, x = newVarWithName[*ast.SExpr](s, "x")
	xvar, _ := s.CastVar(x)
	s = s.Set(xvar, ast.NewSymbol("1"))
	return s
}

func s_xy_y1() *State {
	s := NewState()
	var x, y *ast.SExpr
	s, x = newVarWithName[*ast.SExpr](s, "x")
	s, y = newVarWithName[*ast.SExpr](s, "y")
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)
	s = s.Set(xvar, y)
	s = s.Set(yvar, ast.NewSymbol("1"))
	return s
}

func s_x2() *State {
	s := NewState()
	var x *ast.SExpr
	s, x = newVarWithName[*ast.SExpr](s, "x")
	xvar, _ := s.CastVar(x)
	s = s.Set(xvar, ast.NewSymbol("2"))
	return s
}

func empty() *State {
	return NewState()
}

// newSingletonStream returns the input state as a stream of states containing only the head state.
func newSingletonStream(ctx context.Context, s *State) Stream {
	ss := NewEmptyStream()
	Go(ctx, nil, func() {
		defer close(ss)
		ss.Write(ctx, s)
	})
	return ss
}

func single(ctx context.Context, s *State) Stream {
	return newSingletonStream(ctx, s)
}

func TestMplus1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s1 := single(ctx, s_x1())
	s2 := single(ctx, s_xy_y1())
	ss := NewEmptyStream()
	go func() {
		defer close(ss)
		Mplus(context.Background(), s1, s2, ss)
	}()
	count := 0
	for a := range ss {
		if a == nil {
			t.Fatalf("expected non nil state")
		}
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 states")
	}
}
