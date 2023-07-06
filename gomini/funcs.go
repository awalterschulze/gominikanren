package gomini

// Fmap returns a list where each element of the input list has been morphed by the input function.
func fmap[A, B any](xs []A, f func(A) B) []B {
	ys := make([]B, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func anyOf[A any](xs []A, f func(A) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}
	return false
}

func zip[A, B any](xs []A, ys []B) []func() (A, B) {
	zs := make([]func() (A, B), len(xs))
	for i := 0; i < len(xs); i++ {
		zs[i] = func() (A, B) { return xs[i], ys[i] }
	}
	return zs
}

func zipReduce[A, B any](xs, ys []A, innit B, f func(x, y A, acc B) (B, bool)) (B, bool) {
	if len(xs) != len(ys) {
		return innit, false
	}
	var ok bool
	b := innit
	for i := 0; i < len(xs); i++ {
		b, ok = f(xs[i], ys[i], b)
		if !ok {
			return b, false
		}
	}
	return b, true
}

func keys[A comparable, B any](m map[A]B) []A {
	ks := make([]A, len(m))
	i := 0
	for k := range m {
		ks[i] = k
		i++
	}
	return ks
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

type Stringer interface {
	String() string
}
