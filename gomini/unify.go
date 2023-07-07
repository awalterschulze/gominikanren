package gomini

import (
	"reflect"
)

// unify returns the substitution s extended with zero or more associations,
// unify returns nil, if there is a cycle
func unify(x, y any, s *State) *State {
	x = walk(x, s)
	y = walk(y, s)
	if reflect.DeepEqual(x, y) {
		return s
	} else if haveCycle(x, y, s) {
		return nil
	} else if xvar, ok := s.CastVar(x); ok {
		return s.Set(xvar, y)
	} else if yvar, ok := s.CastVar(y); ok {
		return s.Set(yvar, x)
	}
	return ZipReduce(x, y, s, unify)
}

// walk returns the value of the variable in the substitutions map.
// If the value is a variable, it walks, until it finds a final value or variable without a substitution.
func walk(x any, s *State) any {
	xvar, ok := s.CastVar(x)
	if !ok {
		// This is not a variable, so we can't walk any further.
		return x
	}
	xvalue, ok := s.Get(xvar)
	if !ok {
		// There are no more substitutions to be made, this variable is the final value.
		return x
	}
	return walk(xvalue, s)
}

// haveCycle checks if there is a cycle in either direction.
func haveCycle(x, y any, s *State) bool {
	if xvar, ok := s.CastVar(x); ok {
		if hasCycle(xvar, y, s) {
			return true
		}
	}
	if yvar, ok := s.CastVar(y); ok {
		if hasCycle(yvar, x, s) {
			return true
		}
	}
	return false
}

// hasCycle checks if there is a cycle from xvar to y.
func hasCycle(xvar Var, y any, s *State) bool {
	y = walk(y, s)
	if yvar, ok := s.CastVar(y); ok {
		return xvar == yvar
	}
	return Any(y, func(yelem any) bool {
		return hasCycle(xvar, yelem, s)
	})
}

// rewrite replaces all variables with their values, if it finds any in the substitutions map.
// It only replaces the variables that it finds on it's recursive walk, starting at the input variable.
func rewrite(x any, s *State) any {
	x = walk(x, s)
	if _, ok := s.CastVar(x); ok {
		// we walked down and found a variable, but we didn't find a substitution for it.
		return x
	}
	return Map(x, func(xelem any) any {
		return rewrite(xelem, s)
	})
}
