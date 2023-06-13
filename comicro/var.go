package comicro

import (
	"fmt"
	"reflect"

	"github.com/awalterschulze/gominikanren/sexpr/ast"
)

type Var uintptr

// Var creates a new variable as the string vC
func NewVar(c uint64) Var {
	return Var(uintptr(c))
}

func (v Var) SExpr() *ast.SExpr {
	return ast.NewVar(fmt.Sprintf("v%d", v), uint64(v))
}

func (v Var) String() string {
	return v.SExpr().String()
}

func GetVar(s any) (Var, bool) {
	if isvar(s) {
		return s.(Var), true
	}
	if isvarSExpr(s) {
		return NewVar(s.(*ast.SExpr).Atom.Var.Index), true
	}
	return 0, false
}

func IsVar(s any) bool {
	return isvar(s) || isvarSExpr(s)
}

func isvar(a any) bool {
	v := reflect.ValueOf(a)
	if v.Type() != varType {
		return false
	}
	return true
}

var varType = reflect.TypeOf(NewVar(0))

var sexprType = reflect.TypeOf(ast.SExpr{})
var atomType = reflect.TypeOf(ast.Atom{})

func isvarSExpr(s any) bool {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return false
	}
	if v.IsNil() {
		return false
	}
	v = v.Elem()
	k := v.Kind()
	switch k {
	case reflect.Struct:
		if v.Type() != sexprType {
			return false
		}
		atomValue := v.Field(1)
		if atomValue.IsNil() {
			return false
		}
		if atomValue.Elem().Type() != atomType {
			return false
		}
		varValue := atomValue.Elem().FieldByName("Var")
		if varValue.IsNil() {
			return false
		}
		return true
	}
	return false
}
