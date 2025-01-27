package template

import (
	"fmt"

	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type templateLibrary struct {
	templates map[TemplateName]Template
	logger    *logger.Logger
}

// NewTemplateLibrary creates a new TemplateLibrary instance
func NewTemplateLibrary(logger *logger.Logger) (TemplateLibrary, error) {
	tl := templateLibrary{
		logger: logger.WithFields(map[string]any{
			"component": "template-library",
		}),
		templates: make(map[TemplateName]Template),
	}

	// Register default templates
	if err := registerDefaultTemplates(&tl); err != nil {
		logger.Error("failed to register default templates", zap.Error(err))
		return nil, err
	}

	return &tl, nil
}

type TemplateName string

const (
	// -- user lifecycle templates --
	WelcomeTemplate  TemplateName = "welcome"
	WaitlistTemplate TemplateName = "waitlist"

	// -- workspace user lifecycle templates --
	RequestWorkspaceInviteTemplate   TemplateName = "request_workspace_invite"
	PendingWorkspaceInviteTemplate   TemplateName = "pending_workspace_invite"
	AcceptWorkspaceInviteTemplateFTU TemplateName = "accept_workspace_invite_ftu"
	AcceptWorkspaceInviteTemplateRU  TemplateName = "accept_workspace_invite_ru"
	DeclineWorkspaceInviteTemplate   TemplateName = "decline_workspace_invite"
	DeletedWorkspaceUserTemplate     TemplateName = "deleted_workspace_user"

	// -- trial lifecycle templates --
	TrailSucceededTemplate TemplateName = "trial_succeeded"

	// -- workspace lifecycle templates --
	DeletedWorkspaceTemplate TemplateName = "deleted_workspace"

	// -- subscription lifecycle templates --
	RenewalFailedTemplate    TemplateName = "renewal_failed"
	RenewalSucceededTemplate TemplateName = "renewal_succeeded"
	RenewalCanceledTemplate  TemplateName = "renewal_canceled"

	// -- weekly roundup templates --
	WeeklyRoundupTemplate TemplateName = "weekly_roundup"
)

// registerDefaultTemplates pre-registers all the default email templates
func registerDefaultTemplates(lib TemplateLibrary) error {
	// Register each template
	for name, tmpl := range templates {
		if err := lib.RegisterTemplate(name, tmpl); err != nil {
			return fmt.Errorf("failed to register template %s: %w", name, err)
		}
	}

	return nil
}

// GetTemplate returns a template by name.
func (tl *templateLibrary) GetTemplate(name TemplateName) (Template, error) {
	tl.logger.Debug("getting template", zap.Any("name", name))
	t, ok := tl.templates[name]
	if !ok {
		return nil, ErrTemplateNotFound
	}

	tc, err := t.Copy()
	if err != nil {
		return nil, err
	}

	return tc, nil
}

// RegisterTemplate registers a template by name.
func (tl *templateLibrary) RegisterTemplate(name TemplateName, t Template) error {
	tl.logger.Debug("registering template", zap.Any("name", name))
	if name == "" {
		return ErrEmptyTemplateName
	}
	if t == nil {
		return ErrNilTemplate
	}
	if _, ok := tl.templates[name]; ok {
		return ErrTemplateAlreadyExists
	}
	tl.templates[name] = t
	return nil
}
