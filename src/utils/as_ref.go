package utils

func AsRef[T any](t T) *T {
	return &t
}
