package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOccurs(t *testing.T) {
	tests := []func() (string, string, Substitutions, bool){
		deriveTupleO(",x", ",x", Substitutions(nil), true),
		deriveTupleO(",x", ",y", nil, false),
		deriveTupleO(",x", "(,y)", Substitutions{
            "y": ast.NewVariable("x"),
		}, true),
	}
	for _, test := range tests {
		x, v, s, want := test()
		t.Run("(occurs "+x+" "+v+" "+s.String()+")", func(t *testing.T) {
			xexpr, err := sexpr.Parse(x)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			got := occurs(xexpr.Atom.Var, vexpr, s)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}

func TestExts(t *testing.T) {
	tests := []func() (string, string, Substitutions, Substitutions){
		deriveTupleE(",x", "a", Substitutions(nil), Substitutions{
			"x": ast.NewSymbol("a"),
		}),
		deriveTupleE(",x", "(,x)", Substitutions(nil), Substitutions(nil)),
		deriveTupleE(",x", "(,y)",
			Substitutions{
				"y": ast.NewVariable("x"),
			},
			Substitutions(nil)),
		deriveTupleE(",x", "e",
			Substitutions{
					"z": ast.NewVariable("x"),
					"y": ast.NewVariable("z"),
				},
            Substitutions{
					"x": ast.NewSymbol("e"),
					"z": ast.NewVariable("x"),
					"y": ast.NewVariable("z"),
			},
		),
	}
	for _, test := range tests {
		x, v, s, want := test()
		t.Run("(exts "+x+" "+v+" "+s.String()+")", func(t *testing.T) {
			xexpr, err := sexpr.Parse(x)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			got := ""
			gots, gotok := exts(xexpr.Atom.Var, vexpr, s)
			if gotok {
				got = gots.String()
			}
			if want.String() != got {
				t.Fatalf("got %v <%#v> want %v", got, gots, want.String())
			}
		})
	}
}
