package comicro

type varCreator func(varType any, val string) any

func reifys(v any, s *State, create varCreator) *State {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			varType := s.LookupValue(vvar)
			create(varType, s.GetName(vvar))
			return s.AddKeyValue(vvar, create(varType, ","+s.GetName(vvar)))
		}
	}
	return Fold(v, s, func(state *State, a any) *State {
		return reifys(a, state, create)
	})
}

// reifyFromState is a curried function that reifies the input variable for the given input state.
func reifyFromState(v Var, s *State, create varCreator) any {
	vv := ReplaceAll(v, s)
	noSubs := s.Copy()
	noSubs.Substitutions = make(map[Var]any)
	r := reifys(vv, noSubs, create)
	return ReplaceAll(vv, r)
}

// Reify finds reifications for the first introduced var
// NOTE: the way we've set this up now, vX is a reserved keyword
func Reify(s *State, create varCreator) any {
	vvar := s.GetFirstVar()
	if vvar == nil {
		return nil
	}
	return reifyFromState(*vvar, s, create)
}
