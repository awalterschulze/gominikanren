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

// rewrite replaces all variables with their values, if it finds any in the substitutions map.
// It only replaces the variables that it finds on it's recursive walk, starting at the input variable.
func rewrite(v any, s *State) any {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return s.LookupValue(vvar)
		}
	}
	return Map(v, func(a any) any {
		return rewrite(a, s)
	})
}

// reifyFromState is a curried function that reifies the input variable for the given input state.
func reifyFromState(v Var, s *State) any {
	// rewrite all the variables given the substitutions map
	vv := rewrite(v, s)
	// add values to the substitutions map for all the variables that are left
	emptySubs := s.CopyWithoutSubstitutions()
	r := reifys(vv, emptySubs)
	// rewrite all the variables given the new substitutions map
	return rewrite(vv, r)
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
