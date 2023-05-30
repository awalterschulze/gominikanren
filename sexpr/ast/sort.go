package ast

import "sort"

type Results []*SExpr

func (r Results) Len() int {
	return len(r)
}

func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Results) Less(i, j int) bool {
	return r[i].Compare(r[j]) < 0
}

func Sort(exprs []*SExpr) []*SExpr {
	r := Results(exprs)
	sort.Sort(r)
	return r
}
