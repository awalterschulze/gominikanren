package comicro

// Fmap returns a list where each element of the input list has been morphed by the input function.
func fmap[A, B any](list []A, f func(A) B) []B {
	out := make([]B, len(list))
	for i, elem := range list {
		out[i] = f(elem)
	}
	return out
}

func keys[A comparable, B any](m map[A]B) []A {
	out := make([]A, len(m))
	i := 0
	for k := range m {
		out[i] = k
		i++
	}
	return out
}

func tuple3[A, B, C any](a A, b B, c C) func() (A, B, C) {
	return func() (A, B, C) {
		return a, b, c
	}
}

func tuple4[A, B, C, D any](a A, b B, c C, d D) func() (A, B, C, D) {
	return func() (A, B, C, D) {
		return a, b, c, d
	}
}
