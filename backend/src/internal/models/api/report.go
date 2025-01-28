package models

import models "github.com/wizenheimer/byrd/src/internal/models/core"

type CreateReportRequest struct {
	History []models.PageHistory `json:"history" validate:"required"`
}

type DispatchReportRequest struct {
	CompetitorName   string   `json:"competitor_name" validate:"required"`
	SubscriberEmails []string `json:"subscriber_emails" validate:"required,dive,email"`
}
