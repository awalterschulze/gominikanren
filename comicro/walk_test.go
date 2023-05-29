package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	v := ast.NewVar("v", 1)
	w := ast.NewVar("w", 2)
	x := ast.NewVar("x", 3)
	y := ast.NewVar("y", 4)
	z := ast.NewVar("z", 5)
	zaxwyz := Substitutions{
		SubPair{indexOf(z), ast.NewSymbol("a")},
		SubPair{indexOf(x), w},
		SubPair{indexOf(y), z},
	}
	xyvxwx := Substitutions{
		SubPair{indexOf(x), y},
		SubPair{indexOf(v), x},
		SubPair{indexOf(w), x},
	}
	tests := []func() (*ast.SExpr, Substitutions, string){
		tuple3(z, zaxwyz, "a"),
		tuple3(y, zaxwyz, "a"),
		tuple3(x, zaxwyz, ",w"),
		tuple3(x, xyvxwx, ",y"),
		tuple3(v, xyvxwx, ",y"),
		tuple3(w, xyvxwx, ",y"),
		tuple3(w, Substitutions{
			SubPair{indexOf(x), ast.NewSymbol("b")},
			SubPair{indexOf(z), y},
			SubPair{indexOf(w), ast.NewList(x, ast.NewSymbol("e"), z)},
		}, "(,x e ,z)"),
		tuple3(y, Substitutions{
			SubPair{indexOf(x), ast.NewSymbol("e")},
			SubPair{indexOf(z), x},
			SubPair{indexOf(y), z},
		}, "e"),
	}
	for _, test := range tests {
		v, subs, want := test()
		t.Run("(walk "+v.Atom.Var.Name+" "+subs.String()+")", func(t *testing.T) {
			got := walk(v.Atom.Var, subs).String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestWalkStar(t *testing.T) {
	v := ast.NewVar("v", 1)
	w := ast.NewVar("w", 2)
	x := ast.NewVar("x", 3)
	y := ast.NewVar("y", 4)
	z := ast.NewVar("z", 5)
	zaxwyz := Substitutions{
		{indexOf(z), ast.NewSymbol("a")},
		{indexOf(x), w},
		{indexOf(y), z},
	}
	xyvxwx := Substitutions{
		{indexOf(x), y},
		{indexOf(v), x},
		{indexOf(w), x},
	}
	tests := []func() (*ast.SExpr, Substitutions, string){
		tuple3(z, zaxwyz, "a"),
		tuple3(y, zaxwyz, "a"),
		tuple3(x, zaxwyz, ",w"),
		tuple3(x, xyvxwx, ",y"),
		tuple3(v, xyvxwx, ",y"),
		tuple3(w, xyvxwx, ",y"),
		tuple3(w, Substitutions{
			{indexOf(x), ast.NewSymbol("b")},
			{indexOf(z), y},
			{indexOf(w), ast.NewList(x, ast.NewSymbol("e"), z)},
		}, "(b e ,y)"),
		tuple3(y, Substitutions{
			{indexOf(x), ast.NewSymbol("e")},
			{indexOf(z), x},
			{indexOf(y), z},
		}, "e"),
	}
	for _, test := range tests {
		v, subs, want := test()
		t.Run("(walk "+v.Atom.Var.Name+" "+subs.String()+")", func(t *testing.T) {
			got := walkStar(v, subs).String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}
