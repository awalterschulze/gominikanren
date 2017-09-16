package gominikanren

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	tests := []func() (string, string, string){
		deriveTuple3("z", "`((,z . a) (,x . ,w) (,y . ,z))", "a"),
		deriveTuple3("y", "`((,z . a) (,x . ,w) (,y . ,z))", "a"),
		deriveTuple3("x", "`((,z . a) (,x . ,w) (,y . ,z))", ",w"),
		deriveTuple3("x", "`((,x . ,y) (,v . ,x) (,w . ,x))", ",y"),
		deriveTuple3("v", "`((,x . ,y) (,v . ,x) (,w . ,x))", ",y"),
		deriveTuple3("w", "`((,x . ,y) (,v . ,x) (,w . ,x))", ",y"),
		deriveTuple3("w", "`((,x . b) (,z . ,y) (,w . (,x e ,z)))", "`(,x e ,z)"),
		deriveTuple3("y", "`((,x . e) (,z . ,x) (,y . ,z))", "e"),
	}
	for _, test := range tests {
		q, input, want := test()
		t.Run("(walk "+q+" "+input+")", func(t *testing.T) {
			s, err := sexpr.Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			v := ast.NewVariable(q)
			got := walk(v, s.List).String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}
