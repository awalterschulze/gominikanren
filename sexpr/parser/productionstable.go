// Code generated by gocc; DO NOT EDIT.

package parser

import (
	. "github.com/awalterschulze/gominikanren/sexpr/ast"
	"github.com/awalterschulze/gominikanren/sexpr/token"
)

func getStr(v interface{}) string {
	t := v.(*token.Token)
	return string(t.Lit)
}

func getSExpr(v interface{}) *SExpr {
	if v == nil {
		return nil
	}
	vv := v.(*SExpr)
	if vv.Pair == nil && vv.Atom == nil {
		return nil
	}
	return vv
}

type (
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String: `S' : SExpr	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `SExpr : Atom	<<  >>`,
		Id:         "SExpr",
		NTType:     1,
		Index:      1,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `SExpr : Pair	<<  >>`,
		Id:         "SExpr",
		NTType:     1,
		Index:      2,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Pair : "(" ")"	<< nil, nil >>`,
		Id:         "Pair",
		NTType:     2,
		Index:      3,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `Pair : "(" SExpr ")"	<< Cons(getSExpr(X[1]), nil), nil >>`,
		Id:         "Pair",
		NTType:     2,
		Index:      4,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return Cons(getSExpr(X[1]), nil), nil
		},
	},
	ProdTabEntry{
		String: `Pair : "(" SExpr space ContinueList ")"	<< Cons(getSExpr(X[1]), getSExpr(X[3])), nil >>`,
		Id:         "Pair",
		NTType:     2,
		Index:      5,
		NumSymbols: 5,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return Cons(getSExpr(X[1]), getSExpr(X[3])), nil
		},
	},
	ProdTabEntry{
		String: `Pair : "(" SExpr space "." space SExpr ")"	<< Cons(getSExpr(X[1]), getSExpr(X[5])), nil >>`,
		Id:         "Pair",
		NTType:     2,
		Index:      6,
		NumSymbols: 7,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return Cons(getSExpr(X[1]), getSExpr(X[5])), nil
		},
	},
	ProdTabEntry{
		String: `ContinueList : SExpr	<< Cons(getSExpr(X[0]), nil), nil >>`,
		Id:         "ContinueList",
		NTType:     3,
		Index:      7,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return Cons(getSExpr(X[0]), nil), nil
		},
	},
	ProdTabEntry{
		String: `ContinueList : SExpr space ContinueList	<< Cons(getSExpr(X[0]), getSExpr(X[2])), nil >>`,
		Id:         "ContinueList",
		NTType:     3,
		Index:      8,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return Cons(getSExpr(X[0]), getSExpr(X[2])), nil
		},
	},
	ProdTabEntry{
		String: `Atom : symbol	<< NewSymbol(getStr(X[0])), nil >>`,
		Id:         "Atom",
		NTType:     4,
		Index:      9,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return NewSymbol(getStr(X[0])), nil
		},
	},
	ProdTabEntry{
		String: `Atom : int_lit	<< ParseInt(getStr(X[0])) >>`,
		Id:         "Atom",
		NTType:     4,
		Index:      10,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ParseInt(getStr(X[0]))
		},
	},
	ProdTabEntry{
		String: `Atom : float_lit	<< ParseFloat(getStr(X[0])) >>`,
		Id:         "Atom",
		NTType:     4,
		Index:      11,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ParseFloat(getStr(X[0]))
		},
	},
	ProdTabEntry{
		String: `Atom : string_lit	<< ParseString(getStr(X[0])) >>`,
		Id:         "Atom",
		NTType:     4,
		Index:      12,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ParseString(getStr(X[0]))
		},
	},
	ProdTabEntry{
		String: `Atom : variable	<< ParseVariable(getStr(X[0])) >>`,
		Id:         "Atom",
		NTType:     4,
		Index:      13,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ParseVariable(getStr(X[0]))
		},
	},
}
