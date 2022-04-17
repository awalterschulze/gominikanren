package micro

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
		deriveTuple3SVar(z, zaxwyz, "a"),
		deriveTuple3SVar(y, zaxwyz, "a"),
		deriveTuple3SVar(x, zaxwyz, ",w"),
		deriveTuple3SVar(x, xyvxwx, ",y"),
		deriveTuple3SVar(v, xyvxwx, ",y"),
		deriveTuple3SVar(w, xyvxwx, ",y"),
		deriveTuple3SVar(w, Substitutions{
			SubPair{indexOf(x), ast.NewSymbol("b")},
			SubPair{indexOf(z), y},
			SubPair{indexOf(w), ast.NewList(x, ast.NewSymbol("e"), z)},
		}, "(,x e ,z)"),
		deriveTuple3SVar(y, Substitutions{
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
		deriveTuple3SVar(z, zaxwyz, "a"),
		deriveTuple3SVar(y, zaxwyz, "a"),
		deriveTuple3SVar(x, zaxwyz, ",w"),
		deriveTuple3SVar(x, xyvxwx, ",y"),
		deriveTuple3SVar(v, xyvxwx, ",y"),
		deriveTuple3SVar(w, xyvxwx, ",y"),
		deriveTuple3SVar(w, Substitutions{
			{indexOf(x), ast.NewSymbol("b")},
			{indexOf(z), y},
			{indexOf(w), ast.NewList(x, ast.NewSymbol("e"), z)},
		}, "(b e ,y)"),
		deriveTuple3SVar(y, Substitutions{
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
