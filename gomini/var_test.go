package gomini

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestIntIsVar(t *testing.T) {
	s := NewState()
	s, v := NewVar[*int](s)
	if _, ok := s.CastVar(v); !ok {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsVar(t *testing.T) {
	s := NewState()
	s, v := NewVar[*ast.SExpr](s)
	if _, ok := s.CastVar(v); !ok {
		t.Fatalf("expected %v to be a var", v)
	}
}

func TestSExprIsNotVar(t *testing.T) {
	s := NewState()
	x := ast.NewSymbol("x")
	if _, ok := s.CastVar(x); ok {
		t.Fatalf("expected %v is not a var", x)
	}
}
