package gomini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestHasCycle(t *testing.T) {
	s := NewState()
	var x, y *ast.SExpr
	s, x = newVarWithName(s, "x", &ast.SExpr{})
	s, y = newVarWithName(s, "y", &ast.SExpr{})
	xvar, _ := s.CastVar(x)
	yvar, _ := s.CastVar(y)

	sxy := s.Set(yvar, x)
	tests := []func() (Var, *ast.SExpr, *State, bool){
		tuple4(xvar, x, s, true),
		tuple4(xvar, y, s, false),
		tuple4(xvar, ast.NewList(y), sxy, true),
	}
	for _, test := range tests {
		v, w, s, want := test()
		vname := s.getName(v)
		t.Run("(hasCycle "+vname+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := hasCycle(v, w, s)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}
