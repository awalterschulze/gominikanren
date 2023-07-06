package gomini

import (
	"context"
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func s_x1() *State {
	s := NewState()
	var x *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	xvar, _ := s.CastVar(x)
	s = s.Set(xvar, ast.NewSymbol("1"))
	return s
}

func s_xy_y1() *State {
	s := NewState()
	var x, y *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	s, y = newVarWithName(s, "y", &ast.SExpr{})
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)
	s = s.Set(xvar, y)
	s = s.Set(yvar, ast.NewSymbol("1"))
	return s
}

func s_x2() *State {
	s := NewState()
	var x *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	xvar, _ := s.CastVar(x)
	s = s.Set(xvar, ast.NewSymbol("2"))
	return s
}

func empty() *State {
	return NewState()
}

func single(ctx context.Context, s *State) StreamOfStates {
	return NewSingletonStream(ctx, s)
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
