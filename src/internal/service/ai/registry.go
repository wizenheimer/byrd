package ai

import (
	"fmt"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

// FieldRegistry contains all available predefined fields
type FieldRegistry struct {
	fields map[string]core_models.FieldConfig
}

// NewFieldRegistry creates a new field registry with predefined fields
func NewFieldRegistry() *FieldRegistry {
	registry := &FieldRegistry{
		fields: make(map[string]core_models.FieldConfig),
	}

	// Register predefined fields
	registry.registerPredefinedFields()
	return registry
}

func (r *FieldRegistry) registerPredefinedFields() {
	predefinedFields := []core_models.FieldConfig{
		{
			Name:        "branding",
			Type:        core_models.TypeStringArray,
			Description: "Changes in Visual identity, logos, website design, brand assets",
		},
		{
			Name:        "content",
			Type:        core_models.TypeStringArray,
			Description: "New posts, any content updates or changes",
		},
		{
			Name:        "customers",
			Type:        core_models.TypeStringArray,
			Description: "New/removed customers or customer segments",
		},
		{
			Name:        "integration",
			Type:        core_models.TypeStringArray,
			Description: "New/removed integrations, integration updates",
		},
		{
			Name:        "job_postings",
			Type:        core_models.TypeStringArray,
			Description: "New job postings, job openings, or changes in hiring",
		},
		{
			Name:        "followers",
			Type:        core_models.TypeNumber,
			Description: "Increase or decrease in followers",
		},
		{
			Name:        "messaging",
			Type:        core_models.TypeStringArray,
			Description: "Changes in marketing messages and positioning",
		},
		{
			Name:        "product",
			Type:        core_models.TypeStringArray,
			Description: "List of product updates, new features, removals, or modifications",
		},
		{
			Name:        "pricing",
			Type:        core_models.TypeStringArray,
			Description: "Changes in pricing, new plans, promotions or discounts",
		},
		{
			Name:        "partnerships",
			Type:        core_models.TypeStringArray,
			Description: "New/removed partnerships, collaborations, alliances",
		},
		{
			Name:        "roadmap",
			Type:        core_models.TypeStringArray,
			Description: "Changes in product roadmap, upcoming features",
		},
		{
			Name:        "testimonials",
			Type:        core_models.TypeStringArray,
			Description: "New/removed customer feedback, reviews, testimonials or case studies",
		},
		{
			Name:        "targeting",
			Type:        core_models.TypeStringArray,
			Description: "Changes in target audience, personas or customer segments",
		},
		{
			Name:        "twitter_content",
			Type:        core_models.TypeStringArray,
			Description: "New tweets, or any content related updates on twitter profile",
		},
		{
			Name:        "facebook_content",
			Type:        core_models.TypeStringArray,
			Description: "New posts, or any content related updates on facebook profile or page",
		},
		{
			Name:        "instagram_content",
			Type:        core_models.TypeStringArray,
			Description: "New posts, or any content related updates on instagram profile or page",
		},
		{
			Name:        "content_updates",
			Type:        core_models.TypeStringArray,
			Description: "New posts, any content related updates, or changes in blogs, case studies, whitepapers, etc",
		},
		{
			Name:        "youtube_content",
			Type:        core_models.TypeStringArray,
			Description: "New videos, or any content related updates on youtube channel",
		},
		{
			Name:        "podcast_content",
			Type:        core_models.TypeStringArray,
			Description: "New episodes, or any content related updates on podcast",
		},
		{
			Name:        "subscribers",
			Type:        core_models.TypeNumber,
			Description: "Increase or Decrease in subscribers or followers",
		},
		{
			Name:        "api_changes",
			Type:        core_models.TypeStringArray,
			Description: "Changes in API, new endpoints, updates, or removals",
		},
		{
			Name:        "security_updates",
			Type:        core_models.TypeStringArray,
			Description: "New/removed security features, privacy updates, compliance changes",
		},
		{
			Name:        "compliance_updates",
			Type:        core_models.TypeStringArray,
			Description: "New/removed compliance features, certifications, or standards",
		},
		{
			Name:        "sdk_updates",
			Type:        core_models.TypeStringArray,
			Description: "New/removed SDKs, libraries, developer tools, updates",
		},
	}

	for _, field := range predefinedFields {
		r.fields[field.Name] = field
	}
}

func (r *FieldRegistry) GetField(name string, fallback bool) (core_models.FieldConfig, error) {
	field, exists := r.fields[name]
	if !exists {
		field = core_models.FieldConfig{
			Name:        name,
			Type:        core_models.TypeStringArray,
			Description: "Updates or changes in " + name,
		}
		if fallback {
			return field, nil
		}
		return field, fmt.Errorf("field %s not found", name)
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
