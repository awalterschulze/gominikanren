package comicro

// rewrite replaces all variables with their values, if it finds any in the substitutions map.
// It only replaces the variables that it finds on it's recursive walk, starting at the input variable.
func rewrite(v any, s *State) any {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return lookupValue(s, vvar)
		}
	}
	return Map(v, func(a any) any {
		return rewrite(a, s)
	})
}

// Reify finds reifications for the first introduced variable.
// This means it rewrites all substitutions for the first introduced variable.
// For any variables without substitutions, it adds a placeholder value.
func Reify(s *State) any {
	vvar := s.GetFirstVar()
	if vvar == nil {
		return nil
	}
	return rewrite(*vvar, s)
}
