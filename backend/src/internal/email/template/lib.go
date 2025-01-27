package template

import (
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type templateLibrary struct {
	templates map[string]Template
	logger    *logger.Logger
}

func NewTemplateLibrary(logger *logger.Logger) TemplateLibrary {
	return &templateLibrary{
		logger: logger.WithFields(map[string]any{
			"component": "template-library",
		}),
		templates: make(map[string]Template),
	}
}

func (tl *templateLibrary) GetTemplate(name string) (Template, error) {
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

func (tl *templateLibrary) RegisterTemplate(name string, t Template) error {
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
