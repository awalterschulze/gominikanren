package micro

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

var vars = make(map[string]int64)

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
	id := int64(1)
	if c, ok := vars[name]; ok {
		id = c + 1
		vars[name] = id
	} else {
		vars[name] = id
	}
	v.Atom.Var.ID = id
	fmt.Printf("created: %v %v\n", v, v.Atom.Var.ID)
	return f(v)
}
