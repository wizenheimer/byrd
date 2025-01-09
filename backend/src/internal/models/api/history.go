package models

import (
	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type PageHistoryUpdateRequest struct {
	PageID      uuid.UUID          `json:"page_id"`
	PageHistory models.PageHistory `json:"page_history"`
}

type BatchPageHistoryUpdate []PageHistoryUpdateRequest
