package template

import (
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
