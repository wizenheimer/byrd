package models

type EmailTemplate string

const (
	DiffReportTemplate  EmailTemplate = "diff-report"
	Trial0Day           EmailTemplate = "trial-0-day"
	Trial3Day           EmailTemplate = "trial-3-day"
	Trial7Day           EmailTemplate = "trial-7-day"
	Trial14Day          EmailTemplate = "trial-14-day"
	SuccessConversion   EmailTemplate = "successful-conversion"
	FailedConversion    EmailTemplate = "failed-conversion"
	WaitlistOffboarding EmailTemplate = "waitlist-offboarding"
	WaitlistOnboarding  EmailTemplate = "waitlist-onboarding"
)

type NotificationRequest struct {
	TemplateID          EmailTemplate `json:"templateId" validate:"required"`
	EmailTemplateParams interface{}   `json:"emailTemplateParams" validate:"required"`
	Emails              []string      `json:"emails" validate:"required,min=1,max=100,dive,email"`
}

type NotificationResults struct {
	EmailNotificationResults EmailNotificationResults `json:"email"`
}

type EmailNotificationResults struct {
	Successful []string `json:"successful"`
	Failed     []string `json:"failed"`
}

type DiffReportEmailParams struct {
	Kind       EmailTemplate       `json:"kind"`
	Competitor string              `json:"competitor"`
	FromDate   string              `json:"fromDate"`
	ToDate     string              `json:"toDate"`
	Data       map[string]Category `json:"data"`
}

type Category struct {
	Summary string              `json:"summary"`
	Changes []string            `json:"changes"`
	URLs    map[string][]string `json:"urls"`
}

type BaseEmailParams struct {
	UserName string `json:"userName,omitempty"`
}

type TrialEmailParams struct {
	BaseEmailParams
	UpgradeLink string `json:"upgradeLink" validate:"required,url"`
}

type WaitlistEmailParams struct {
	BaseEmailParams
	InviteLink string `json:"inviteLink" validate:"required,url"`
}

type EmailParams struct {
	From        string   `json:"from" validate:"required,email"`
	To          []string `json:"to" validate:"required,min=1,dive,email"`
	Subject     string   `json:"subject" validate:"required"`
	HTML        string   `json:"html" validate:"required"`
	Text        string   `json:"text,omitempty"`
	ReplyTo     string   `json:"replyTo,omitempty" validate:"omitempty,email"`
	Attachments []struct {
		Filename string `json:"filename"`
		Content  []byte `json:"content"`
		Type     string `json:"type"`
	} `json:"attachments,omitempty"`
}
