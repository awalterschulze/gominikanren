package gominikanren

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	tests := []func() (string, string, string){
		deriveTuple("z", "`((,z . a) (,x . ,w) (,y . ,z))", "a"),
		deriveTuple("y", "`((,z . a) (,x . ,w) (,y . ,z))", "a"),
		deriveTuple("x", "`((,z . a) (,x . ,w) (,y . ,z))", ",w"),
		deriveTuple("x", "`((,x . ,y) (,v . ,x) (,w . ,x))", ",y"),
		deriveTuple("v", "`((,x . ,y) (,v . ,x) (,w . ,x))", ",y"),
		deriveTuple("w", "`((,x . ,y) (,v . ,x) (,w . ,x))", ",y"),
		deriveTuple("w", "`((,x . b) (,z . ,y) (,w . (,x e ,z)))", "`(,x e ,z)"),
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
