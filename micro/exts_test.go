package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOccurs(t *testing.T) {
	tests := []func() (string, string, string, bool){
		deriveTupleO(",x", ",x", "()", true),
		deriveTupleO(",x", ",y", "()", false),
		deriveTupleO(",x", "(,y)", "((,y . ,x))", true),
	}
	for _, test := range tests {
		x, v, s, want := test()
		t.Run("(occurs "+x+" "+v+" "+s+")", func(t *testing.T) {
			xexpr, err := sexpr.Parse(x)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			sexpr, err := sexpr.Parse(s)
			if err != nil {
				t.Fatal(err)
			}
			var pair *ast.Pair
			if sexpr != nil {
				pair = sexpr.Pair
			}
			got := occurs(xexpr, vexpr, pair)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}

func TestExts(t *testing.T) {
	tests := []func() (string, string, string, string){
		deriveTupleE(",x", "a", "()", "((,x . a))"),
		deriveTupleE(",x", "(,x)", "()", ""),
		deriveTupleE(",x", "(,y)", "((,y . ,x))", ""),
		deriveTupleE(",x", "e", "((,z . ,x) (,y . ,z))", "((,x . e) (,z . ,x) (,y . ,z))"),
	}
	for _, test := range tests {
		x, v, s, want := test()
		t.Run("(exts "+x+" "+v+" "+s+")", func(t *testing.T) {
			xexpr, err := sexpr.Parse(x)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			sexpr, err := sexpr.Parse(s)
			if err != nil {
				t.Fatal(err)
			}
			var pair *ast.Pair
			if sexpr != nil {
				pair = sexpr.Pair
			}
			got := ""
			gots, gotok := exts(xexpr, vexpr, pair)
			if gotok {
				got = (&ast.SExpr{Pair: gots}).String()
			}
			if want != got {
				t.Fatalf("got %v <%#v> want %v", got, gots, want)
			}
		})
	}
}
