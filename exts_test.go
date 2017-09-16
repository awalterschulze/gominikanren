package gominikanren

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr"
)

func TestOccurs(t *testing.T) {
	tests := []func() (string, string, string, bool){
		deriveTupleO(",x", ",x", "`()", true),
		deriveTupleO(",x", ",y", "`()", false),
		deriveTupleO(",x", "`(,y)", "`((,y . ,x))", true),
	}
	for _, test := range tests {
		x, v, s, want := test()
		t.Run("(occurs "+x+" "+v+" "+s+")", func(t *testing.T) {
			xexpr, err := sexpr.Parse(x)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			sexpr, err := sexpr.Parse(s)
			if err != nil {
				t.Fatal(err)
			}
			got := occurs(xexpr, vexpr, sexpr.List)
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}

func TestExts(t *testing.T) {
	tests := []func() (string, string, string, string){
		deriveTupleE(",x", "a", "`()", "`((,x . a))"),
		deriveTupleE(",x", "`(,x)", "`()", ""),
		deriveTupleE(",x", "`(,y)", "`((,y . ,x))", ""),
	}
	for _, test := range tests {
		x, v, s, want := test()
		t.Run("(occurs "+x+" "+v+" "+s+")", func(t *testing.T) {
			xexpr, err := sexpr.Parse(x)
			if err != nil {
				t.Fatal(err)
			}
			vexpr, err := sexpr.Parse(v)
			if err != nil {
				t.Fatal(err)
			}
			sexpr, err := sexpr.Parse(s)
			if err != nil {
				t.Fatal(err)
			}
			gots, gotok := exts(xexpr, vexpr, sexpr.List)
			got := ""
			if gotok {
				got = gots.String()
			}
			if want != got {
				t.Fatalf("got %v want %v", got, want)
			}
		})
	}
}
