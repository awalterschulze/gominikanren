package comicro

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestMapInc(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(x any) any { return x.(int) + 1 }).([]int)
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

type A struct {
	A int
	B string
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

func TestFoldSum(t *testing.T) {
	got := Fold([]int{1, 2, 3}, 0, func(x, sum any) any { return sum.(int) + x.(int) }).(int)
	want := 6
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestFoldConcat(t *testing.T) {
	got := Fold([]int{1, 2, 3}, "", func(x any, s any) any { return s.(string) + fmt.Sprintf("%d", x) }).(string)
	want := "123"
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestFoldStruct(t *testing.T) {
	got := Fold(&A{1, "2"}, 0, func(x, sum any) any {
		switch x := x.(type) {
		case int:
			return sum.(int) + x
		case string:
			xi, err := strconv.Atoi(x)
			if err != nil {
				return sum.(int)
			}
			return sum.(int) + xi
		}
		return sum
	}).(int)
	want := 3
	if got != want {
		t.Fatalf("expected %v but got %v", want, got)
	}
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

func TestZipFoldSum(t *testing.T) {
	got, gotok := ZipFold([]int{1, 2, 3}, []int{3, 2, 1}, 0, func(x1, x2 any, sum int) (int, bool) { return sum + x1.(int) + x2.(int), true })
	want, wantok := 12, true
	if got != want || gotok != wantok {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestZipFoldConcat(t *testing.T) {
	got, gotok := ZipFold([]int{1, 2, 3}, []int{3, 2, 1}, "", func(x1, x2 any, s string) (string, bool) { return s + fmt.Sprintf("%d%d", x1, x2), true })
	want, wantok := "132231", true
	if got != want || gotok != wantok {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestZipFoldFalse(t *testing.T) {
	got, gotok := ZipFold([]int{1, 2, 3}, []int{3, 2, 1}, 0,
		func(x1, x2 any, sum int) (int, bool) {
			if x1 == 2 {
				return 0, false
			} else {
				return sum + x1.(int) + x2.(int), true
			}
		})
	want, wantok := 0, false
	if got != want || gotok != wantok {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestZipFoldStruct(t *testing.T) {
	got, gotok := ZipFold(&A{1, "2"}, &A{2, "3"}, 0, func(x1, x2 any, sum int) (int, bool) {
		switch x1 := x1.(type) {
		case int:
			return sum + x1 + x2.(int), true
		case string:
			x1i, err := strconv.Atoi(x1)
			if err != nil {
				return sum, false
			}
			x2i, err := strconv.Atoi(x2.(string))
			if err != nil {
				return sum, false
			}
			return sum + x1i + x2i, true
		}
		return sum, false
	})
	want, wantok := 8, true
	if got != want || gotok != wantok {
		t.Fatalf("expected %v but got %v", want, got)
	}
}

func TestZipFoldDeepEqual(t *testing.T) {
	got, gotok := ZipFold(&A{1, "2"}, &A{1, "2"}, true, func(x1, x2 any, eq bool) (bool, bool) {
		switch x1 := x1.(type) {
		case int:
			return eq && x1 == x2.(int), true
		case string:
			return eq && x1 == x2.(string), true
		}
		return eq, false
	})
	want, wantok := true, true
	if got != want || gotok != wantok {
		t.Fatalf("expected %v but got %v", want, got)
	}
}
