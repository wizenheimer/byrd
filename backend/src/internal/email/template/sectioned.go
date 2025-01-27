package template

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// SectionedTemplate represents a template with multiple sections
type SectionedTemplate struct {
	Competitor  string
	FromDate    time.Time
	ToDate      time.Time
	GeneratedAt time.Time
	Summary     string
	Sections    map[string]Section
}

// Section represents a section in the email
type Section struct {
	Title   string
	Summary string
	Bullets []BulletPoint
}

// BulletPoint represents a single bullet point with its link
type BulletPoint struct {
	Text    string
	LinkURL string
}

// formatDate formats a date as DD/MM/YYYY
func formatDate(t time.Time) string {
	return t.Format("02/01/2006")
}

func (st *SectionedTemplate) RenderHTML() (string, error) {
	const emailTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html dir="ltr" lang="en">
<head>
    <meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
    <meta name="x-apple-disable-message-reformatting" />
</head>

<body style="background-color:#ffffff;font-family:-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif">
    <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="max-width:600px;margin:0 auto;padding:40px 20px">
        <tbody>
            <tr style="width:100%">
                <td>
                    <!-- Logo -->
                    <p style="font-size:24px;line-height:24px;margin:16px 0;text-align:center;color:#000;margin-bottom:40px">byrd</p>

                    <!-- Title -->
                    <p style="font-size:24px;line-height:1.3;margin:16px 0;font-weight:700;color:#000;margin-bottom:20px;text-align:center">Weekly Roundup for {{.Competitor}}</p>

                    <!-- Date Range -->
                    <p style="font-size:14px;line-height:24px;margin:16px 0;color:#666;margin-bottom:24px;text-align:center">{{formatDate .FromDate}} → {{formatDate .ToDate}}</p>

                    {{if .Summary}}
                    <!-- Summary -->
                    <p style="font-size:16px;line-height:1.5;margin:16px 0;color:#333;margin-bottom:32px;text-align:center">{{.Summary}}</p>
                    {{end}}

                    {{range $sectionKey, $section := .Sections}}
                    <!-- {{$sectionKey}} Section -->
                    <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="margin-bottom:24px">
                        <tbody>
                            <tr>
                                <td>
                                    <p style="font-size:12px;font-weight:700;color:#666;text-transform:uppercase;letter-spacing:0.05em;margin-bottom:12px">{{$section.Title}}</p>
                                    <p style="font-size:14px;line-height:1.5;color:#333;margin-bottom:16px">{{$section.Summary}}</p>
                                    <div style="margin:12px 0">
                                        {{range $bullet := $section.Bullets}}
                                        <p style="font-size:14px;color:#000;line-height:1.4;padding-left:16px;text-indent:-16px;margin:8px 0">
                                            • {{$bullet.Text}} {{if $bullet.LinkURL}}<a href="{{$bullet.LinkURL}}" style="font-size:13px;color:#666;text-decoration:none;font-weight:400">· Learn more</a>{{end}}
                                        </p>
                                        {{end}}
                                    </div>
                                    <hr style="border:none;border-top:1px solid #eaeaea;margin:24px 0" />
                                </td>
                            </tr>
                        </tbody>
                    </table>
                    {{end}}

                    <!-- Footer -->
                    <p style="font-size:12px;line-height:24px;margin:16px 0;color:#666;text-align:center">{{formatDate .GeneratedAt}}</p>
                    <p style="font-size:12px;line-height:24px;margin:16px 0;color:#666;text-align:center">ByrdLabs • San Francisco</p>
                </td>
            </tr>
        </tbody>
    </table>
</body>
</html>`

	tmpl, err := template.New("email").Funcs(template.FuncMap{
		"formatDate": formatDate,
	}).Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, st); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}

// Copy implements the Template interface
func (st *SectionedTemplate) Copy() (Template, error) {
	copied := &SectionedTemplate{
		Competitor:  st.Competitor,
		FromDate:    st.FromDate,
		ToDate:      st.ToDate,
		GeneratedAt: time.Now(), // Refresh the timestamp on copy
		Summary:     st.Summary,
		Sections:    make(map[string]Section),
	}

	// Deep copy sections
	for key, section := range st.Sections {
		// Copy bullets
		copiedBullets := make([]BulletPoint, len(section.Bullets))
		copy(copiedBullets, section.Bullets)

		copied.Sections[key] = Section{
			Title:   section.Title,
			Summary: section.Summary,
			Bullets: copiedBullets,
		}
	}

	return copied, nil
}
