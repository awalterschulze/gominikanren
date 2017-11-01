package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	zaxwyz := Substitutions{
		&Substitution{
			Var:   "z",
			Value: ast.NewSymbol("a"),
		},
		&Substitution{
			Var:   "x",
			Value: ast.NewVariable("w"),
		},
		&Substitution{
			Var:   "y",
			Value: ast.NewVariable("z"),
		},
	}
	xyvxwx := Substitutions{
		&Substitution{
			Var:   "x",
			Value: ast.NewVariable("y"),
		},
		&Substitution{
			Var:   "v",
			Value: ast.NewVariable("x"),
		},
		&Substitution{
			Var:   "w",
			Value: ast.NewVariable("x"),
		},
	}
	tests := []func() (string, Substitutions, string){
		deriveTuple3S("z", zaxwyz, "a"),
		deriveTuple3S("y", zaxwyz, "a"),
		deriveTuple3S("x", zaxwyz, ",w"),
		deriveTuple3S("x", xyvxwx, ",y"),
		deriveTuple3S("v", xyvxwx, ",y"),
		deriveTuple3S("w", xyvxwx, ",y"),
		deriveTuple3S("w", Substitutions{
			&Substitution{
				Var:   "x",
				Value: ast.NewSymbol("b"),
			},
			&Substitution{
				Var:   "z",
				Value: ast.NewVariable("y"),
			},
			&Substitution{
				Var:   "w",
				Value: ast.NewList(ast.NewVariable("x"), ast.NewSymbol("e"), ast.NewVariable("z")),
			},
		}, "(,x e ,z)"),
		deriveTuple3S("y", Substitutions{
			&Substitution{
				Var:   "x",
				Value: ast.NewSymbol("e"),
			},
			&Substitution{
				Var:   "z",
				Value: ast.NewVariable("x"),
			},
			&Substitution{
				Var:   "y",
				Value: ast.NewVariable("z"),
			},
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
		&Substitution{
			Var:   "z",
			Value: ast.NewSymbol("a"),
		},
		&Substitution{
			Var:   "x",
			Value: ast.NewVariable("w"),
		},
		&Substitution{
			Var:   "y",
			Value: ast.NewVariable("z"),
		},
	}
	xyvxwx := Substitutions{
		&Substitution{
			Var:   "x",
			Value: ast.NewVariable("y"),
		},
		&Substitution{
			Var:   "v",
			Value: ast.NewVariable("x"),
		},
		&Substitution{
			Var:   "w",
			Value: ast.NewVariable("x"),
		},
	}
	tests := []func() (string, Substitutions, string){
		deriveTuple3S("z", zaxwyz, "a"),
		deriveTuple3S("y", zaxwyz, "a"),
		deriveTuple3S("x", zaxwyz, ",w"),
		deriveTuple3S("x", xyvxwx, ",y"),
		deriveTuple3S("v", xyvxwx, ",y"),
		deriveTuple3S("w", xyvxwx, ",y"),
		deriveTuple3S("w", Substitutions{
			&Substitution{
				Var:   "x",
				Value: ast.NewSymbol("b"),
			},
			&Substitution{
				Var:   "z",
				Value: ast.NewVariable("y"),
			},
			&Substitution{
				Var:   "w",
				Value: ast.NewList(ast.NewVariable("x"), ast.NewSymbol("e"), ast.NewVariable("z")),
			},
		}, "(b e ,y)"),
		deriveTuple3S("y", Substitutions{
			&Substitution{
				Var:   "x",
				Value: ast.NewSymbol("e"),
			},
			&Substitution{
				Var:   "z",
				Value: ast.NewVariable("x"),
			},
			&Substitution{
				Var:   "y",
				Value: ast.NewVariable("z"),
			},
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
