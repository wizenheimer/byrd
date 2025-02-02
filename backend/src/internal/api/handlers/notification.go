package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	logger      *logger.Logger
	emailClient email.EmailClient
	library     template.TemplateLibrary
}

func NewNotificationHandler(
	logger *logger.Logger,
	emailClient email.EmailClient,
	library template.TemplateLibrary,
) (*NotificationHandler, error) {
	return &NotificationHandler{
		logger:      logger,
		emailClient: emailClient,
		library:     library,
	}, nil
}

type CommonTemplateRequest struct {
	PreviewText string                 `json:"preview_text"`
	Logo        string                 `json:"logo,omitempty"`
	Title       string                 `json:"title"`
	Subtitle    string                 `json:"subtitle,omitempty"`
	Body        []string               `json:"body"`
	BulletTitle string                 `json:"bullet_title,omitempty"`
	Bullets     []string               `json:"bullets,omitempty"`
	CTA         *template.CallToAction `json:"cta,omitempty"`
	ClosingText string                 `json:"closing_text,omitempty"`
	Footer      template.Footer        `json:"footer"`
}

type SectionedTemplateRequest struct {
	Competitor string                      `json:"competitor"`
	FromDate   time.Time                   `json:"from_date"`
	ToDate     time.Time                   `json:"to_date"`
	Summary    string                      `json:"summary,omitempty"`
	Sections   map[string]template.Section `json:"sections"`
}

// SendNotification sends a notification to the user
func (nh *NotificationHandler) SendNotification(c *fiber.Ctx) error {
	nh.logger.Debug("sending notification")
	userEmail := c.Query("email", "")
	if userEmail == "" {
		return sendErrorResponse(c, nh.logger, fiber.StatusBadRequest, "Email not provided", "Email not provided")
	}

	templateTypeString := strings.ToLower(c.Query("template_type", "common"))
	var emailHTML string
	switch templateTypeString {
	case "common":
		var req CommonTemplateRequest
		err := c.BodyParser(&req)
		if err != nil {
			return sendErrorResponse(c, nh.logger, fiber.StatusBadRequest, "Couldn't parse the common template request", err.Error())
		}
		tmpl := &template.CommonTemplate{
			PreviewText: req.PreviewText,
			Logo:        req.Logo,
			Title:       req.Title,
			Subtitle:    req.Subtitle,
			Body:        req.Body,
			BulletTitle: req.BulletTitle,
			Bullets:     req.Bullets,
			CTA:         req.CTA,
			ClosingText: req.ClosingText,
			Footer:      req.Footer,
			GeneratedAt: time.Now(),
		}
		emailHTML, err = tmpl.RenderHTML()
		if err != nil {
			return sendErrorResponse(c, nh.logger, fiber.StatusInternalServerError, "Couldn't render the email template", err.Error())
		}

	case "sectioned":
		var req SectionedTemplateRequest
		err := c.BodyParser(&req)
		if err != nil {
			return sendErrorResponse(c, nh.logger, fiber.StatusBadRequest, "Couldn't parse the sectioned template request", err.Error())
		}

		tmpl := &template.SectionedTemplate{
			Competitor:  req.Competitor,
			FromDate:    req.FromDate,
			ToDate:      req.ToDate,
			Summary:     req.Summary,
			Sections:    req.Sections,
			GeneratedAt: time.Now(),
		}

		emailHTML, err = tmpl.RenderHTML()
		if err != nil {
			return sendErrorResponse(c, nh.logger, fiber.StatusInternalServerError, "Couldn't render the email template", err.Error())
		}
	}

	go func() {
		nh.logger.Debug("sending email", zap.String("email", userEmail))
		email := models.Email{
			To:           []string{userEmail},
			EmailContent: emailHTML,
			EmailSubject: "Test",
			EmailFormat:  models.EmailFormatHTML,
		}
		// TODO: fix this
		if err := nh.emailClient.Send(context.Background(), email); err != nil {
			nh.logger.Error("Error sending email", zap.Error(err))
		}
	}()

	return sendDataResponse(c, fiber.StatusOK, "Email sent successfully", map[string]string{
		"email": userEmail,
		"html":  emailHTML,
	})
}
