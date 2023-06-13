package comicro

import (
	"reflect"
)

// Map maps a function over a the elements of a slice, the fields of a struct or the values of a map.
// For example:
//   - Map([]{1,2,3}, func(x int) int { return x + 1 }) = []{2,3,4}
//   - Map([]{1,2,3}, func(x int) string { return fmt.Sprintf("%d", x) }) = []{"1","2","3"}
func Map(a any, f func(a any) any) any {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr && v.IsNil() {
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
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return a
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.Elem().Kind() == reflect.Struct {
			r := reflect.New(v.Elem().Type())
			for i := 0; i < v.Elem().NumField(); i++ {
				fa := f(v.Elem().Field(i).Interface())
				r.Elem().Field(i).Set(reflect.ValueOf(fa))
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
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return b
	}
	if u, ui := GetUnionValue(a); ui != -1 {
		return fold(u, b, f)
	}
	return fold(a, b, f)
}

func fold[B any](a any, b B, f func(b B, a any) B) B {
	ra := reflect.ValueOf(a)
	if ra.Kind() == reflect.Ptr && ra.IsNil() {
		return b
	}
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
func Any(i any, pred func(a any) bool) bool {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return false
	}
	if u, ui := GetUnionValue(i); ui != -1 {
		return anyOf(u, pred)
	}
	return anyOf(i, pred)
}

func anyOf(a any, pred func(a any) bool) bool {
	ra := reflect.ValueOf(a)
	if ra.Kind() == reflect.Ptr && ra.IsNil() {
		return false
	}
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
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return a, -1
	}
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
		if !v.Field(i).IsNil() {
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
