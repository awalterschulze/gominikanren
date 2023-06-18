package comicro

import (
	"context"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func s_x1() *State {
	s := NewEmptyState()
	var x Var
	s, x = s.NewVarWithName("x")
	s = s.AddKeyValue(x, ast.NewSymbol("1"))
	return s
}

func s_xy_y1() *State {
	s := NewEmptyState()
	var x, y Var
	s, x = s.NewVarWithName("x")
	s, y = s.NewVarWithName("y")
	s = s.AddKeyValue(x, y)
	s = s.AddKeyValue(y, ast.NewSymbol("1"))
	return s
}

func s_x2() *State {
	s := NewEmptyState()
	var x Var
	s, x = s.NewVarWithName("x")
	s = s.AddKeyValue(x, ast.NewSymbol("2"))
	return s
}

func empty() *State {
	return NewEmptyState()
}

func single(ctx context.Context, s *State) StreamOfStates {
	return NewSingletonStream(ctx, s)
}
