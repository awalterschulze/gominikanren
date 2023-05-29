package comicro

import "github.com/awalterschulze/gominikanren/sexpr/ast"

func s_x1() *State {
	return &State{
		Substitutions: Substitutions{
			SubPair{indexOf(ast.NewVar("x", 1)), ast.NewSymbol("1")},
		},
		Counter: 2,
	}
}

func s_xy_y1() *State {
	return &State{
		Substitutions: Substitutions{
			SubPair{indexOf(ast.NewVar("x", 1)), ast.NewVariable("y")},
			SubPair{indexOf(ast.NewVar("y", 2)), ast.NewSymbol("1")},
		},
		Counter: 3,
	}
}

func s_x2() *State {
	return &State{
		Substitutions: Substitutions{
			SubPair{indexOf(ast.NewVar("x", 1)), ast.NewSymbol("2")},
		},
		Counter: 2,
	}
}

func empty() *State {
	return EmptyState()
}

func single(s *State) StreamOfStates {
	return NewSingletonStream(s)
}

func cons(s *State, ss StreamOfStates) StreamOfStates {
	return ConsStream(s, func() StreamOfStates {
		return ss
	})
}

func suspend(ss StreamOfStates) StreamOfStates {
	return Suspension(func() StreamOfStates {
		return ss
	})
}
