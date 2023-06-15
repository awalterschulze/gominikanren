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
	u, ui := GetUnionValue(a)
	if ui != -1 {
		fu := mapOverAny(u, f)
		return SetUnionValue(a, ui, fu)
	}
	return mapOverAny(a, f)
}

func mapOverAny(a any, f func(a any) any) any {
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
func Fold[B any](a any, b B, f func(b B, a any) B) B {
	if IsNil(a) {
		return b
	}
	if u, ui := GetUnionValue(a); ui != -1 {
		return fold(u, b, f)
	}
	return fold(a, b, f)
}

func fold[B any](a any, b B, f func(b B, a any) B) B {
	if IsNil(a) {
		return b
	}
	ra := reflect.ValueOf(a)
	switch ra.Kind() {
	case reflect.Ptr:
		if ra.Elem().Kind() == reflect.Struct {
			for i := 0; i < ra.Elem().NumField(); i++ {
				b = f(b, ra.Elem().Field(i).Interface())
			}
			return b
		}
	case reflect.Slice:
		for i := 0; i < ra.Len(); i++ {
			b = f(b, ra.Index(i).Interface())
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
	if u, ui := GetUnionValue(a); ui != -1 {
		return anyOf(u, pred)
	}
	return anyOf(a, pred)
}

func anyOf(a any, pred func(a any) bool) bool {
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

func GetUnionValue(a any) (any, int) {
	if IsNil(a) {
		return a, -1
	}
	v := reflect.ValueOf(a)
	typ := reflect.TypeOf(a)
	if v.Kind() != reflect.Ptr {
		return a, -1
	}
	v = v.Elem()
	typ = typ.Elem()
	if v.Kind() != reflect.Struct {
		return a, -1
	}
	if v.NumField() < 2 {
		return a, -1
	}
	nonnil := -1
	for i := 0; i < v.NumField(); i++ {
		if typ.Field(i).Type.Kind() != reflect.Ptr {
			return a, -1
		}
		if !IsNil(v.Field(i).Interface()) {
			if nonnil != -1 {
				return a, -1
			}
			nonnil = i
		}
	}
	return v.Field(nonnil).Interface(), nonnil
}

func SetUnionValue(a any, i int, b any) any {
	v := reflect.ValueOf(a)
	r := reflect.New(v.Type().Elem())
	r.Elem().Field(i).Set(reflect.ValueOf(b))
	return r.Interface()
}

func IsContainer(a any) bool {
	if IsNil(a) {
		return false
	}
	ra := reflect.ValueOf(a)
	if ra.Kind() == reflect.Ptr {
		ra = ra.Elem()
	}
	switch ra.Kind() {
	case reflect.Slice, reflect.Struct:
		return true
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
	if u1, ui1 := GetUnionValue(a1); ui1 != -1 {
		if u2, ui2 := GetUnionValue(a2); ui2 != -1 {
			return zipFold(u1, u2, b, f)
		}
		return b, false
	}
	return zipFold(a1, a2, b, f)
}

func zipFold[B any](a1, a2 any, b B, f func(a1, a2 any, b B) (B, bool)) (B, bool) {
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
	return b, true
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
