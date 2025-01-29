package models

type DispatchReportRequest struct {
	Emails []string `json:"emails" validate:"required,dive,email"`
}
