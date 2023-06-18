package comicro

func reifys(v any, s *State) *State {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if _, ok := v.(Var); ok {
			// terminate when a var has no mapping
			return s
		}
	}
	return Fold(v, s, func(state *State, a any) *State {
		return reifys(a, state)
	})
}

// reifyFromState is a curried function that reifies the input variable for the given input state.
func reifyFromState(v Var, s *State) any {
	vv := ReplaceAll(v, s)
	r := reifys(vv, nil)
	return ReplaceAll(vv, r)
}

// Reify finds reifications for the first introduced var
// NOTE: the way we've set this up now, vX is a reserved keyword
func Reify(s *State) any {
	return reifyFromState(0, s)
}
