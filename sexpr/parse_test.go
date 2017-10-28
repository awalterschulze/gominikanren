package sexpr

import (
	"strings"
	"testing"
)

var pass = []string{
	`(x . y)`,
	`(x y z)`,
	`(1 () (2 . 3) (4))`,
	`with-hyphen`,
	`?@!$`,
	`a\ symbol\ with\ spaces`,
	`"Hello, world!"`,
	`-9876543210`,
	`-0.01`,
	`6.28318`,
	`6.023e23`,
	`(((S) (NP VP))
	 ((VP) (V))
	 ((VP) (V NP))
	 ((V) died)
	 ((V) employed)
	 ((NP) nurses)
	 ((NP) patients)
	 ((NP) Medicenter)
	 ((NP) "Dr Chan"))`,
	`(foo bar "1" "quote non tokens" true)`,
	`(sub (sub) (sub (sub)))`,
	`(,x y)`,
}

func TestParseSuccess(t *testing.T) {
	for _, p := range pass {
		t.Run(p, func(t *testing.T) {
			sexpr, err := Parse(p)
			if err != nil {
				t.Fatal(err)
			}
			got := sexpr.String()
			pp := strings.Split(p, " ")
			wants := []string{}
			for i := range pp {
				pp[i] = strings.Replace(pp[i], "\n", "", -1)
				pp[i] = strings.Replace(pp[i], "\t", "", -1)
				if len(pp[i]) > 0 {
					wants = append(wants, pp[i])
				}
			}
			want := strings.Join(wants, " ")
			if got != want {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestParseDiff(t *testing.T) {
	input := `(x . (y . (z . ())))`
	want := "(x y z)"
	sexpr, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	got := sexpr.String()
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

var fail = []string{
	"(",
	")",
	"())",
	"(()",
	"(()))",
	"((())",
	"(()()()()))",
}

func TestParseFailure(t *testing.T) {
	for _, f := range fail {
		t.Run(f, func(t *testing.T) {
			_, err := Parse(f)
			if err == nil {
				t.Fatal("expected failure")
			}
		})
	}
}
