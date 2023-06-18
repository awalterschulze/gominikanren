package comicro

import (
	"context"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func s_x1() *State {
	s := NewEmptyState()
	var x *ast.SExpr
	s, x = NewVarWithName(s, "x", &ast.SExpr{})
	xvar, _ := s.GetVar(x)
	s = s.AddKeyValue(xvar, ast.NewSymbol("1"))
	return s
}

func s_xy_y1() *State {
	s := NewEmptyState()
	var x, y *ast.SExpr
	s, x = NewVarWithName(s, "x", &ast.SExpr{})
	s, y = NewVarWithName(s, "y", &ast.SExpr{})
	xvar, _ := s.GetVar(x)
	yvar, _ := s.GetVar(y)
	s = s.AddKeyValue(xvar, y)
	s = s.AddKeyValue(yvar, ast.NewSymbol("1"))
	return s
}

func s_x2() *State {
	s := NewEmptyState()
	var x *ast.SExpr
	s, x = NewVarWithName(s, "x", &ast.SExpr{})
	xvar, _ := s.GetVar(x)
	s = s.AddKeyValue(xvar, ast.NewSymbol("2"))
	return s
}

func empty() *State {
	return NewEmptyState()
}

func single(ctx context.Context, s *State) StreamOfStates {
	return NewSingletonStream(ctx, s)
}
