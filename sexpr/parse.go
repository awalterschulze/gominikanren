package sexpr

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"
	"github.com/awalterschulze/gominikanren/sexpr/lexer"
	"github.com/awalterschulze/gominikanren/sexpr/parser"
)

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
