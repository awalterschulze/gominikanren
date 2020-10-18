package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	zaxwyz := Substitutions{
			"z": ast.NewSymbol("a"),
			"x": ast.NewVariable("w"),
			"y": ast.NewVariable("z"),
	}
	xyvxwx := Substitutions{
			"x": ast.NewVariable("y"),
			"v": ast.NewVariable("x"),
			"w": ast.NewVariable("x"),
	}
	tests := []func() (string, Substitutions, string){
		deriveTuple3S("z", zaxwyz, "a"),
		deriveTuple3S("y", zaxwyz, "a"),
		deriveTuple3S("x", zaxwyz, ",w"),
		deriveTuple3S("x", xyvxwx, ",y"),
		deriveTuple3S("v", xyvxwx, ",y"),
		deriveTuple3S("w", xyvxwx, ",y"),
		deriveTuple3S("w", Substitutions{
				"x": ast.NewSymbol("b"),
				"z": ast.NewVariable("y"),
				"w": ast.NewList(ast.NewVariable("x"), ast.NewSymbol("e"), ast.NewVariable("z")),
		}, "(,x e ,z)"),
		deriveTuple3S("y", Substitutions{
				"x": ast.NewSymbol("e"),
				"z": ast.NewVariable("x"),
				"y": ast.NewVariable("z"),
		}, "e"),
	}
	for _, test := range tests {
		q, subs, want := test()
		t.Run("(walk "+q+" "+subs.String()+")", func(t *testing.T) {
			v := ast.NewVariable(q)
			got := walk(v.Atom.Var, subs).String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestWalkStar(t *testing.T) {
	zaxwyz := Substitutions{
			"z": ast.NewSymbol("a"),
			"x": ast.NewVariable("w"),
			"y": ast.NewVariable("z"),
	}
	xyvxwx := Substitutions{
			"x": ast.NewVariable("y"),
			"v": ast.NewVariable("x"),
			"w": ast.NewVariable("x"),
	}
	tests := []func() (string, Substitutions, string){
		deriveTuple3S("z", zaxwyz, "a"),
		deriveTuple3S("y", zaxwyz, "a"),
		deriveTuple3S("x", zaxwyz, ",w"),
		deriveTuple3S("x", xyvxwx, ",y"),
		deriveTuple3S("v", xyvxwx, ",y"),
		deriveTuple3S("w", xyvxwx, ",y"),
		deriveTuple3S("w", Substitutions{
				"x": ast.NewSymbol("b"),
				"z": ast.NewVariable("y"),
				"w": ast.NewList(ast.NewVariable("x"), ast.NewSymbol("e"), ast.NewVariable("z")),
		}, "(b e ,y)"),
		deriveTuple3S("y", Substitutions{
				"x": ast.NewSymbol("e"),
				"z": ast.NewVariable("x"),
				"y": ast.NewVariable("z"),
		}, "e"),
	}
	for _, test := range tests {
		q, subs, want := test()
		t.Run("(walk "+q+" "+subs.String()+")", func(t *testing.T) {
			v := ast.NewVariable(q)
			got := walkStar(v, subs).String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}
