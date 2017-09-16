package micro

import (
	"github.com/awalterschulze/gominikanren/sexpr/ast"

	"strings"
)

type Substitution = *ast.List

type StreamOfSubstitutions func() (Substitution, StreamOfSubstitutions)

func (stream StreamOfSubstitutions) String() string {
	buf := []string{}
	var s Substitution
	for stream != nil {
		s, stream = stream()
		buf = append(buf, s.String())
	}
	return "(" + strings.Join(buf, " ") + ")"
}

type Goal func(Substitution) StreamOfSubstitutions

func Success() Goal {
	return func(s Substitution) StreamOfSubstitutions {
		return NewSingletonStream(s)
	}
}

func NewSingletonStream(s Substitution) StreamOfSubstitutions {
	return func() (Substitution, StreamOfSubstitutions) {
		return s, nil
	}
}

func Failure() Goal {
	return func(s Substitution) StreamOfSubstitutions {
		return nil
	}
}

func Suspension(s StreamOfSubstitutions) StreamOfSubstitutions {
	return func() (Substitution, StreamOfSubstitutions) {
		return nil, s
	}
}

func EmptySubstitution() Substitution {
	return ast.NewList(true).List
}

func Equal(u, v *ast.SExpr) Goal {
	return func(s Substitution) StreamOfSubstitutions {
		ss, sok := unify(u, v, s)
		if sok {
			return NewSingletonStream(ss)
		}
		return nil
	}
}

func eqv(s, ss *ast.SExpr) bool {
	return s.Equal(ss)
}