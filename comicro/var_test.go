package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestVarIsVar(t *testing.T) {
	s := NewEmptyState()
	s, v := s.NewVar()
	if _, ok := s.GetVar(v); !ok {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsVar(t *testing.T) {
	s := NewEmptyState()
	s, v := s.NewVar()
	if _, ok := s.GetVar(v.SExpr()); !ok {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsNotVar(t *testing.T) {
	s := NewEmptyState()
	x := ast.NewSymbol("x")
	if _, ok := s.GetVar(x); ok {
		t.Fatalf("expected %v is not a var", x)
	}
}
