package models

type EmailFormat string

const (
	EmailFormatHTML EmailFormat = "html"
	EmailFormatText EmailFormat = "text"
)

type Email struct {
	To           []string
	EmailFormat  EmailFormat
	EmailContent string
	EmailSubject string
}
