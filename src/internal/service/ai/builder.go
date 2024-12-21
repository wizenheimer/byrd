package ai

import (
	"encoding/json"
	"fmt"

	"github.com/wizenheimer/iris/src/internal/domain/models"
)

// UserProfileRequest represents the user's request to create a profile
type ProfileRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	FieldNames  []string `json:"fields"`
}

// ProfileBuilder helps build profiles from user requests
type ProfileBuilder struct {
	registry *FieldRegistry
}

// NewProfileBuilder creates a new profile builder
func NewProfileBuilder(registry *FieldRegistry) *ProfileBuilder {
	return &ProfileBuilder{
		registry: registry,
	}
}

// BuildProfileFromJSON builds a profile from a JSON request
func (pb *ProfileBuilder) BuildProfileFromJSON(jsonData string) (models.Profile, error) {
	var request ProfileRequest
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		return models.Profile{}, fmt.Errorf("invalid request JSON: %w", err)
	}

	return pb.BuildProfile(request, false)
}

// BuildProfile builds a profile from a request
func (pb *ProfileBuilder) BuildProfile(request ProfileRequest, fallback bool) (models.Profile, error) {
	if request.Name == "" {
		return DefaultUpdates, fmt.Errorf("profile name is required")
	}

	// Deduplicate field names
	fieldNames := make(map[string]struct{})
	for _, fieldName := range request.FieldNames {
		fieldNames[fieldName] = struct{}{}
	}

	request.FieldNames = make([]string, 0, len(fieldNames))
	for fieldName := range fieldNames {
		request.FieldNames = append(request.FieldNames, fieldName)
	}

	fields := make([]models.FieldConfig, 0, len(request.FieldNames))
	for _, fieldName := range request.FieldNames {
		field, err := pb.registry.GetField(fieldName, fallback)
		if err != nil {
			return models.Profile{}, fmt.Errorf("error getting field '%s': %w", fieldName, err)
		}
		fields = append(fields, field)
	}

	return models.Profile{
		Name:        request.Name,
		Description: request.Description,
		Fields:      fields,
	}, nil
}
