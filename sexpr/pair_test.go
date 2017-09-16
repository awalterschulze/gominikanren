package sexpr

import "testing"

func TestCar(t *testing.T) {
	tests := []func() (string, string){
		deriveTuple("`(,z . a)", ",z"),
		deriveTuple(`(z a)`, "z"),
	}
	for _, test := range tests {
		input, want := test()
		t.Run(input, func(t *testing.T) {
			s, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			got := s.List.Car().String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestCdr(t *testing.T) {
	tests := []func() (string, string){
		deriveTuple("`(,z . b)", "b"),
		deriveTuple("`(,z . (,x e ,y))", "`(,x e ,y)"),
		deriveTuple(`(z a)`, "(a)"),
		deriveTuple("`((,z . b) (,x . ,y))", "`((,x . ,y))"),
		deriveTuple("`((,z . b) (,x . ,y) (,y . a))", "`((,x . ,y) (,y . a))"),
	}
	for _, test := range tests {
		input, want := test()
		t.Run(input, func(t *testing.T) {
			s, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			got := s.List.Cdr().String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}
