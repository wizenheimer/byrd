package report

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/email/template"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/internal/repository/report"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type reportService struct {
	// logger is the logger used by the service.
	logger *logger.Logger

	// aiService is the AI service used by the service.
	aiService ai.AIService

	// repo
	repo report.ReportRepository

	// library
	library template.TemplateLibrary

	// emailChannel
	emailClient email.EmailClient

	// errorRecorder
	errorRecorder *recorder.ErrorRecorder
}

// NewReportService creates a new report service.
func NewReportService(
	aiService ai.AIService,
	emailClient email.EmailClient,
	library template.TemplateLibrary,
	repo report.ReportRepository,
	logger *logger.Logger,
	errorRecorder *recorder.ErrorRecorder,
) (ReportService, error) {
	rs := reportService{
		logger: logger.WithFields(map[string]any{
			"service": "report",
		}),
		aiService:     aiService,
		emailClient:   emailClient,
		errorRecorder: errorRecorder,
		library:       library,
		repo:          repo,
	}
	return &rs, nil
}

// Get returns the report with the given ID.
func (s *reportService) Get(ctx context.Context, reportID uuid.UUID) (*models.Report, error) {
	s.logger.Debug("Getting report", zap.Any("reportID", reportID))
	return s.repo.Get(ctx, reportID)
}

// GetLatest returns the latest report for the given workspace and competitor
func (s *reportService) GetLatest(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Report, error) {
	s.logger.Debug("Getting latest report", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitorID))
	return s.repo.GetLatest(ctx, workspaceID, competitorID)
}

// List returns a list of reports for the given workspace and competitor
func (s *reportService) List(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error) {
	s.logger.Debug("Listing reports", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitorID), zap.Any("limit", limit), zap.Any("offset", offset))
	return s.repo.List(ctx, workspaceID, competitorID, limit, offset)
}

// Create creates a new report for the given workspace and competitor
func (s *reportService) Create(ctx context.Context, workspaceID, competitorID uuid.UUID, history []models.PageHistory) (*models.Report, error) {
	// Calculate the period boundaries
	// Check for reports in the last week
	oneWeekAgo := time.Now().UTC().AddDate(0, 0, -7)

	existingReport, exists, err := s.repo.GetForPeriod(ctx, workspaceID, competitorID, oneWeekAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing report: %w", err)
	}

	// If report exists, return it
	if existingReport != nil && exists {
		s.logger.Debug("Found existing report for period",
			zap.Any("reportID", existingReport.ID))
		return existingReport, nil
	}

	// No existing report found, create a new one
	s.logger.Debug("Creating new report for period",
		zap.Any("workspaceID", workspaceID),
		zap.Any("competitorID", competitorID),
		zap.Any("history", history))

	changeList := make([]*models.DynamicChanges, 0)
	for _, pageHistory := range history {
		changeList = append(changeList, &pageHistory.DiffContent)
	}

	changes, err := s.aiService.SummarizeChanges(ctx, changeList)
	if err != nil {

		return nil, err
	}

	report := models.NewReport(workspaceID, competitorID, changes)
	if err := s.repo.Set(ctx, report); err != nil {
		s.errorRecorder.RecordError(ctx, err, zap.Any("competitorID", competitorID), zap.Any("workspaceID", workspaceID))
		return report, err
	}

	return report, nil
}

// Dispatch send the report to it's subscribers.
func (s *reportService) Dispatch(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string, subscriberEmails []string) error {
	report, err := s.GetLatest(ctx, workspaceID, competitorID)
	if err != nil {
		return err
	}

	tmp, err := s.library.GetTemplate(template.WeeklyRoundupTemplate)
	if err != nil {
		return err
	}

	// ERROR: template.SectionedTemplate is not a typecompilerNotAType
	sectionedTemplate, ok := tmp.(*template.SectionedTemplate)
	if !ok {
		return errors.New("failed to assert template to SectionedTemplate")
	}

	// Override template with report data
	sectionedTemplate.Competitor = competitorName
	sectionedTemplate.GeneratedAt = time.Now()
	sectionedTemplate.FromDate = report.Time.AddDate(0, 0, -7) // Assuming weekly report
	sectionedTemplate.ToDate = report.Time

	// Create sections map
	sectionedTemplate.Sections = make(map[string]template.Section)

	// Map each CategoryChange to a Section
	for _, change := range report.Changes {
		bullets := make([]template.BulletPoint, len(change.Changes))

		// Convert changes to bullet points
		for i, changeText := range change.Changes {
			bullets[i] = template.BulletPoint{
				Text: changeText,
				// TODO: LinkURL can be set later
			}
		}

		sectionedTemplate.Sections[change.Category] = template.Section{
			Title:   change.Category,
			Summary: change.Summary,
			Bullets: bullets,
		}
	}

	// Render the template
	htmlContent, err := sectionedTemplate.RenderHTML()
	if err != nil {
		return err
	}

	email := models.Email{
		To:           subscriberEmails,
		EmailSubject: "Weekly Roundup for " + competitorName,
		EmailContent: htmlContent,
		EmailFormat:  models.EmailFormatHTML,
	}

	go s.sendEmail(email)

	return nil
}

func (s *reportService) sendEmail(email models.Email) {
	// Create a context with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Important to avoid context leak

	if err := s.emailClient.Send(ctx, email); err != nil {
		s.errorRecorder.RecordError(ctx, err, zap.Any("subscriberEmails", email.To), zap.Any("emailSubject", email.EmailSubject))
	}
}
