package regex

import (
	"context"
	"reflect"
	"strings"
	"testing"

	. "github.com/awalterschulze/gominikanren/gomini"
)

func testo(t *testing.T, f func(q *Regex) Goal, want *Regex) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewState(CreateVarRegex)
	exprs := Run(ctx, 1, s, f)
	if len(exprs) != 1 {
		t.Fatalf("expected len %d result, but got %d instead", 1, len(exprs))
	}
	got := exprs[0]
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, but got %v instead", want, got)
	}
}

func CreateVarRegex(varTyp any, name string) (any, bool) {
	switch varTyp.(type) {
	case *Regex:
		chars := fmap([]rune(name), func(r rune) *Regex {
			return Char(r)
		})
		res := chars[len(chars)-1]
		for i := len(chars) - 2; i >= 0; i-- {
			res = Concat(chars[i], res)
		}
		return res, true
	case *rune:
		name, _ = strings.CutPrefix(name, ",")
		name, _ = strings.CutPrefix(name, "v")
		r := []rune(name)[0]
		return &r, true
	}
	return nil, false
}
