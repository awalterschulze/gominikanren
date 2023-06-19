package comicro

func addPlaceHolders(v any, s *State) *State {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return s.AddKeyValue(vvar, s.GetReifyName(vvar))
		}
	}
	return Fold(v, s, addPlaceHolders)
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
func reifyFromState(start Var, state *State) any {
	// rewrite all the variables given the substitutions map
	substitutedValue := rewrite(start, state)
	// add values to the substitutions map for all the variables that are left
	placeholders := addPlaceHolders(substitutedValue, state.CopyWithoutSubstitutions())
	// rewrite all the variables given the new placeholder substitutions map
	return rewrite(substitutedValue, placeholders)
}

// Reify finds reifications for the first introduced variable.
// This means it rewrites all substitutions for the first introduced variable.
// For any variables without substitutions, it adds a placeholder value.
func Reify(s *State) any {
	vvar := s.GetFirstVar()
	if vvar == nil {
		return nil
	}
	return reifyFromState(*vvar, s)
}
