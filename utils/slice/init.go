package slice

func Of[T any](n int, v T) []T {
	s := make([]T, n)
	for i := 0; i < n; i++ {
		s[i] = v
	}
	return s
}
