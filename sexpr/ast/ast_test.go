package ast

import (
	"testing"
)

type test struct {
	desc  string
	input *SExpr
	want  string
}

func TestCons(t *testing.T) {
	tests := []test{
		{
			desc:  "var",
			input: NewVariable("a"),
			want:  ",a",
		},
		{
			desc: "1",
			input: Cons(NewVariable("a"),
				Cons(NewVariable("b"), nil),
			),
			want: "(,a ,b)",
		},
		{
			desc: "1 cdr",
			input: Cons(NewVariable("a"),
				Cons(NewVariable("b"), nil),
			).Cdr(),
			want: "(,b)",
		},
		{
			desc:  "2",
			input: Cons(NewVariable("a"), NewVariable("b")),
			want:  "(,a . ,b)",
		},
		{
			desc:  "2 cdr",
			input: Cons(NewVariable("a"), NewVariable("b")).Cdr(),
			want:  ",b",
		},
		{
			desc: "3",
			input: Cons(NewVariable("a"),
				Cons(NewVariable("b"),
					Cons(NewVariable("c"), nil),
				),
			),
			want: "(,a ,b ,c)",
		},
		{
			desc: "3 cdr",
			input: Cons(NewVariable("a"),
				Cons(NewVariable("b"),
					Cons(NewVariable("c"), nil),
				),
			).Cdr(),
			want: "(,b ,c)",
		},
		{
			desc: "4",
			input: Cons(NewVariable("a"),
				Cons(Cons(NewVariable("b"),
					Cons(NewVariable("c"), nil),
				), nil),
			),
			want: "(,a (,b ,c))",
		},
		{
			desc: "4 cdr",
			input: Cons(NewVariable("a"),
				Cons(Cons(NewVariable("b"),
					Cons(NewVariable("c"), nil),
				), nil),
			).Cdr(),
			want: "((,b ,c))",
		},
		{
			desc: "5",
			input: Cons(NewVariable("a"),
				Cons(Cons(NewVariable("b"), nil), nil),
			),
			want: "(,a (,b))",
		},
		{
			desc: "5 cdr",
			input: Cons(NewVariable("a"),
				Cons(Cons(NewVariable("b"), nil), nil),
			).Cdr(),
			want: "((,b))",
		},
		{
			desc: "6",
			input: Cons(NewVariable("a"),
				Cons(NewVariable("b"), nil),
			),
			want: "(,a ,b)",
		},
		{
			desc: "6 cdr",
			input: Cons(NewVariable("a"),
				Cons(NewVariable("b"), nil),
			).Cdr(),
			want: "(,b)",
		},
	}
	for _, test := range tests {
		got := test.input.String()
		if got != test.want {
			t.Errorf("%s: got <%s> != want <%s>", test.desc, got, test.want)
		}
	}
}
