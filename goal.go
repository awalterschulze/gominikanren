package gominikanren

import (
	"fmt"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

type Substitution = *ast.List

type StreamOfSubstitutions func() (Substitution, StreamOfSubstitutions)

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

// TODO compare IDs rather than names
func eqv(v, vv *ast.Variable) bool {
	fmt.Printf("(eqv? %v %v)\n", v, vv)
	return v.Name == vv.Name
}
