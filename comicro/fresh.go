package comicro

import (
	"context"
)

// CallFresh expects a function that expects a variable and returns a Goal.
func CallFresh(f func(Var) Goal) Goal {
	return func(ctx context.Context, s *State, ss StreamOfStates) {
		v := Var(s.Counter)
		s1 := s.AddCounter()
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
			v := Var(s1.Counter)
			vars[i] = v
			s1 = s1.AddCounter()
		}
		f(vars...)(ctx, s1, ss)
	}
}
