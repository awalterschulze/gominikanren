package comicro

import (
	"context"
)

// CallFresh expects a function that expects a variable and returns a Goal.
func CallFresh(f func(Var) Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		s1, v := s.NewVar()
		f(v)(ctx, s1, ss)
	}
}

// calls n new fresh vars, it is up to the user to use them
// (ie by assigning them to go vars in f)
// not ideal, but better than 5 nested callfresh calls...
func Fresh(n int, f func(...Var) Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		s1 := s
		vars := make([]Var, n)
		for i := 0; i < n; i++ {
			s1, vars[i] = s1.NewVar()
		}
		f(vars...)(ctx, s1, ss)
	}
}
