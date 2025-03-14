package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// <-----------------------	AI Service Types	------------------------------>

// FieldType represents the supported field types
type FieldType string

const (
	TypeStringArray FieldType = "string_array"
	TypeNumber      FieldType = "number"
	TypeBoolean     FieldType = "boolean"
	TypeObject      FieldType = "object"
	TypeNumberArray FieldType = "number_array"
)

// FieldConfig represents configuration for a single field
type FieldConfig struct {
	Name        string        `json:"name"`
	Type        FieldType     `json:"type"`
	Description string        `json:"description"`
	EnumValues  []string      `json:"enum_values,omitempty"` // For enum type
	Properties  []FieldConfig `json:"properties,omitempty"`  // For object type
}

// Profile represents a predefined set of fields to analyze
type Profile struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Fields      []FieldConfig `json:"fields"`
}

// DynamicChanges is a wrapper around map[string]interface{} to store dynamic fields
type DynamicChanges struct {
	Fields map[string]interface{}
}

// CompareOptions holds configuration for the comparison operation
type CompareOptions struct {
	SystemPrompt string
	Model        string
	Temperature  float64
	MaxTokens    int64
}

// DefaultCompareOptions returns default comparison options
func DefaultCompareOptions() CompareOptions {
	return CompareOptions{
		Model:       openai.ChatModelGPT4oMini,
		Temperature: 0.7,
		MaxTokens:   2048,
	}
}

// <-----------------------	AI Serivce Type Helpers	------------------------------>

// MarshalJSON implements custom JSON marshaling
func (d DynamicChanges) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Fields)
}

// UnmarshalJSON implements custom JSON unmarshaling
func (d *DynamicChanges) UnmarshalJSON(data []byte) error {
	if d.Fields == nil {
		d.Fields = make(map[string]interface{})
	}
	return json.Unmarshal(data, &d.Fields)
}

// Pretty prints the changes in a markdown-like format
func (d *DynamicChanges) Pretty() {
	if d == nil {
		footer := "\nGenerated at: " + time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("No changes detected in the content." + footer)
		return
	}

	caser := cases.Title(language.English)

	markdownContent := "## Changes\n"
	for fieldName, fieldValue := range d.Fields {
		markdownContent += fmt.Sprintf("### %s\n", caser.String(fieldName))
		markdownContent += formatValue(fieldValue, "")
		markdownContent += "\n"
	}

	markdownContent += "\n\nGenerated at: " + time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(markdownContent)
}

// JSON returns the JSON representation of the dynamic changes
func (d *DynamicChanges) JSON(pretty bool) (string, error) {
	if d == nil {
		return "{}", nil
	}

	prettyJSON := make(map[string]interface{})
	caser := cases.Title(language.English)

	for fieldName, fieldValue := range d.Fields {
		prettyJSON[caser.String(fieldName)] = fieldValue
	}

	var data []byte
	var err error
	if !pretty {
		data, err = json.Marshal(prettyJSON)
	} else {
		data, err = json.MarshalIndent(prettyJSON, "", "  ")
	}

	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GenerateFieldSchema generates a schema for a single field based on its type
func GenerateFieldSchema(field FieldConfig) *jsonschema.Schema {
	switch field.Type {
	case TypeStringArray:
		return &jsonschema.Schema{
			Type:        "array",
			Items:       &jsonschema.Schema{Type: "string"},
			Description: field.Description,
		}
	case TypeNumber:
		return &jsonschema.Schema{
			Type:        "number",
			Description: field.Description,
		}
	case TypeBoolean:
		return &jsonschema.Schema{
			Type:        "boolean",
			Description: field.Description,
		}
	case TypeObject:
		properties := jsonschema.NewProperties()
		required := make([]string, 0, len(field.Properties))
		for _, prop := range field.Properties {
			properties.Set(prop.Name, GenerateFieldSchema(prop))
			required = append(required, prop.Name)
		}
		return &jsonschema.Schema{
			Type:        "object",
			Properties:  properties,
			Required:    required,
			Description: field.Description,
		}
	case TypeNumberArray:
		return &jsonschema.Schema{
			Type:        "array",
			Items:       &jsonschema.Schema{Type: "number"},
			Description: field.Description,
		}
	default:
		// Default to string for unknown types
		return &jsonschema.Schema{
			Type:        "string",
			Description: field.Description,
		}
	}
}

// GenerateDynamicSchema generates a JSON schema based on provided field configurations
func GenerateDynamicSchema(fields []FieldConfig) interface{} {
	properties := jsonschema.NewProperties()
	required := make([]string, 0, len(fields))

	for _, field := range fields {
		properties.Set(field.Name, GenerateFieldSchema(field))
		required = append(required, field.Name)
	}

	return &jsonschema.Schema{
		Type:                 "object",
		Properties:           properties,
		Required:             required,
		AdditionalProperties: jsonschema.FalseSchema,
	}
}

// BuildCompetitorSystemPrompt generates the system prompt based on field configurations
func BuildCompetitorSystemPrompt(fields []FieldConfig) string {
	var builder strings.Builder
	builder.WriteString("Analyze the content for changes across these categories:\n# Process\n1. Review content for changes in:\n")

	for i, field := range fields {
		builder.WriteString(fmt.Sprintf("1.%d %s: %s\n", i+1, field.Name, field.Description))
	}

	builder.WriteString(`
2. For each change identified:
2.1 Document the specific change with clear description
2.2 For changes found:
2.2.1 Start with action verbs or clear transition phrases
2.2.2 List each change as a complete, detailed statement
2.2.3 Include relevant context (numbers, timeframes, features)
2.2.4 Separate related but distinct changes into individual items
2.2.5 Structure complex changes into bullet points when needed
2.3 In case no significant changes, DO NOT hallucinate or invent changes`)

	return builder.String()
}

// formatValue formats a field value based on its type for display
func formatValue(value interface{}, indent string) string {
	switch v := value.(type) {
	case []interface{}:
		if len(v) == 0 {
			return "No changes detected"
		}
		var builder strings.Builder
		for _, item := range v {
			builder.WriteString(fmt.Sprintf("%s- %v\n", indent, item))
		}
		return strings.TrimSuffix(builder.String(), "\n")
	case float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case map[string]interface{}:
		var builder strings.Builder
		for k, val := range v {
			formattedVal := formatValue(val, indent+"  ")
			builder.WriteString(fmt.Sprintf("%s%s: %s\n", indent, k, formattedVal))
		}
		return strings.TrimSuffix(builder.String(), "\n")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// normalizeKey cleans and standardizes a key string
func normalizeKey(key string) string {
	// Convert to title case first
	caser := cases.Title(language.English)
	normalized := caser.String(key)

	// Remove leading/trailing whitespace
	normalized = strings.TrimSpace(normalized)

	// Remove common punctuation
	normalized = strings.Map(func(r rune) rune {
		switch r {
		case '.', ',', ':', ';', '-', '_', '/', '\\', '(', ')', '[', ']', '{', '}', '!', '?', '"', '\'':
			return ' '
		default:
			return r
		}
	}, normalized)

	// Replace multiple spaces with single space
	normalized = strings.Join(strings.Fields(normalized), " ")

	// Remove any remaining non-alphanumeric characters except spaces
	normalized = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			return r
		}
		return -1
	}, normalized)

	return normalized
}

// MergeDynamicChanges combines multiple DynamicChanges instances into a single one
// It normalizes keys by removing punctuation, extra spaces, and standardizing case
func MergeDynamicChanges(changes []*DynamicChanges) (*DynamicChanges, error) {
	if len(changes) == 0 {
		return &DynamicChanges{Fields: make(map[string]interface{})}, nil
	}

	merged := &DynamicChanges{
		Fields: make(map[string]interface{}),
	}

	// Process each DynamicChanges instance
	for _, change := range changes {
		if change == nil {
			continue
		}

		for key, value := range change.Fields {
			// Apply comprehensive normalization to the key
			normalizedKey := normalizeKey(key)
			if normalizedKey == "" {
				continue // Skip empty keys after normalization
			}

			// If the key already exists, merge the values
			if existingValue, exists := merged.Fields[normalizedKey]; exists {
				mergedValue, err := mergeValues(existingValue, value)
				if err != nil {
					return nil, fmt.Errorf("error merging values for key '%s': %w", normalizedKey, err)
				}
				merged.Fields[normalizedKey] = mergedValue
			} else {
				// For new keys, just add the value
				merged.Fields[normalizedKey] = value
			}
		}
	}

	return merged, nil
}

// mergeValues combines two values based on their types
func mergeValues(existing, new interface{}) (interface{}, error) {
	// Handle nil cases
	if existing == nil {
		return new, nil
	}
	if new == nil {
		return existing, nil
	}

	switch existingVal := existing.(type) {
	case []interface{}:
		// If existing is an array, append new values
		switch newVal := new.(type) {
		case []interface{}:
			return append(existingVal, newVal...), nil
		default:
			return append(existingVal, newVal), nil
		}

	case map[string]interface{}:
		// If existing is a map, merge with new map
		newMap, ok := new.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot merge map with non-map value")
		}

		result := make(map[string]interface{})
		for k, v := range existingVal {
			normalizedKey := normalizeKey(k)
			if normalizedKey != "" {
				result[normalizedKey] = v
			}
		}

		for k, v := range newMap {
			normalizedKey := normalizeKey(k)
			if normalizedKey == "" {
				continue
			}

			if existingMapVal, exists := result[normalizedKey]; exists {
				mergedMapVal, err := mergeValues(existingMapVal, v)
				if err != nil {
					return nil, err
				}
				result[normalizedKey] = mergedMapVal
			} else {
				result[normalizedKey] = v
			}
		}
		return result, nil

	default:
		// For other types, prefer the new value
		return new, nil
	}
}

func NewEmptyDynamicChanges(fieldNames []string) (*DynamicChanges, error) {
	fields := make(map[string]interface{})
	for _, name := range fieldNames {
		fields[name] = make([]interface{}, 0)
	}
	return &DynamicChanges{Fields: fields}, nil
}
