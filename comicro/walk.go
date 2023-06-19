package comicro

// Lookup returns the value of the variable,
// if it has any substitutions it recurses,
// otherwise it returns the final variable or final variable
// that it reaches in it's walk of the substitutions map.
func Lookup(key Var, s *State) any {
	value, ok := s.Get(key)
	if !ok {
		// There are no more substitutions to be made, this variable is the final value.
		return key
	}
	if newvar, ok := s.GetVar(value); ok {
		// If this value is a variable, then we need to Lookup it.
		return Lookup(newvar, s)
	}
	return value
}

// Replace replaces all variables with their values, if it finds any in the substitutions map.
// It only replaces the variables that it finds on it's recursive walk, starting at the input variable.
func ReplaceAll(v any, s *State) any {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			return s.LookupValue(vvar)
		}
	}
	return Map(v, func(a any) any {
		return ReplaceAll(a, s)
	})
}
