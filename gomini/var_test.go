package gomini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestIntIsVar(t *testing.T) {
	varTyp := int(0)
	s := NewEmptyState()
	s, v := NewVar(s, &varTyp)
	if _, ok := s.castVar(v); !ok {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsVar(t *testing.T) {
	s := NewEmptyState()
	s, v := NewVar(s, &ast.SExpr{})
	if _, ok := s.castVar(v); !ok {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsNotVar(t *testing.T) {
	s := NewEmptyState()
	x := ast.NewSymbol("x")
	if _, ok := s.castVar(x); ok {
		t.Fatalf("expected %v is not a var", x)
	}
}
