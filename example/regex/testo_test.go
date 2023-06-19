package regex

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/awalterschulze/gominikanren/comicro"
)

func testo(t *testing.T, f func(q *Regex) comicro.Goal, want *Regex) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exprs := comicro.Run(ctx, 1, VarCreator, f)
	if len(exprs) != 1 {
		t.Fatalf("expected len %d result, but got %d instead", 1, len(exprs))
	}
	got := exprs[0]
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, but got %v instead", want, got)
	}
}

func VarCreator(varTyp any, name string) any {
	switch varTyp.(type) {
	case *Regex:
		chars := fmap([]rune(name), func(r rune) *Regex {
			return Char(r)
		})
		res := chars[len(chars)-1]
		for i := len(chars) - 2; i >= 0; i-- {
			res = Concat(chars[i], res)
		}
		return res
	case *rune:
		name, _ = strings.CutPrefix(name, ",")
		name, _ = strings.CutPrefix(name, "v")
		return []rune(name)[0]
	}
	panic("unknown type")
}
