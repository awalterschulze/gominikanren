package comicro

import (
	"testing"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

func TestWalk(t *testing.T) {
	v := Var(1).SExpr()
	w := Var(2).SExpr()
	x := Var(3).SExpr()
	y := Var(4).SExpr()
	z := Var(5).SExpr()
	zaxwyz := Substitutions{
		indexOf(z): ast.NewSymbol("a"),
		indexOf(x): w,
		indexOf(y): z,
	}
	xyvxwx := Substitutions{
		indexOf(x): y,
		indexOf(v): x,
		indexOf(w): x,
	}
	tests := []func() (*ast.SExpr, Substitutions, string){
		tuple3(z, zaxwyz, "a"),
		tuple3(y, zaxwyz, "a"),
		tuple3(x, zaxwyz, w.String()),
		tuple3(x, xyvxwx, y.String()),
		tuple3(v, xyvxwx, y.String()),
		tuple3(w, xyvxwx, y.String()),
		tuple3(w, Substitutions{
			indexOf(x): ast.NewSymbol("b"),
			indexOf(z): y,
			indexOf(w): ast.NewList(x, ast.NewSymbol("e"), z),
		}, "("+x.String()+" e "+z.String()+")"),
		tuple3(y, Substitutions{
			indexOf(x): ast.NewSymbol("e"),
			indexOf(z): x,
			indexOf(y): z,
		}, "e"),
	}
	for _, test := range tests {
		v, subs, want := test()
		t.Run("(walk "+v.Atom.Var.Name+" "+subs.String()+")", func(t *testing.T) {
			got := Lookup(NewVar(v.Atom.Var.Index), subs).(interface{ String() string }).String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}

func TestReplaceAll(t *testing.T) {
	v := Var(1).SExpr()
	w := Var(2).SExpr()
	x := Var(3).SExpr()
	y := Var(4).SExpr()
	z := Var(5).SExpr()
	zaxwyz := Substitutions{
		indexOf(z): ast.NewSymbol("a"),
		indexOf(x): w,
		indexOf(y): z,
	}
	xyvxwx := Substitutions{
		indexOf(x): y,
		indexOf(v): x,
		indexOf(w): x,
	}
	tests := []func() (*ast.SExpr, Substitutions, string){
		tuple3(z, zaxwyz, "a"),
		tuple3(y, zaxwyz, "a"),
		tuple3(x, zaxwyz, w.String()),
		tuple3(x, xyvxwx, y.String()),
		tuple3(v, xyvxwx, y.String()),
		tuple3(w, xyvxwx, y.String()),
		tuple3(w, Substitutions{
			indexOf(x): ast.NewSymbol("b"),
			indexOf(z): y,
			indexOf(w): ast.NewList(x, ast.NewSymbol("e"), z),
		}, "(b e "+y.String()+")"),
		tuple3(y, Substitutions{
			indexOf(x): ast.NewSymbol("e"),
			indexOf(z): x,
			indexOf(y): z,
		}, "e"),
	}
	for _, test := range tests {
		v, subs, want := test()
		t.Run("(walk "+v.Atom.Var.Name+" "+subs.String()+")", func(t *testing.T) {
			got := ReplaceAll(v, subs).String()
			if want != got {
				t.Fatalf("got %s want %s", got, want)
			}
		})
	}
}
