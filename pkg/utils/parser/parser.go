package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func SliceParser[T any](separator string, elementParser func(string) (T, error)) func(string) ([]T, error) {
	return func(s string) ([]T, error) {
		if s == "" {
			return []T{}, nil
		}

		parts := strings.Split(s, separator)
		result := make([]T, 0, len(parts))

		for _, part := range parts {
			parsed, err := elementParser(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("failed to parse element '%s': %w", part, err)
			}
			result = append(result, parsed)
		}

		return result, nil
	}
}

var StrParser = func(s string) (string, error) { return s, nil }

var IntParser = func(s string) (int, error) { return strconv.Atoi(s) }

var Int64Parser = func(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }

var BoolParser = func(s string) (bool, error) { return strconv.ParseBool(s) }

var Float64Parser = func(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

var IntSliceParser = SliceParser(",", IntParser)

var BoolSliceParser = SliceParser(",", BoolParser)

var Float64SliceParser = SliceParser(",", Float64Parser)
