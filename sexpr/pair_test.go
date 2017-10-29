package sexpr

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestCons(t *testing.T) {
	tests := []func() (string, string, string){
		deriveTuple3("a", "b", "(a . b)"),
		deriveTuple3("(a)", "b", "((a) . b)"),
		deriveTuple3("(a)", "(b)", "((a) b)"),
		deriveTuple3("a", "(b)", "(a b)"),
		deriveTuple3("(a)", "(b c)", "((a) b c)"),
		deriveTuple3("a", "(b c)", "(a b c)"),
		deriveTuple3("(a b)", "(c d)", "((a b) c d)"),
		deriveTuple3("a", "()", "(a)"),
		deriveTuple3("(a)", "()", "((a))"),
		deriveTuple3("(a b)", "()", "((a b))"),
		deriveTuple3("()", "()", "(())"),
		deriveTuple3("()", "(b)", "(() b)"),
		deriveTuple3("()", "(b c)", "(() b c)"),
	}
	for _, test := range tests {
		car, cdr, want := test()
		t.Run("(cons "+car+" "+cdr+")", func(t *testing.T) {
			carexpr, err := Parse(car)
			if err != nil {
				t.Fatal(err)
			}
			cdrexpr, err := Parse(cdr)
			if err != nil {
				t.Fatal(err)
			}
			gotexpr := ast.Cons(carexpr, cdrexpr)
			got := gotexpr.String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
			gotcarexpr := gotexpr.Car()
			gotcar := gotcarexpr.String()
			if gotcar != car {
				t.Fatalf("car: got %s want %s", gotcar, car)
			}
			gotcdrexpr := gotexpr.Cdr()
			gotcdr := gotcdrexpr.String()
			if gotcdr != cdr {
				t.Fatalf("cdr: got %s want %s", gotcdr, cdr)
			}
		})
	}
}

func TestCar(t *testing.T) {
	tests := []func() (string, string){
		deriveTuple("(,z . a)", ",z"),
		deriveTuple(`(z a)`, "z"),
	}
	for _, test := range tests {
		input, want := test()
		t.Run(input, func(t *testing.T) {
			s, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			got := s.Car().String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestCdr(t *testing.T) {
	tests := []func() (string, string){
		deriveTuple("(,z . b)", "b"),
		deriveTuple("(,z . (,x e ,y))", "(,x e ,y)"),
		deriveTuple(`(z a)`, "(a)"),
		deriveTuple("((,z . b) (,x . ,y))", "((,x . ,y))"),
		deriveTuple("((,z . b) (,x . ,y) (,y . a))", "((,x . ,y) (,y . a))"),
	}
	for _, test := range tests {
		input, want := test()
		t.Run(input, func(t *testing.T) {
			s, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			got := s.Cdr().String()
			if want != got {
				t.Fatalf("got %s want %s parsed %s", got, want, s)
			}
		})
	}
}
