package comicro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func reifys(v any, s *State) *State {
	if vvar, ok := s.GetVar(v); ok {
		v = Lookup(vvar, s)
		if vvar, ok := v.(Var); ok {
			// TODO: we need a way to create a struct that is not null, but have the name of a var there
			return s.AddKeyValue(vvar, ast.NewSymbol(","+s.GetName(vvar)))
		}
	}
	return Fold(v, s, func(state *State, a any) *State {
		return reifys(a, state)
	})
}

// reifyFromState is a curried function that reifies the input variable for the given input state.
func reifyFromState(v Var, s *State) any {
	vv := ReplaceAll(v, s)
	noSubs := s.Copy()
	noSubs.Substitutions = make(map[Var]any)
	r := reifys(vv, noSubs)
	return ReplaceAll(vv, r)
}

// Reify finds reifications for the first introduced var
// NOTE: the way we've set this up now, vX is a reserved keyword
func Reify(s *State, typ any) any {
	vvar := s.GetFirstVar()
	if vvar == nil {
		return typ
	}
	return reifyFromState(*vvar, s)
}
