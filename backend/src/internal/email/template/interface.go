package template

import (
	"errors"
)

var (
	ErrTemplateNotFound      = errors.New("template not found")
	ErrTemplateAlreadyExists = errors.New("template already exists")
	ErrEmptyTemplateName     = errors.New("template name cannot be empty")
	ErrNilTemplate           = errors.New("template cannot be nil")
)

type Template interface {
	RenderHTML() (string, error)
	Copy() (Template, error)
}

// TemplateLibrary is a collection of pre-parameterized email templates.
type TemplateLibrary interface {
	// GetTemplate returns a template by name.
	// Returns a copy of the template, so that the original template is not modified.
	GetTemplate(name TemplateName) (Template, error)
	// RegisterTemplate registers a template by name
	// Returns an error if a template with the same name already exists.
	RegisterTemplate(name TemplateName, t Template) error
}
