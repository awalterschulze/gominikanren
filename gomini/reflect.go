package gomini

import (
	"reflect"
)

// Map maps a function over a the elements of a slice, the fields of a struct or the values of a map.
// For example:
//   - Map([]int{1, 2, 3}, func(x any) any { return x.(int) + 1 }).([]int) = []int{2,3,4}
//   - See Tests for an example of how to map over a struct
//
// The following is not possible, because we can only Map back to the same types:
//   - Map([]{1,2,3}, func(x int) string { return fmt.Sprintf("%d", x) }) = []{"1","2","3"}
func Map(x any, f func(a any) any) any {
	if IsNil(x) {
		return x
	}
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Ptr:
		if v.Elem().Kind() == reflect.Struct {
			r := reflect.New(v.Elem().Type())
			for i := 0; i < v.Elem().NumField(); i++ {
				a := v.Elem().Field(i).Interface()
				b := f(a)
				r.Elem().Field(i).Set(reflect.ValueOf(b))
			}
			return r.Interface()
		}
	case reflect.Slice:
		r := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			a := v.Index(i).Interface()
			b := f(a)
			r.Index(i).Set(reflect.ValueOf(b))
		}
		return r.Interface()
	case reflect.Map:
		r := reflect.MakeMap(v.Type())
		for _, k := range v.MapKeys() {
			a := v.MapIndex(k).Interface()
			b := f(a)
			r.SetMapIndex(k, reflect.ValueOf(b))
		}
		return r.Interface()
	}
	return x
}

// Any returns true if the predicate is true for any of the elements in the slice of fields of a struct.
// For example:
//   - Any([]int{1, 2, 3}, func(x any) bool { return x.(int) == 2 }) = true
//   - Any([]int{1, 2, 3}, func(x any) bool { return x.(int) == 4 }) = false
func Any(x any, pred func(a any) bool) bool {
	if IsNil(x) {
		return false
	}
	ra := reflect.ValueOf(x)
	switch ra.Kind() {
	case reflect.Ptr:
		if ra.Elem().Kind() == reflect.Struct {
			for i := 0; i < ra.Elem().NumField(); i++ {
				a := ra.Elem().Field(i).Interface()
				if pred(a) {
					return true
				}
			}
			return false
		}
	case reflect.Slice:
		for i := 0; i < ra.Len(); i++ {
			a := ra.Index(i).Interface()
			if pred(a) {
				return true
			}
		}
		return false
	}
	return false
}

// ZipReduce reduces a function over the elements of two slices or the fields of two structs together
// It allows you to shortcircuit the reduce by returning false from the function
// For example:
//   - ZipReduce([]int{1, 2, 3}, []int{3, 2, 1}, 0, func(x1, x2 any, sum int) (int, bool) { return sum + x1.(int) + x2.(int), true }) = (12, true)
//   - ZipReduce([]int{1, 2, 3}, []int{3, 2, 1}, "", func(x1, x2 any, s string) (string, bool) { return s + fmt.Sprintf("%d%d", x1, x2), true }) = ("132231", true)
func ZipReduce[B any](x, y any, innit *B, f func(x, y any, acc *B) *B) *B {
	if IsNil(x) {
		if IsNil(y) {
			return innit
		}
		return nil
	}
	if IsNil(y) {
		return nil
	}
	rx := reflect.ValueOf(x)
	ry := reflect.ValueOf(y)
	if rx.Kind() != ry.Kind() {
		return nil
	}
	switch rx.Kind() {
	case reflect.Ptr:
		if rx.Elem().Kind() != ry.Elem().Kind() {
			return nil
		}
		if rx.Elem().Kind() == reflect.Struct {
			if rx.Elem().NumField() != ry.Elem().NumField() {
				return nil
			}
			b := innit
			for i := 0; i < rx.Elem().NumField(); i++ {
				b = f(rx.Elem().Field(i).Interface(), ry.Elem().Field(i).Interface(), b)
				if b == nil {
					return nil
				}
			}
			return b
		}
	case reflect.Slice:
		if rx.Len() != ry.Len() {
			return nil
		}
		b := innit
		for i := 0; i < rx.Len(); i++ {
			b = f(rx.Index(i).Interface(), ry.Index(i).Interface(), b)
			if b == nil {
				return nil
			}
		}
		return b
	}
	return nil
}

func IsNil(x any) bool {
	if x == nil {
		return true
	}
	v := reflect.ValueOf(x)
	if !v.IsValid() {
		return true
	}
	return v.Kind() == reflect.Ptr && v.IsNil()
}
