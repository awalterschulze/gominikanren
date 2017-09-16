package gominikanren

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
)

func TestEqual(t *testing.T) {
	tests := []func() (string, string, string){
		deriveTuple3("#f", "#f", "(`())"),
		deriveTuple3("#f", "#t", "()"),
	}
	for _, test := range tests {
		u, v, want := test()
		t.Run("((== "+u+" "+v+") empty-s)", func(t *testing.T) {
			uexpr, err := sexpr.Parse(u)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			stream := Equal(uexpr, vexpr)(EmptySubstitution())
			got := stream.String()
			if got != want {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestFailure(t *testing.T) {
	if got, want := Failure()(EmptySubstitution()).String(), "()"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestSuccess(t *testing.T) {
	if got, want := Success()(EmptySubstitution()).String(), "(`())"; got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
