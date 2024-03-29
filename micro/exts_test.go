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
		tuple4(x, x, Substitutions(nil), true),
		tuple4(x, y, Substitutions(nil), false),
		tuple4(x, ast.NewList(y), Substitutions{
			SubPair{indexOf(y), x},
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
		tuple4(x, ast.NewSymbol("a"), Substitutions(nil), Substitutions{
			SubPair{indexOf(x), ast.NewSymbol("a")},
		}),
		tuple4(x, ast.NewList(x), Substitutions(nil), Substitutions(nil)),
		tuple4(x, ast.NewList(y),
			Substitutions{
				SubPair{indexOf(y), x},
			},
			Substitutions(nil)),
		tuple4(x, ast.NewSymbol("e"),
			Substitutions{
				SubPair{indexOf(z), x},
				SubPair{indexOf(y), z},
			},
			Substitutions{
				SubPair{indexOf(x), ast.NewSymbol("e")},
				SubPair{indexOf(z), x},
				SubPair{indexOf(y), z},
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
