package ast

func CreateVar(varTyp any, name string) (any, bool) {
	switch varTyp.(type) {
	case *SExpr:
		return NewSymbol(name), true
	}
	return nil, false
}
