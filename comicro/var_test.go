package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestVarIsVar(t *testing.T) {
	v := NewVar(0)
	if !IsVar(v) {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsVar(t *testing.T) {
	v := NewVar(0).SExpr()
	if !IsVar(v) {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsNotVar(t *testing.T) {
	s := ast.NewSymbol("x")
	if IsVar(s) {
		t.Fatalf("expected %v is not a var", s)
	}
}
