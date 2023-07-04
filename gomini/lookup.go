package gomini

// Lookup returns the value of the variable in the substitutions map.
// If the value is a variable, it recurses, until it finds a final value or variable without a substitution.
func Lookup(key Var, s *State) any {
	value, ok := s.FindSubstitution(key)
	if !ok {
		// There are no more substitutions to be made, this variable is the final value.
		return key
	}
	if newvar, ok := s.CastVar(value); ok {
		// If this value is a variable, then we need to Lookup it.
		return Lookup(newvar, s)
	}
	return value
}
