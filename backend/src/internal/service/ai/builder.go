package ai

import (
	"encoding/json"

	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/err"
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
func (pb *ProfileBuilder) BuildProfileFromJSON(jsonData string) (models.Profile, err.Error) {
	var request ProfileRequest
	pErr := err.New()
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		pErr.Add(svc.ErrProfileParsing, map[string]any{
			"jsonData": jsonData,
		})
		return models.Profile{}, pErr // TODO: make it non-fatal
	}

	return pb.BuildProfile(request, false)
}

// BuildProfile builds a profile from a request
func (pb *ProfileBuilder) BuildProfile(request ProfileRequest, fallback bool) (models.Profile, err.Error) {
	pErr := err.New()
	if request.Name == "" {
		pErr.Add(svc.ErrProfileNameMissing, map[string]any{
			"request":  request,
			"fallback": fallback,
		})
		return DefaultUpdates, pErr // TODO: make it non-fatal
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
			pErr.Add(svc.ErrProfileFieldNotFound, map[string]any{
				"fieldName": fieldName,
				"fallback":  fallback,
			})
			return models.Profile{}, pErr // TODO: make it non-fatal
		}
		fields = append(fields, field)
	}

	return models.Profile{
		Name:        request.Name,
		Description: request.Description,
		Fields:      fields,
	}, nil
}
