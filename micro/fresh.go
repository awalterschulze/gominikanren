package micro

import (
	"math/rand"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

/*
(define (call/fresh name f)
	(f (var name))
)
*/
// call/fresh second arguments is a function that expects a varaible and returns a Goal.
func CallFresh(name string, f func(v *ast.SExpr) Goal) Goal {
	// TODO this is where a variable should be assigned an ID
	// Could we do this purely or are we simply assigning a random number?
	v := ast.NewVariable(name)
	v.Atom.Var.ID = rand.Int63()
	return f(v)
}
