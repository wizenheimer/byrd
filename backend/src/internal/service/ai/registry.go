// ./src/internal/service/ai/registry.go
package ai

import (
	"fmt"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// DefaultField is the default field for profiles
var DefaultField = "product"

// fields is a list of predefined fields for profiles
var fields = []models.FieldConfig{
	{
		Name:        "branding",
		Type:        models.TypeStringArray,
		Description: "Changes in Visual identity, logos, website design, brand assets",
	},
	{
		Name:        "content",
		Type:        models.TypeStringArray,
		Description: "New posts, any content updates or changes",
	},
	{
		Name:        "customers",
		Type:        models.TypeStringArray,
		Description: "New/removed customers or customer segments",
	},
	{
		Name:        "integration",
		Type:        models.TypeStringArray,
		Description: "New/removed integrations, integration updates",
	},
	{
		Name:        "job_postings",
		Type:        models.TypeStringArray,
		Description: "New job postings, job openings, or changes in hiring",
	},
	{
		Name:        "followers",
		Type:        models.TypeNumber,
		Description: "Increase or decrease in followers",
	},
	{
		Name:        "messaging",
		Type:        models.TypeStringArray,
		Description: "Changes in marketing messages and positioning",
	},
	{
		Name:        "product",
		Type:        models.TypeStringArray,
		Description: "List of product updates, new features, removals, or modifications",
	},
	{
		Name:        "pricing",
		Type:        models.TypeStringArray,
		Description: "Changes in pricing, new plans, promotions or discounts",
	},
	{
		Name:        "partnerships",
		Type:        models.TypeStringArray,
		Description: "New/removed partnerships, collaborations, alliances",
	},
	{
		Name:        "roadmap",
		Type:        models.TypeStringArray,
		Description: "Changes in product roadmap, upcoming features",
	},
	{
		Name:        "testimonials",
		Type:        models.TypeStringArray,
		Description: "New/removed customer feedback, reviews, testimonials or case studies",
	},
	{
		Name:        "targeting",
		Type:        models.TypeStringArray,
		Description: "Changes in target audience, personas or customer segments",
	},
	{
		Name:        "twitter_content",
		Type:        models.TypeStringArray,
		Description: "New tweets, or any content related updates on twitter profile",
	},
	{
		Name:        "facebook_content",
		Type:        models.TypeStringArray,
		Description: "New posts, or any content related updates on facebook profile or page",
	},
	{
		Name:        "instagram_content",
		Type:        models.TypeStringArray,
		Description: "New posts, or any content related updates on instagram profile or page",
	},
	{
		Name:        "content_updates",
		Type:        models.TypeStringArray,
		Description: "New posts, any content related updates, or changes in blogs, case studies, whitepapers, etc",
	},
	{
		Name:        "youtube_content",
		Type:        models.TypeStringArray,
		Description: "New videos, or any content related updates on youtube channel",
	},
	{
		Name:        "podcast_content",
		Type:        models.TypeStringArray,
		Description: "New episodes, or any content related updates on podcast",
	},
	{
		Name:        "subscribers",
		Type:        models.TypeNumber,
		Description: "Increase or Decrease in subscribers or followers",
	},
	{
		Name:        "api_changes",
		Type:        models.TypeStringArray,
		Description: "Changes in API, new endpoints, updates, or removals",
	},
	{
		Name:        "security_updates",
		Type:        models.TypeStringArray,
		Description: "New/removed security features, privacy updates, compliance changes",
	},
	{
		Name:        "compliance_updates",
		Type:        models.TypeStringArray,
		Description: "New/removed compliance features, certifications, or standards",
	},
	{
		Name:        "sdk_updates",
		Type:        models.TypeStringArray,
		Description: "New/removed SDKs, libraries, developer tools, updates",
	},
}

// fieldMap is a map of field name to field config
// This is used to quickly lookup field config by name
var AvailableFields = make(map[string]models.FieldConfig)

// Sanitize removes duplicates from a profile.
// It also validates the profile fields.
// If a field is not valid, it returns an error.
// If a field is duplicated, it removes the duplicates.
func Sanitize(profile models.DiffProfile) (models.DiffProfile, error) {
	// Use a map to track unique fields
	fields := make(map[string]bool)
	var result models.DiffProfile

	// First validate all fields
	for _, f := range profile {
		if _, exists := AvailableFields[f]; !exists {
			return nil, fmt.Errorf("field %s is not a valid field", f)
		}
	}

	// Then deduplicate
	for _, f := range profile {
		if !fields[f] {
			fields[f] = true
			result = append(result, f)
		}
	}

	return result, nil
}

// FieldRegistry contains all available predefined fields
type FieldRegistry struct {
	fields map[string]models.FieldConfig
}

// NewFieldRegistry creates a new field registry with predefined fields
func NewFieldRegistry() *FieldRegistry {
	registry := &FieldRegistry{
		fields: make(map[string]models.FieldConfig),
	}

	// Register predefined fields
	registry.registerPredefinedFields()
	return registry
}

func (r *FieldRegistry) registerPredefinedFields() {
	predefinedFields := fields
	for _, field := range predefinedFields {
		r.fields[field.Name] = field
	}
}

func (r *FieldRegistry) GetField(name string, fallback bool) (models.FieldConfig, error) {
	field, exists := r.fields[name]
	if !exists {
		field = models.FieldConfig{
			Name:        name,
			Type:        models.TypeStringArray,
			Description: "Updates or changes in " + name,
		}
		if fallback {
			return field, nil
		}
		return field, ErrProfileFieldNotFound
	}
	return field, nil
}

func (r *FieldRegistry) ListAvailableFields() []string {
	fields := make([]string, 0, len(r.fields))
	for name := range r.fields {
		fields = append(fields, name)
	}
	return fields
}
