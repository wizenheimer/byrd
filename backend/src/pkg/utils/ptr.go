// ./src/pkg/utils/ptr.go
package utils

// To returns a pointer to the given value
func ToPtr[T any](v T) *T {
	return &v
}
