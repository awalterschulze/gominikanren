// Code generated by goderive DO NOT EDIT.

package micro

import (
	ast "github.com/awalterschulze/gominikanren/sexpr/ast"
	"sort"
)

// deriveTuple3SVars returns a function, which returns the input values.
// Since tuples are not first class citizens in Go, this is a way to fake it, because functions that return tuples are first class citizens.
func deriveTuple3SVars(v0 *ast.SExpr, v1 *ast.SExpr, v2 Substitutions, v3 bool) func() (*ast.SExpr, *ast.SExpr, Substitutions, bool) {
	return func() (*ast.SExpr, *ast.SExpr, Substitutions, bool) {
		return v0, v1, v2, v3
	}
}

// deriveTupleE returns a function, which returns the input values.
// Since tuples are not first class citizens in Go, this is a way to fake it, because functions that return tuples are first class citizens.
func deriveTupleE(v0 *ast.SExpr, v1 *ast.SExpr, v2 Substitutions, v3 Substitutions) func() (*ast.SExpr, *ast.SExpr, Substitutions, Substitutions) {
	return func() (*ast.SExpr, *ast.SExpr, Substitutions, Substitutions) {
		return v0, v1, v2, v3
	}
}

// deriveTuple3 returns a function, which returns the input values.
// Since tuples are not first class citizens in Go, this is a way to fake it, because functions that return tuples are first class citizens.
func deriveTuple3(v0 string, v1 string, v2 string) func() (string, string, string) {
	return func() (string, string, string) {
		return v0, v1, v2
	}
}

// deriveTuple3SVar returns a function, which returns the input values.
// Since tuples are not first class citizens in Go, this is a way to fake it, because functions that return tuples are first class citizens.
func deriveTuple3SVar(v0 *ast.SExpr, v1 Substitutions, v2 string) func() (*ast.SExpr, Substitutions, string) {
	return func() (*ast.SExpr, Substitutions, string) {
		return v0, v1, v2
	}
}

// deriveSorted sorts the slice inplace and also returns it.
func deriveSorted(list []uint64) []uint64 {
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	return list
}

// deriveKeys returns the keys of the input map as a slice.
func deriveKeys(m map[uint64]*ast.SExpr) []uint64 {
	keys := make([]uint64, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// deriveFmapRs returns a list where each element of the input list has been morphed by the input function.
func deriveFmapRs(f func(*State) *ast.SExpr, list []*State) []*ast.SExpr {
	out := make([]*ast.SExpr, len(list))
	for i, elem := range list {
		out[i] = f(elem)
	}
	return out
}

// deriveFmaps returns a list where each element of the input list has been morphed by the input function.
func deriveFmaps(f func(*State) string, list []*State) []string {
	out := make([]string, len(list))
	for i, elem := range list {
		out[i] = f(elem)
	}
	return out
}
