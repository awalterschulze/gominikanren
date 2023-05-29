package sexpr

func tuple2[A, B any](a A, b B) func() (A, B) {
	return func() (A, B) {
		return a, b
	}
}

func tuple3[A, B, C any](a A, b B, c C) func() (A, B, C) {
	return func() (A, B, C) {
		return a, b, c
	}
}
