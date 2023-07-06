package gomini

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func ExampleMap() {
	xs := []int{1, 2, 3}
	f := func(x any) any { return x.(int) + 1 }
	ys := Map(xs, f).([]int)
	fmt.Println(ys)
	// Output: [2 3 4]
}

func TestMapInc(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(x any) any { return x.(int) + 1 }).([]int)
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

type A struct {
	Field1 int
	Field2 string
}

func TestMapStruct(t *testing.T) {
	got := Map(&A{1, "1"}, func(x any) any {
		switch x := x.(type) {
		case int:
			return x + 1
		case string:
			return x + "+1"
		}
		return x
	}).(*A)
	want := &A{2, "1+1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func ExampleMap_struct() {
	xs := &A{1, "1"}
	f := func(x any) any {
		switch x := x.(type) {
		case int:
			return x + 1
		case string:
			return x + "+1"
		}
		return x
	}
	ys := Map(xs, f).(*A)
	fmt.Printf("%#v\n", ys)
	// Output: &gomini.A{Field1:2, Field2:"1+1"}
}

func TestAnyIntTrue(t *testing.T) {
	got := Any([]int{1, 2, 3}, func(x any) bool { return x.(int) == 2 })
	want := true
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestAnyIntFalse(t *testing.T) {
	got := Any([]int{1, 2, 3}, func(x any) bool { return x.(int) == 4 })
	want := false
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestAnyStruct(t *testing.T) {
	got := Any(&A{1, "2"}, func(x any) bool {
		switch x := x.(type) {
		case int:
			return false
		case string:
			return x == "2"
		}
		return false
	})
	want := true
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

type AStruct struct {
	Field1 int
	Field2 string
	Field3 int
}

func ExampleAny_struct() {
	xs := &AStruct{1, "2", 3}
	f := func(x any) bool {
		switch x := x.(type) {
		case int:
			return x == 2
		case string:
			return x == "2"
		}
		return false
	}
	isAny := Any(xs, f)
	// isAny = true
	fmt.Println(isAny)
	// Output: true
}

func TestZipReduceSum(t *testing.T) {
	got := ZipReduce([]int{1, 2, 3}, []int{3, 2, 1}, option(0), func(x1, x2 any, sum *int) *int { return option(*sum + x1.(int) + x2.(int)) })
	want := 12
	if *got != want {
		t.Fatalf("expected %v but got %v", want, *got)
	}
}

func TestZipReduceConcat(t *testing.T) {
	got := ZipReduce([]int{1, 2, 3}, []int{3, 2, 1}, option(""), func(x1, x2 any, s *string) *string { return option(*s + fmt.Sprintf("%d%d", x1, x2)) })
	want := "132231"
	if *got != want {
		t.Fatalf("expected %v but got %v", want, *got)
	}
}

func TestZipReduceFalse(t *testing.T) {
	got := ZipReduce([]int{1, 2, 3}, []int{3, 2, 1}, option(0),
		func(x1, x2 any, sum *int) *int {
			if x1 == 2 {
				return nil
			} else {
				return option(*sum + x1.(int) + x2.(int))
			}
		})
	if got != nil {
		t.Fatalf("expected nil but got %v", got)
	}
}

func TestZipReduceStruct(t *testing.T) {
	got := ZipReduce(&A{1, "2"}, &A{2, "3"}, option(0), func(x1, x2 any, sum *int) *int {
		switch x1 := x1.(type) {
		case int:
			return option(*sum + x1 + x2.(int))
		case string:
			x1i, err := strconv.Atoi(x1)
			if err != nil {
				return nil
			}
			x2i, err := strconv.Atoi(x2.(string))
			if err != nil {
				return nil
			}
			return option(*sum + x1i + x2i)
		}
		return nil
	})
	want := 8
	if *got != want {
		t.Fatalf("expected %v but got %v", want, *got)
	}
}

type MyStruct struct {
	Field1 int
	Field2 string
	More   *MyStruct
}

func DeepEqualOriginal(x, y any) bool {
	equalValues := ZipReduce(x, y, option(true),
		func(xfield, yfield any, eq *bool) *bool {
			switch xfield := xfield.(type) {
			case int:
				if yfield, ok := yfield.(int); ok {
					return option(*eq && xfield == yfield)
				}
				return nil
			case string:
				if yfield, ok := yfield.(string); ok {
					return option(*eq && xfield == yfield)
				}
				return nil
			case *MyStruct:
				if yfield, ok := yfield.(*MyStruct); ok {
					return option(*eq && DeepEqualOriginal(xfield, yfield))
				}
				return nil
			}
			return nil
		})
	return equalValues != nil && *equalValues
}

func TestZipReduceDeepEqualOriginal(t *testing.T) {
	got := DeepEqualOriginal(&MyStruct{1, "2", &MyStruct{3, "4", nil}}, &MyStruct{1, "2", &MyStruct{3, "4", nil}})
	want := true
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func DeepEqual[A any](x, y A) bool {
	eq := ZipReduce(x, y, option(true),
		func(xfield, yfield any, eq *bool) *bool {
			switch xfield := xfield.(type) {
			case int:
				return option(*eq && xfield == yfield.(int))
			case string:
				return option(*eq && xfield == yfield.(string))
			case *MyStruct:
				return option(*eq && DeepEqual(xfield, yfield.(*MyStruct)))
			}
			return nil
		})
	return eq != nil && *eq
}

func TestZipReduceDeepEqual(t *testing.T) {
	got := DeepEqual(&MyStruct{1, "2", &MyStruct{3, "4", nil}}, &MyStruct{1, "2", &MyStruct{3, "4", nil}})
	want := true
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
}
