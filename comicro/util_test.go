package comicro

import (
	"context"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func s_x1() *State {
	return &State{
		Substitutions: Substitutions{
			indexOf(ast.NewVar("x", 1)): ast.NewSymbol("1"),
		},
		Counter: 2,
	}
}

func s_xy_y1() *State {
	return &State{
		Substitutions: Substitutions{
			indexOf(ast.NewVar("x", 1)): ast.NewVariable("y"),
			indexOf(ast.NewVar("y", 2)): ast.NewSymbol("1"),
		},
		Counter: 3,
	}
}

func s_x2() *State {
	return &State{
		Substitutions: Substitutions{
			indexOf(ast.NewVar("x", 1)): ast.NewSymbol("2"),
		},
		Counter: 2,
	}
}

func empty() *State {
	return NewEmptyState()
}

func single(ctx context.Context, s *State) StreamOfStates {
	return NewSingletonStream(ctx, s)
}
