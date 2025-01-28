package template

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

type CommonTemplate struct {
	PreviewText string
	Logo        string // Optional, defaults to "byrd"
	Title       string
	Subtitle    string
	Body        []string      // List of paragraphs
	BulletTitle string        // Optional
	Bullets     []string      // Optional
	CTA         *CallToAction // Optional
	ClosingText string        // Optional
	Footer      Footer
	GeneratedAt time.Time
}

// CallToAction represents a button with its text and URL
type CallToAction struct {
	ButtonText string
	ButtonURL  string
	FooterText string // Text below the button
}

// Footer represents the email footer
type Footer struct {
	ContactMessage string // Optional
	ContactEmail   string // Optional
	Location       string // defaults to "ByrdLabs • San Francisco"
}

func (data *CommonTemplate) RenderHTML() (string, error) {
	const emailTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html dir="ltr" lang="en">
  <head>
    <meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
    <meta name="x-apple-disable-message-reformatting" />
  </head>
  <div style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0">{{.PreviewText}}<div></div></div>

  <body style="background-color:#ffffff;font-family:-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif">
    <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="max-width:600px;margin:0 auto;padding:40px 20px">
      <tbody>
        <tr style="width:100%">
          <td>
            <!-- Logo -->
            <p style="font-size:24px;line-height:24px;margin:16px 0;text-align:center;color:#000;margin-bottom:40px">{{if .Logo}}{{.Logo}}{{else}}byrd{{end}}</p>

            <!-- Title & Subtitle -->
            <p style="font-size:32px;line-height:1.3;margin:16px 0;font-weight:700;color:#000;margin-bottom:12px;text-align:center">{{.Title}}</p>
            {{if .Subtitle}}
            <p style="font-size:18px;line-height:24px;margin:16px 0;color:#666;margin-bottom:32px;text-align:center">{{.Subtitle}}</p>
            {{end}}

            <!-- Body Paragraphs -->
            <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="margin-bottom:32px">
              <tbody>
                <tr>
                  <td>
                    {{range .Body}}
                    <p style="font-size:16px;line-height:1.5;margin:16px 0;color:#333;margin-bottom:16px">{{.}}</p>
                    {{end}}
                  </td>
                </tr>
              </tbody>
            </table>

            {{if and .BulletTitle .Bullets}}
            <!-- Bullet Points Section -->
            <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="margin-bottom:32px">
              <tbody>
                <tr>
                  <td>
                    <p style="font-size:12px;line-height:24px;margin:16px 0;font-weight:500;color:#666;text-transform:uppercase;letter-spacing:0.05em;margin-bottom:16px">{{.BulletTitle}}</p>
                    <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="margin-bottom:24px">
                      <tbody>
                        <tr>
                          <td>
                            {{range .Bullets}}
                            <p style="font-size:14px;line-height:1.6;margin:16px 0;color:#333;margin-bottom:8px">• {{.}}</p>
                            {{end}}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </td>
                </tr>
              </tbody>
            </table>
            {{end}}

            {{if .CTA}}
            <!-- Call to Action -->
            <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="text-align:center;margin-bottom:32px">
              <tbody>
                <tr>
                  <td>
                    <a href="{{.CTA.ButtonURL}}" style="color:#fff;text-decoration-line:none;background-color:#000;padding:12px 32px;border-radius:6px;text-decoration:none;font-size:16px;font-weight:500;display:inline-block;margin-bottom:12px" target="_blank">{{.CTA.ButtonText}}</a>
                    {{if .CTA.FooterText}}
                    <p style="font-size:14px;line-height:24px;margin:16px 0;color:#666;font-style:italic;margin-top:12px">{{.CTA.FooterText}}</p>
                    {{end}}
                  </td>
                </tr>
              </tbody>
            </table>
            {{end}}

            {{if .ClosingText}}
            <!-- Closing Text -->
            <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="text-align:center;margin-bottom:32px">
              <tbody>
                <tr>
                  <td>
                    <p style="font-size:14px;line-height:24px;margin:16px 0;color:#666">{{.ClosingText}}</p>
                  </td>
                </tr>
              </tbody>
            </table>
            {{end}}

            <!-- Footer -->
            {{if or .Footer.ContactMessage .Footer.ContactEmail}}
            <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="text-align:center;margin-bottom:32px">
              <tbody>
                <tr>
                  <td>
                    {{if .Footer.ContactMessage}}
                    <p style="font-size:14px;line-height:24px;margin:16px 0;color:#666;margin-bottom:8px">{{.Footer.ContactMessage}}</p>
                    {{end}}
                    {{if .Footer.ContactEmail}}
                    <a href="mailto:{{.Footer.ContactEmail}}" style="color:#000;text-decoration-line:none;text-decoration:none;font-weight:500" target="_blank">{{.Footer.ContactEmail}}</a>
                    {{end}}
                  </td>
                </tr>
              </tbody>
            </table>
            {{end}}

            <p style="font-size:12px;line-height:24px;margin:16px 0;color:#666;text-align:center">{{if .Footer.Location}}{{.Footer.Location}}{{else}}ByrdLabs • San Francisco{{end}}</p>
          </td>
        </tr>
      </tbody>
    </table>
  </body>
</html>`

	tmpl, err := template.New("trial").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}

func (t *CommonTemplate) Copy() (Template, error) {
	copied := CommonTemplate{
		PreviewText: t.PreviewText,
		Logo:        t.Logo,
		Title:       t.Title,
		Subtitle:    t.Subtitle,
		BulletTitle: t.BulletTitle,
		ClosingText: t.ClosingText,
		Footer:      t.Footer,   // Struct copy is fine since Footer contains only strings
		GeneratedAt: time.Now(), // Refresh the timestamp on copy
	}

	// Deep copy slices
	if t.Body != nil {
		copied.Body = make([]string, len(t.Body))
		copy(copied.Body, t.Body)
	}

	if t.Bullets != nil {
		copied.Bullets = make([]string, len(t.Bullets))
		copy(copied.Bullets, t.Bullets)
	}

	// Deep copy CTA if it exists
	if t.CTA != nil {
		copied.CTA = &CallToAction{
			ButtonText: t.CTA.ButtonText,
			ButtonURL:  t.CTA.ButtonURL,
			FooterText: t.CTA.FooterText,
		}
	}

	return &copied, nil
}
