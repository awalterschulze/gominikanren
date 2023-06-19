package ast

import (
	"reflect"
)

func ReifyName(varTyp any, name string) (any, bool) {
	typ := reflect.TypeOf(&SExpr{})
	typOfVarTyp := reflect.TypeOf(varTyp)
	if !typ.AssignableTo(typOfVarTyp) {
		return nil, false
	}
	return NewSymbol(name), true
}
