package gomini

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestFmap(t *testing.T) {
	xs := []int{1, 2, 3}
	f := func(x int) int { return x + 1 }
	got := fmap(xs, f)
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func ExampleFmap() {
	xs := []int{1, 2, 3}
	f := func(x int) int { return x + 1 }
	ys := fmap(xs, f)
	// ys = []int{1, 2, 3}

	fmt.Println(ys)
	// Output: [2 3 4]
}

func ExampleFmap_string() {
	xs := []int{1, 2, 3}
	f := func(x int) string { return strconv.Itoa(x) }
	ys := fmap(xs, f)
	// ys = []string{"1", "2", "3"}

	fmt.Printf("%#v\n", ys)
	// Output: []string{"1", "2", "3"}
}

func ExampleAnyOf() {
	xs := []int{1, 2, 3}
	f := func(x int) bool { return x == 2 }
	ys := anyOf(xs, f)
	// ys = true

	fmt.Println(ys)
	// Output: true
}

func ExampleZip() {
	xs := []int{1, 2, 3}
	ys := []string{"1", "2", "3"}
	zs := zip(xs, ys)
	zsStrs := fmap(zs, func(f func() (int, string)) string {
		x, y := f()
		return fmt.Sprintf("%d:%s", x, y)
	})

	fmt.Printf("%#v\n", zsStrs)
	// Output: []string{"1:1", "2:2", "3:3"}
}

func ExampleZipReduce() {
	xs := []int{1, 2, 3}
	ys := []int{1, 2, 3}
	zs, ok := zipReduce(xs, ys, true,
		func(x, y int, acc bool) (bool, bool) { return acc && x == y, true },
	)
	equal := zs && ok
	fmt.Println(equal)
	// Output: true
}
