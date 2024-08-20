package xslice

func FindBy[T comparable](s []T, f func(v T) bool) (T, bool) {
	for _, v := range s {
		if f(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}
