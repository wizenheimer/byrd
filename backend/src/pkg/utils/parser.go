package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
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

func DeduplicateElements[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func Contains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func ParseTextFromHTML(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	var extract func(*html.Node)

	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			// Get text content
			text := strings.TrimSpace(n.Data)
			if text != "" {
				buf.WriteString(text)
				buf.WriteString("\n")
			}
		}

		if n.Type == html.ElementNode {
			// Skip script and style elements
			switch n.Data {
			case "script", "style", "noscript":
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)

	// Clean up the extracted text
	text := buf.String()

	// Remove extra whitespace
	text = strings.Join(strings.Fields(text), " ")

	// Remove multiple newlines
	text = strings.ReplaceAll(text, "\n\n", "\n")

	return text, nil
}
