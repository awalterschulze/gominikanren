package comicro

import (
	"reflect"
)

// Map maps a function over a the elements of a slice, the fields of a struct or the values of a map.
// For example:
//   - Map([]{1,2,3}, func(x int) int { return x + 1 }) = []{2,3,4}
//   - Map([]{1,2,3}, func(x int) string { return fmt.Sprintf("%d", x) }) = []{"1","2","3"}
func Map(a any, f func(a any) any) any {
	if IsNil(a) {
		return a
	}
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Ptr:
		if v.Elem().Kind() == reflect.Struct {
			r := reflect.New(v.Elem().Type())
			for i := 0; i < v.Elem().NumField(); i++ {
				fv := v.Elem().Field(i).Interface()
				b := f(fv)
				r.Elem().Field(i).Set(reflect.ValueOf(b))
			}
			return r.Interface()
		}
	case reflect.Slice:
		r := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			fa := f(v.Index(i).Interface())
			r.Index(i).Set(reflect.ValueOf(fa))
		}
		return r.Interface()
	case reflect.Map:
		r := reflect.MakeMap(v.Type())
		for _, k := range v.MapKeys() {
			fa := f(v.MapIndex(k).Interface())
			r.SetMapIndex(k, reflect.ValueOf(fa))
		}
		return r.Interface()
	}
	return a
}

// Fold folds a function over the elements of a slice or the fields of a struct
// For example:
//   - Fold([]{1,2,3}, 0, func(x, sum int) int { return sum + x }) = 6
//   - Fold([]{1,2,3}, 0, func(x int, s string) string { return s + fmt.Sprintf("%d", x) }) = "123"
func Fold[B any](a any, b B, f func(a any, b B) B) B {
	if IsNil(a) {
		return b
	}
	ra := reflect.ValueOf(a)
	switch ra.Kind() {
	case reflect.Ptr:
		if ra.Elem().Kind() == reflect.Struct {
			for i := 0; i < ra.Elem().NumField(); i++ {
				b = f(ra.Elem().Field(i).Interface(), b)
			}
			return b
		}
	case reflect.Slice:
		for i := 0; i < ra.Len(); i++ {
			b = f(ra.Index(i).Interface(), b)
		}
		return b
	}
	return b
}

// Any returns true if the predicate is true for any of the elements in the slice of fields of a struct.
// For example:
//   - Any([]{1,2,3}, func(x int) bool { return x == 2 }) = true
//   - Any([]{1,2,3}, func(x int) bool { return x == 4 }) = false
func Any(a any, pred func(a any) bool) bool {
	if IsNil(a) {
		return false
	}
	ra := reflect.ValueOf(a)
	switch ra.Kind() {
	case reflect.Ptr:
		if ra.Elem().Kind() == reflect.Struct {
			for i := 0; i < ra.Elem().NumField(); i++ {
				if pred(ra.Elem().Field(i).Interface()) {
					return true
				}
			}
			return false
		}
	case reflect.Slice:
		for i := 0; i < ra.Len(); i++ {
			if pred(ra.Index(i).Interface()) {
				return true
			}
		}
		return false
	}
	return false
}

// ZipFold folds a function over the elements of two slices or the fields of two structs together
// It allows you to shortcircuit the fold by returning false from the function
// For example:
//   - ZipFold([]{1,2,3}, []{3,2,1}, 0, func(x1, x2, sum int) (int, bool) { return sum + x1 + x2, true }) = (12, true)
//   - ZipFold([]{1,2,3}, []{3,2,1}, 0, func(x1, x2, int, s string) (string, bool) { return s + fmt.Sprintf("%d%d", x1, x2), true }) = ("132231", true)
//   - ZipFold([]{1,2,3}, []{3,2,1}, 0, func(x1, x2, sum int) (int, bool) { if x1 == 2 { return 0, false } else { return sum + x1 + x2, true }) = (0, false)
func ZipFold[B any](a1, a2 any, b B, f func(a1, a2 any, b B) (B, bool)) (B, bool) {
	if IsNil(a1) {
		if IsNil(a2) {
			return b, true
		}
		return b, false
	}
	if IsNil(a2) {
		return b, false
	}
	ra1 := reflect.ValueOf(a1)
	ra2 := reflect.ValueOf(a2)
	if ra1.Kind() != ra2.Kind() {
		return b, false
	}
	switch ra1.Kind() {
	case reflect.Ptr:
		if ra1.Elem().Kind() != ra2.Elem().Kind() {
			return b, false
		}
		if ra1.Elem().Kind() == reflect.Struct {
			if ra1.Elem().NumField() != ra2.Elem().NumField() {
				return b, false
			}
			var ok bool
			for i := 0; i < ra1.Elem().NumField(); i++ {
				b, ok = f(ra1.Elem().Field(i).Interface(), ra2.Elem().Field(i).Interface(), b)
				if !ok {
					return b, false
				}
			}
			return b, true
		}
	case reflect.Slice:
		if ra1.Len() != ra2.Len() {
			return b, false
		}
		var ok bool
		for i := 0; i < ra1.Len(); i++ {
			b, ok = f(ra1.Index(i).Interface(), ra2.Index(i).Interface(), b)
			if !ok {
				return b, false
			}
		}
		return b, true
	}
	return b, false
}

func IsNil(a any) bool {
	if a == nil {
		return true
	}
	v := reflect.ValueOf(a)
	if !v.IsValid() {
		return true
	}
	return v.Kind() == reflect.Ptr && v.IsNil()
}
