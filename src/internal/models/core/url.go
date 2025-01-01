package models

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID        *uuid.UUID `json:"id"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `json:"created_at"`
}

type URLBatch struct {
	URLs     []URL      `json:"urls"`
	HasMore  bool       `json:"has_more"`
	LastSeen *uuid.UUID `json:"last_seen,omitempty"`
}

type URLInput struct {
	URL string `json:"url"`
}

type URLListInput struct {
	BatchSize  int        `json:"batch_size"`
	LastSeenID *uuid.UUID `json:"last_seen_id"`
}
