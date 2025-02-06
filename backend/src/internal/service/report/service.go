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

const MaxReportQueryLimit = 25

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
	report, err := s.repo.Get(ctx, reportID)
	if err != nil {
		return nil, err
	}

	if report == nil {
		return nil, fmt.Errorf("report with ID %s not found", reportID)
	}

	return report, nil
}

// GetLatest returns the latest report for the given workspace and competitor
func (s *reportService) GetLatest(ctx context.Context, workspaceID, competitorID uuid.UUID) (*models.Report, error) {
	report, err := s.repo.GetLatest(ctx, workspaceID, competitorID)
	if err != nil {
		return nil, err
	}

	if report == nil {
		return nil, fmt.Errorf("report not found for workspace %s and competitor %s", workspaceID, competitorID)
	}

	return report, nil
}

// GetContent returns the content of the report with the given ID.
func (s *reportService) GetContent(ctx context.Context, reportURI string) (string, error) {
	reportContent, err := s.repo.GetReportContent(ctx, reportURI)
	if err != nil {
		return "", err
	}

	return reportContent, nil
}

// List returns a list of reports for the given workspace and competitor
func (s *reportService) List(ctx context.Context, workspaceID, competitorID uuid.UUID, limit, offset *int) ([]models.Report, bool, error) {
	// Check if limit and offset are valid
	if limit != nil {
		if *limit < 0 {
			return nil, false, errors.New("limit cannot be negative")
		} else if *limit > MaxReportQueryLimit {
			return nil, false, fmt.Errorf("limit cannot exceed %d", MaxReportQueryLimit)
		}
	}

	if offset != nil && *offset < 0 {
		return nil, false, errors.New("offset cannot be negative")
	}

	// List reports from the repository
	return s.repo.List(ctx, workspaceID, competitorID, limit, offset)
}

// Create creates a new report for the given workspace and competitor
func (s *reportService) Create(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string, history []models.PageHistory) (*models.Report, error) {
	// Calculate the period boundaries
	// Check for reports in the last week
	oneWeekAgo := time.Now().UTC().AddDate(0, 0, -7)

	existingReport, exists, err := s.repo.GetForPeriod(ctx, workspaceID, competitorID, oneWeekAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing report: %w", err)
	}

	// If report exists, return it
	if existingReport != nil && exists {
		s.logger.Debug("existing report found for the time period, using it instead", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitorID))
		return existingReport, nil
	}
	changeList := make([]*models.DynamicChanges, 0)
	for _, pageHistory := range history {
		changeList = append(changeList, &pageHistory.DiffContent)
	}

	changes, err := s.aiService.SummarizeChanges(ctx, changeList)
	if err != nil {

		return nil, err
	}

	reportContent, err := s.renderHTML(competitorName, changes)
	if err != nil {
		return nil, err
	}

	report, err := s.repo.Set(ctx, workspaceID, competitorID, changes, reportContent)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// Dispatch send the report to it's subscribers.
func (s *reportService) Dispatch(ctx context.Context, workspaceID, competitorID uuid.UUID, competitorName string, subscriberEmails []string) error {
	report, err := s.GetLatest(ctx, workspaceID, competitorID)
	if err != nil {
		return err
	}

	reportContent, err := s.GetContent(ctx, report.URI)
	if err != nil {
		return err
	}

	email := models.Email{
		To:           subscriberEmails,
		EmailSubject: "Weekly Roundup for " + competitorName,
		EmailContent: reportContent,
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

func (s *reportService) renderHTML(competitorName string, changes []models.CategoryChange) (string, error) {
	tmp, err := s.library.GetTemplate(template.WeeklyRoundupTemplate)
	if err != nil {
		return "", err
	}

	// ERROR: template.SectionedTemplate is not a typecompilerNotAType
	sectionedTemplate, ok := tmp.(*template.SectionedTemplate)
	if !ok {
		return "", errors.New("failed to assert template to SectionedTemplate")
	}

	// Override template with report data
	sectionedTemplate.Competitor = competitorName
	sectionedTemplate.GeneratedAt = time.Now()
	sectionedTemplate.FromDate = time.Now().AddDate(0, 0, -7) // Assuming weekly report
	sectionedTemplate.ToDate = time.Now()

	// Create sections map
	sectionedTemplate.Sections = make(map[string]template.Section)

	// Map each CategoryChange to a Section
	for _, change := range changes {
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
		return "", err
	}

	return htmlContent, nil
}
