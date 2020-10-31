package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// helper func used in all tests that use Substitution of vars
func indexOf(x *ast.SExpr) uint64 {
	return x.Atom.Var.Index
}

func TestOccurs(t *testing.T) {
	x := ast.NewVar("x", 0)
	y := ast.NewVar("y", 1)
	tests := []func() (*ast.SExpr, *ast.SExpr, Substitutions, bool){
		deriveTuple3SVars(x, x, Substitutions(nil), true),
		deriveTuple3SVars(x, y, nil, false),
		deriveTuple3SVars(x, ast.NewList(y), Substitutions{
			indexOf(y): x,
		}, true),
	}
	for _, test := range tests {
		v, w, s, want := test()
		t.Run("(occurs "+v.String()+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := occurs(v.Atom.Var, w, s)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}

func TestExts(t *testing.T) {
	x := ast.NewVar("x", 0)
	y := ast.NewVar("y", 1)
	z := ast.NewVar("z", 2)
	tests := []func() (*ast.SExpr, *ast.SExpr, Substitutions, Substitutions){
		deriveTupleE(x, ast.NewSymbol("a"), Substitutions(nil), Substitutions{
			indexOf(x): ast.NewSymbol("a"),
		}),
		deriveTupleE(x, ast.NewList(x), Substitutions(nil), Substitutions(nil)),
		deriveTupleE(x, ast.NewList(y),
			Substitutions{
				indexOf(y): x,
			},
			Substitutions(nil)),
		deriveTupleE(x, ast.NewSymbol("e"),
			Substitutions{
				indexOf(z): x,
				indexOf(y): z,
			},
			Substitutions{
				indexOf(x): ast.NewSymbol("e"),
				indexOf(z): x,
				indexOf(y): z,
			},
		),
	}
	for _, test := range tests {
		v, w, s, want := test()
		t.Run("(exts "+v.String()+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := ""
			gots, gotok := exts(v.Atom.Var, w, s)
			if gotok {
				got = gots.String()
			}
			if want.String() != got {
				t.Fatalf("got %v <%#v> want %v", got, gots, want.String())
			}
		})
	}
}
