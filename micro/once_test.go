package micro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestOnce(t *testing.T) {
	ifte := IfThenElseO(
		OnceO(DisjointO(
			EqualO(ast.NewSymbol("#t"), ast.NewVariable("x")),
			EqualO(ast.NewSymbol("#f"), ast.NewVariable("x")),
		)),
		EqualO(ast.NewSymbol("#f"), ast.NewVariable("y")),
		EqualO(ast.NewSymbol("#t"), ast.NewVariable("y")),
	)
	ss := ifte(EmptySubstitution())
	got := ss.String()
	want := "(`((,y . #f) (,x . #t)))"
	if got != want {
		t.Fatalf("got %v != want %v", got, want)
	}
}
