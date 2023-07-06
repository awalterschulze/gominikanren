package gomini

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestFuncsFmap(t *testing.T) {
	xs := []int{1, 2, 3}
	f := func(x int) int { return x + 1 }
	got := fmap(xs, f)
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestFuncsFmap_string(t *testing.T) {
	xs := []int{1, 2, 3}
	f := func(x int) string { return strconv.Itoa(x) }
	ys := fmap(xs, f)
	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(ys, want) {
		t.Fatalf("expected %v but got %v", want, ys)
	}
}

func TestFuncsAnyOf(t *testing.T) {
	xs := []int{1, 2, 3}
	f := func(x int) bool { return x == 2 }
	ys := anyOf(xs, f)
	want := true

	if ys != want {
		t.Fatalf("expected %v but got %v", want, ys)
	}
}

func TestFuncsZip(t *testing.T) {
	xs := []int{1, 2, 3}
	ys := []string{"1", "2", "3"}
	zs := zip(xs, ys)
	zsStrs := fmap(zs, func(f func() (int, string)) string {
		x, y := f()
		return fmt.Sprintf("%d:%s", x, y)
	})
	want := []string{"1:1", "2:2", "3:3"}
	if !reflect.DeepEqual(zsStrs, want) {
		t.Fatalf("expected %v but got %v", want, zsStrs)
	}
}

func TestFuncsZipReduce(t *testing.T) {
	xs := []int{1, 2, 3}
	ys := []int{1, 2, 3}
	zs, ok := zipReduce(xs, ys, true,
		func(x, y int, acc bool) (bool, bool) { return acc && x == y, true },
	)
	equal := zs && ok
	want := true
	if equal != want {
		t.Fatalf("expected %v but got %v", want, equal)
	}
}
