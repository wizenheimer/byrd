package utils

// To returns a pointer to the given value
func ToPtr[T any](v T) *T {
	return &v
}

func FromPtr[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
