package comicro

import "reflect"

func Apply(a any, f func(a any) any) any {
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
	return f(a)
}

func Fold[B any](i any, b B, f func(b B, a any) B) B {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return b
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.Elem().Kind() == reflect.Struct {
			for i := 0; i < v.Elem().NumField(); i++ {
				b = f(b, v.Elem().Field(i).Interface())
			}
			return b
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			b = f(b, v.Index(i).Interface())
		}
		return b
	}
	return f(b, i)
}
