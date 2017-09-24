package mini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/micro"
	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOnce(t *testing.T) {
	ifte := IfThenElseO(
		OnceO(micro.DisjointO(
			micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("x")),
			micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("x")),
		)),
		micro.EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		micro.EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(micro.EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f) (,x . #t)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
