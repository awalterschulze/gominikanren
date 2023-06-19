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
