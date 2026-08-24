// Package sexpr implements a parser of symbolic expressions.
package sexpr

import (
	"awalterschulze.org/go/gominikanren/sexpr/ast"
	"awalterschulze.org/go/gominikanren/sexpr/lexer"
	"awalterschulze.org/go/gominikanren/sexpr/parser"
)

// Parse parses an symbol expression.
// It assumes that all expressions are already quoted, so variables need to be escaped with a comma.
func Parse(s string) (*ast.SExpr, error) {
	r, err := parser.NewParser().Parse(lexer.NewLexer([]byte(s)))
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	return r.(*ast.SExpr), nil
}
