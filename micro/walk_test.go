package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

// u v w x y z
// 0 1 2 3 4 5
func TestWalk(t *testing.T) {
	zaxwyz := Substitutions{
		5: ast.NewSymbol("a"),
		3: ast.NewVar("w", 2),
		4: ast.NewVar("z", 5),
	}
	xyvxwx := Substitutions{
		3: ast.NewVar("y", 4),
		1: ast.NewVar("x", 3),
		2: ast.NewVar("x", 3),
	}
	tests := []func() (*ast.SExpr, Substitutions, string){
		deriveTuple3SVar(ast.NewVar("z", 5), zaxwyz, "a"),
		deriveTuple3SVar(ast.NewVar("y", 4), zaxwyz, "a"),
		deriveTuple3SVar(ast.NewVar("x", 3), zaxwyz, ",w"),
		deriveTuple3SVar(ast.NewVar("x", 3), xyvxwx, ",y"),
		deriveTuple3SVar(ast.NewVar("v", 1), xyvxwx, ",y"),
		deriveTuple3SVar(ast.NewVar("w", 2), xyvxwx, ",y"),
		deriveTuple3SVar(ast.NewVar("w", 2), Substitutions{
			3: ast.NewSymbol("b"),
			5: ast.NewVar("y", 4),
			2: ast.NewList(ast.NewVar("x", 3), ast.NewSymbol("e"), ast.NewVar("z", 5)),
		}, "(,x e ,z)"),
		deriveTuple3SVar(ast.NewVar("y", 4), Substitutions{
			3: ast.NewSymbol("e"),
			5: ast.NewVar("x", 3),
			4: ast.NewVar("z", 5),
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
	zaxwyz := Substitutions{
		5: ast.NewSymbol("a"),
		3: ast.NewVar("w", 2),
		4: ast.NewVar("z", 5),
	}
	xyvxwx := Substitutions{
		3: ast.NewVar("y", 4),
		1: ast.NewVar("x", 3),
		2: ast.NewVar("x", 3),
	}
	tests := []func() (*ast.SExpr, Substitutions, string){
		deriveTuple3SVar(ast.NewVar("z", 5), zaxwyz, "a"),
		deriveTuple3SVar(ast.NewVar("y", 4), zaxwyz, "a"),
		deriveTuple3SVar(ast.NewVar("x", 3), zaxwyz, ",w"),
		deriveTuple3SVar(ast.NewVar("x", 3), xyvxwx, ",y"),
		deriveTuple3SVar(ast.NewVar("v", 1), xyvxwx, ",y"),
		deriveTuple3SVar(ast.NewVar("w", 2), xyvxwx, ",y"),
		deriveTuple3SVar(ast.NewVar("w", 2), Substitutions{
			3: ast.NewSymbol("b"),
			5: ast.NewVar("y", 4),
			2: ast.NewList(ast.NewVar("x", 3), ast.NewSymbol("e"), ast.NewVar("z", 5)),
		}, "(b e ,y)"),
		deriveTuple3SVar(ast.NewVar("y", 4), Substitutions{
			3: ast.NewSymbol("e"),
			5: ast.NewVar("x", 3),
			4: ast.NewVar("z", 5),
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
