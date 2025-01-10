// ./src/internal/service/ai/profiles.go
package ai

import (
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

var DefaultUpdates = models.Profile{
	Name:        "competitor_updates",
	Description: "Changes in competitor offerings, features, pricing, integrations",
	Fields: []models.FieldConfig{
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
			Name:        "integration",
			Type:        models.TypeStringArray,
			Description: "New/removed integrations, integration updates",
		},
		{
			Name:        "partnerships",
			Type:        models.TypeStringArray,
			Description: "New/removed partnerships, collaborations, alliances",
		},
		{
			Name:        "customers",
			Type:        models.TypeStringArray,
			Description: "New/removed customers or customer segments",
		},
		{
			Name:        "testimonials",
			Type:        models.TypeStringArray,
			Description: "New/removed customer feedback, reviews, testimonials or case studies",
		},
		{
			Name:        "messaging",
			Type:        models.TypeStringArray,
			Description: "Changes in marketing messages and positioning",
		},
	},
}

// ProductUpdates represents changes in product features, updates, new features
var ProductUpdates = models.Profile{
	Name:        "product_updates",
	Description: "changes in product features, updates, new features",
	Fields: []models.FieldConfig{
		{
			Name:        "product",
			Type:        models.TypeStringArray,
			Description: "List of product updates, new features, removals, or modifications",
		},
		{
			Name:        "roadmap",
			Type:        models.TypeStringArray,
			Description: "Changes in product roadmap, upcoming features",
		},
	},
}

// CustomerUpdates represents changes in customer base, customer feedback, reviews
var CustomerUpdates = models.Profile{
	Name:        "customer_updates",
	Description: "changes in customer base, customer feedback, reviews",
	Fields: []models.FieldConfig{
		{
			Name:        "testimonials",
			Type:        models.TypeStringArray,
			Description: "New/removed customer feedback, reviews, testimonials or case studies",
		},
		{
			Name:        "customers",
			Type:        models.TypeStringArray,
			Description: "New/removed customers logos or segments",
		},
	},
}

// PartnershipUpdates represents changes in partnerships, collaborations, alliances
var PartnershipUpdates = models.Profile{
	Name:        "partnership_updates",
	Description: "changes in partnerships, collaborations, alliances",
	Fields: []models.FieldConfig{
		{
			Name:        "partnerships",
			Type:        models.TypeStringArray,
			Description: "New/removed partnerships, collaborations, alliances or affiliations",
		},
		{
			Name:        "integrations",
			Type:        models.TypeStringArray,
			Description: "New/removed integrations, integration updates",
		},
	},
}

// PricingUpdates represents changes in pricing, new plans, promotions or discounts
var PricingUpdates = models.Profile{
	Name:        "pricing_updates",
	Description: "changes in pricing, new plans, promotions or discounts",
	Fields: []models.FieldConfig{
		{
			Name:        "pricing",
			Type:        models.TypeStringArray,
			Description: "Changes in pricing, new plans, promotions or discounts",
		},
	},
}

// PositioningUpdates represents changes in positioning, messaging, brand identity
var PositioningUpdates = models.Profile{
	Name:        "positioning_updates",
	Description: "changes in positioning, messaging, brand identity",
	Fields: []models.FieldConfig{
		{
			Name:        "targeting",
			Type:        models.TypeStringArray,
			Description: "Changes in target audience, personas or customer segments",
		},
		{
			Name:        "messaging",
			Type:        models.TypeStringArray,
			Description: "Changes in marketing messages and positioning",
		},
		{
			Name:        "branding",
			Type:        models.TypeStringArray,
			Description: "Changes in Visual identity, logos, website design, brand assets",
		},
	},
}

// LinkedInUpdates represents changes in LinkedIn posts, content updates, followers
var LinkedInUpdates = models.Profile{
	Name:        "linkedin_updates",
	Description: "changes in LinkedIn posts, content updates, followers",
	Fields: []models.FieldConfig{
		{
			Name:        "content",
			Type:        models.TypeStringArray,
			Description: "New posts, content updates",
		},
		{
			Name:        "followers",
			Type:        models.TypeNumber,
			Description: "Increase in followers",
		},
	},
}

// TwitterUpdates represents changes in Twitter posts, content updates, followers
var TwitterUpdates = models.Profile{
	Name:        "twitter_updates",
	Description: "changes in Twitter posts, content updates, followers",
	Fields: []models.FieldConfig{
		{
			Name:        "content",
			Type:        models.TypeStringArray,
			Description: "New tweets, content updates",
		},
		{
			Name:        "followers",
			Type:        models.TypeNumber,
			Description: "Increase in followers",
		},
	},
}

// FacebookUpdates represents changes in Facebook posts, content updates, followers
var FacebookUpdates = models.Profile{
	Name:        "facebook_updates",
	Description: "changes in Facebook posts, content updates, followers",
	Fields: []models.FieldConfig{
		{
			Name:        "content",
			Type:        models.TypeStringArray,
			Description: "New posts, content updates",
		},
		{
			Name:        "followers",
			Type:        models.TypeNumber,
			Description: "Increase in followers",
		},
	},
}

// InstagramUpdates represents changes in Instagram posts, content updates, followers
var InstagramUpdates = models.Profile{
	Name:        "instagram_updates",
	Description: "changes in Instagram posts, content updates, followers",
	Fields: []models.FieldConfig{
		{
			Name:        "content",
			Type:        models.TypeStringArray,
			Description: "New posts, content updates",
		},
		{
			Name:        "followers",
			Type:        models.TypeNumber,
			Description: "Increase in followers",
		},
	},
}

// ContentUpdates represents changes in blog posts, content updates, readers
var ContentUpdates = models.Profile{
	Name:        "content_updates",
	Description: "content related updates, changes in blogs, case studies, whitepapers, etc",
	Fields: []models.FieldConfig{
		{
			Name:        "content",
			Type:        models.TypeStringArray,
			Description: "New posts, content updates",
		},
	},
}

// YouTubeUpdates represents changes in YouTube videos, content updates, subscribers
var YouTubeUpdates = models.Profile{
	Name:        "youtube_updates",
	Description: "changes in YouTube videos, content updates, subscribers",
	Fields: []models.FieldConfig{
		{
			Name:        "content",
			Type:        models.TypeStringArray,
			Description: "New videos, content updates",
		},
		{
			Name:        "subscribers",
			Type:        models.TypeNumber,
			Description: "Increase in subscribers",
		},
	},
}

// PodcastUpdates represents changes in podcast episodes, content updates, listeners
var PodcastUpdates = models.Profile{
	Name:        "podcast_updates",
	Description: "changes in podcast episodes, content updates, listeners",
	Fields: []models.FieldConfig{
		{
			Name:        "content",
			Type:        models.TypeStringArray,
			Description: "New episodes, content updates",
		},
		{
			Name:        "listeners",
			Type:        models.TypeNumber,
			Description: "Increase in listeners",
		},
	},
}

// DevUpdates represents changes in API, SDK, libraries, developer tools
var DevUpdates = models.Profile{
	Name:        "dev_updates",
	Description: "changes in API, SDK, libraries, developer tools",
	Fields: []models.FieldConfig{
		{
			Name:        "developer_updates",
			Type:        models.TypeStringArray,
			Description: "New/removed APIs, SDKs, libraries, developer tools",
		},
	},
}

// SecurityUpdates represents changes in security, privacy, compliance, certifications
var SecurityUpdates = models.Profile{
	Name:        "security_updates",
	Description: "changes in security, privacy, compliance, certifications",
	Fields: []models.FieldConfig{
		{
			Name:        "security",
			Type:        models.TypeStringArray,
			Description: "New/removed security features, privacy updates, compliance changes",
		},
	},
}

// CareerPageUpdates represents changes in career page, job postings, hiring updates
var CareerPageUpdates = models.Profile{
	Name:        "career_page_updates",
	Description: "changes in career page, job postings, hiring updates",
	Fields: []models.FieldConfig{
		{
			Name:        "job_postings",
			Type:        models.TypeStringArray,
			Description: "New/removed job postings, hiring updates",
		},
	},
}
