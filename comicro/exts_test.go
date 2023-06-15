package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// helper func used in all tests that use Substitution of vars
func indexOf(x *ast.SExpr) Var {
	return NewVar(x.Atom.Var.Index)
}

func TestOccurs(t *testing.T) {
	x := ast.NewVar("x", 0)
	y := ast.NewVar("y", 1)
	tests := []func() (*ast.SExpr, *ast.SExpr, Substitutions, bool){
		tuple4(x, x, Substitutions(nil), true),
		tuple4(x, y, Substitutions(nil), false),
		tuple4(x, ast.NewList(y), Substitutions{
			indexOf(y): x,
		}, true),
	}
	for _, test := range tests {
		v, w, subs, want := test()
		s := &State{Substitutions: subs, Counter: uint64(len(subs))}
		t.Run("(occurs "+v.String()+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := occurs(indexOf(v), w, s)
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
			indexOf(x): ast.NewSymbol("a"),
		}),
		tuple4(x, ast.NewList(x), Substitutions(nil), Substitutions(nil)),
		tuple4(x, ast.NewList(y),
			Substitutions{
				indexOf(y): x,
			},
			Substitutions(nil)),
		tuple4(x, ast.NewSymbol("e"),
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
		v, w, subs, want := test()
		s := &State{Substitutions: subs, Counter: uint64(len(subs))}
		t.Run("(exts "+v.String()+" "+w.String()+" "+s.String()+")", func(t *testing.T) {
			got := ""
			gots, gotok := exts(NewVar(v.Atom.Var.Index), w, s)
			if gotok {
				got = gots.Substitutions.String()
			}
			if want.String() != got {
				t.Fatalf("got %v <%#v> want %v", got, gots, want.String())
			}
		})
	}
}
