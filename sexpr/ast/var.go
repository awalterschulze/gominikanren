package ast

func VarCreator(varTyp any, name string) any {
	return NewSymbol(name)
}
