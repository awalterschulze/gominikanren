package comicro

func reifys(v any, s *State) *State {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return s.AddKeyValue(vvar, s.GetReifyName(vvar))
		}
	}
	return Fold(v, s, func(state *State, a any) *State {
		return reifys(a, state)
	})
}

// reifyFromState is a curried function that reifies the input variable for the given input state.
func reifyFromState(v Var, s *State) any {
	vv := ReplaceAll(v, s)
	emptySubs := s.CopyWithoutSubstitutions()
	r := reifys(vv, emptySubs)
	return ReplaceAll(vv, r)
}

// Reify finds reifications for the first introduced var
// NOTE: the way we've set this up now, vX is a reserved keyword
func Reify(s *State) any {
	vvar := s.GetFirstVar()
	if vvar == nil {
		return nil
	}
	return reifyFromState(*vvar, s)
}
