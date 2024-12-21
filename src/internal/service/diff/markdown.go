package diff

import (
	"fmt"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// MarkdownProcessor is a processor for markdown content
type MarkdownProcessor struct {
	minifier *MarkdownMinifier
}

// NewMarkdownProcessor creates a new markdown processor
func NewMarkdownProcessor() (*MarkdownProcessor, error) {
	minifier, err := NewMarkdownMinifier()
	if err != nil {
		return nil, fmt.Errorf("failed to create new markdown minifier: %w", err)
	}

	return &MarkdownProcessor{
		minifier: minifier,
	}, nil
}

// Process converts HTML content to markdown and minifies it for diffing
func (m *MarkdownProcessor) Process(htmlContent string) (string, error) {
	markdownContent, err := htmltomarkdown.ConvertString(htmlContent)
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to markdown: %w", err)
	}

	return m.minifier.Minify(markdownContent), nil
}

// MarkdownMinifier is a minifier for markdown content
type MarkdownMinifier struct {
	// State tracking
	inCodeBlock      bool
	inFencedBlock    bool
	lastLineEmpty    bool
	listIndentSize   int
	currentParagraph []string

	// Link cleaning patterns
	linkPatterns []*regexp.Regexp

	// Minification patterns
	listMarkerRegex      *regexp.Regexp
	orderedListRegex     *regexp.Regexp
	headerSpaceRegex     *regexp.Regexp
	blockquoteRegex      *regexp.Regexp
	emphasisRegex        []*regexp.Regexp
	linkSpaceRegex       []*regexp.Regexp
	multipleNewlineRegex *regexp.Regexp
	trailingSpaceRegex   *regexp.Regexp
	htmlCommentRegex     *regexp.Regexp
	markdownCommentRegex *regexp.Regexp
}

// NewMarkdownMinifier creates a new markdown minifier
func NewMarkdownMinifier() (*MarkdownMinifier, error) {
	m := &MarkdownMinifier{
		listIndentSize:   2,
		currentParagraph: make([]string, 0),
	}

	// Initialize link cleaning patterns
	linkPatterns := []string{
		// Empty links [text]()
		`\[([^\]]+)\]\(\)`,

		// Inline links [text](url)
		`\[([^\]]+)\]\([^)]+\)`,

		// Reference links [text][ref]
		`\[([^\]]+)\]\[[^\]]*\]`,

		// Reference definitions [ref]: url
		`^\[[^\]]+\]:\s*http[s]?://.*$`,

		// Empty images with URL ![](url)
		`!\[\]\([^)]+\)`,

		// Image with simple caption ![alt text](url)
		`!\[([^\]]+)\]\([^)]+\)`,

		// Image with complex multi-line caption and URL
		`!\[((?:[^\]]|\\\n)+)\]\([^)]+\)`,

		// Image with complex multi-line caption without URL
		`!\[((?:[^\]]|\\\n)+)\]\(\)`,

		// Nested image in link [![alt](img-url)](link-url)
		`\[!\[[^\]]*\]\([^)]+\)\]\([^)]+\)`,

		// Complex LinkedIn-style nested content
		`!\[ !\[[^\]]*\]\([^)]+\)[^\]]*\]\([^)]+\)`,

		// Raw image links with numbers
		`!\[ \d+\]\([^)]+\)`,

		// Complex multi-line LinkedIn-style content
		`!\[[^\]]*(?:\\\n[^\]]*)*\]\([^)]+\)`,

		// LinkedIn-style titled image links
		`\[!\[[^\]]+\]\([^)]+\)\s*\]`,

		// Catch-all for any remaining URLs in parentheses
		`\([^()]*(?:http|https):\/\/[^)]*\)`,
	}

	m.linkPatterns = make([]*regexp.Regexp, len(linkPatterns))
	for i, pattern := range linkPatterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile link pattern '%s': %v", pattern, err)
		}
		m.linkPatterns[i] = regex
	}

	// Initialize minification patterns
	var err error

	// Initialize comment removal patterns
	if m.htmlCommentRegex, err = regexp.Compile(`(?s)<!--.*?-->`); err != nil {
		return nil, fmt.Errorf("failed to compile HTML comment regex: %v", err)
	}

	if m.markdownCommentRegex, err = regexp.Compile(`(?m)^[^\S\r\n]*%%.*$`); err != nil {
		return nil, fmt.Errorf("failed to compile markdown comment regex: %v", err)
	}

	if m.listMarkerRegex, err = regexp.Compile(`^[-*+]\s+`); err != nil {
		return nil, fmt.Errorf("failed to compile list marker regex: %v", err)
	}

	if m.orderedListRegex, err = regexp.Compile(`^\d+\.\s+`); err != nil {
		return nil, fmt.Errorf("failed to compile ordered list regex: %v", err)
	}

	if m.headerSpaceRegex, err = regexp.Compile(`^(#{1,6})\s+`); err != nil {
		return nil, fmt.Errorf("failed to compile header space regex: %v", err)
	}

	if m.blockquoteRegex, err = regexp.Compile(`^>\s+`); err != nil {
		return nil, fmt.Errorf("failed to compile blockquote regex: %v", err)
	}

	if m.multipleNewlineRegex, err = regexp.Compile(`\n{3,}`); err != nil {
		return nil, fmt.Errorf("failed to compile newline regex: %v", err)
	}

	if m.trailingSpaceRegex, err = regexp.Compile(`[ \t]+\n`); err != nil {
		return nil, fmt.Errorf("failed to compile trailing space regex: %v", err)
	}

	// Compile emphasis patterns
	emphasisPatterns := []string{
		`\*\*([^*]+)\*\*`,
		`__([^_]+)__`,
		`\*([^*]+)\*`,
		`_([^_]+)_`,
	}

	m.emphasisRegex = make([]*regexp.Regexp, len(emphasisPatterns))
	for i, pattern := range emphasisPatterns {
		if m.emphasisRegex[i], err = regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("failed to compile emphasis regex: %v", err)
		}
	}

	// Compile link space patterns
	linkSpacePatterns := []string{
		`\[\s+([^\]]+)\s+\]`,
		`\(\s+([^)]+)\s+\)`,
	}

	m.linkSpaceRegex = make([]*regexp.Regexp, len(linkSpacePatterns))
	for i, pattern := range linkSpacePatterns {
		if m.linkSpaceRegex[i], err = regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("failed to compile link space regex: %v", err)
		}
	}

	return m, nil
}

func (m *MarkdownMinifier) removeComments(content string) string {
	// Remove HTML comments (including multi-line)
	content = m.htmlCommentRegex.ReplaceAllString(content, "")

	// Remove Markdown comments (lines starting with %%)
	content = m.markdownCommentRegex.ReplaceAllString(content, "")

	// Remove extra newlines that might have been created
	return m.multipleNewlineRegex.ReplaceAllString(content, "\n\n")
}

func (m *MarkdownMinifier) extractTextFromImageLink(content string) string {
	parts := strings.Split(content, "\\")

	var cleaned []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || strings.HasPrefix(part, "!") {
			continue
		}
		cleaned = append(cleaned, part)
	}

	return strings.Join(cleaned, "\n")
}

func (m *MarkdownMinifier) cleanLinks(content string) string {
	// Remove reference link definitions first
	lines := strings.Split(content, "\n")
	var cleanedLines []string

	for _, line := range lines {
		if !m.linkPatterns[3].MatchString(line) { // Reference definition pattern
			cleanedLines = append(cleanedLines, line)
		}
	}
	content = strings.Join(cleanedLines, "\n")

	// Handle inline links first
	content = m.extractTextFromImageLink(content)

	// Handle nested patterns first
	content = m.linkPatterns[8].ReplaceAllString(content, "$1") // Nested image in link
	content = m.linkPatterns[9].ReplaceAllString(content, "")   // Complex LinkedIn nested content
	content = m.linkPatterns[10].ReplaceAllString(content, "")  // Raw image links with numbers

	// Handle complex multi-line content
	content = m.linkPatterns[11].ReplaceAllStringFunc(content, func(match string) string {
		// Extract text between \n that isn't part of the URL
		parts := strings.Split(match, "\\")
		var cleaned []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" && !strings.HasPrefix(part, "http") && !strings.HasPrefix(part, "!") {
				cleaned = append(cleaned, part)
			}
		}
		return strings.Join(cleaned, " ")
	})

	// Handle all other patterns
	for i, pattern := range m.linkPatterns {
		if i != 3 && i != 8 && i != 9 && i != 10 && i != 11 {
			if i == 4 {
				// Remove empty images completely
				content = pattern.ReplaceAllString(content, "")
			} else {
				content = pattern.ReplaceAllString(content, "$1")
			}
		}
	}

	// Clean up any leftover backslashes and extra whitespace
	content = strings.ReplaceAll(content, "\\", " ")
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")

	// Final cleanup for any remaining URL patterns
	content = m.linkPatterns[len(m.linkPatterns)-2].ReplaceAllString(content, "") // LinkedIn titled images
	content = m.linkPatterns[len(m.linkPatterns)-1].ReplaceAllString(content, "") // Catch-all URLs in parentheses

	// Clean up any leftover brackets that might be empty now
	content = regexp.MustCompile(`\[\s*\]`).ReplaceAllString(content, "")

	return strings.TrimSpace(content)
}

func (m *MarkdownMinifier) flushParagraph(minified *[]string) {
	if len(m.currentParagraph) > 0 {
		*minified = append(*minified, strings.Join(m.currentParagraph, " "))
		m.currentParagraph = m.currentParagraph[:0]
	}
}

func (m *MarkdownMinifier) isStructuralLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "#") ||
		strings.HasPrefix(trimmed, "-") ||
		strings.HasPrefix(trimmed, "*") ||
		strings.HasPrefix(trimmed, "+") ||
		strings.HasPrefix(trimmed, ">") ||
		strings.HasPrefix(trimmed, "```") ||
		strings.HasPrefix(trimmed, "    ") ||
		strings.HasPrefix(trimmed, "\t") ||
		m.orderedListRegex.MatchString(trimmed)
}

// Minify minifies the markdown content
// It removes comments, cleans up links, and minifies the content
// It also handles list items, headers, blockquotes, emphasis, and horizontal rules etc.
// It is purpose built for reducing token count, looses some information during interchange
func (m *MarkdownMinifier) Minify(content string) string {
	// First remove all comments
	content = m.removeComments(content)

	// Then clean all links
	content = m.cleanLinks(content)

	// Reset state
	m.inCodeBlock = false
	m.inFencedBlock = false
	m.lastLineEmpty = false
	m.currentParagraph = m.currentParagraph[:0]

	lines := strings.Split(content, "\n")
	var minified []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "```") {
			m.flushParagraph(&minified)
			m.inFencedBlock = !m.inFencedBlock
			minified = append(minified, line)
			continue
		}

		if m.inFencedBlock {
			m.flushParagraph(&minified)
			minified = append(minified, line)
			continue
		}

		processedLine := m.processLine(line)

		if processedLine == "" {
			m.flushParagraph(&minified)
			if !m.lastLineEmpty && i > 0 {
				minified = append(minified, "")
				m.lastLineEmpty = true
			}
			continue
		}

		if m.isStructuralLine(processedLine) {
			m.flushParagraph(&minified)
			minified = append(minified, processedLine)
		} else {
			m.currentParagraph = append(m.currentParagraph, processedLine)
		}

		m.lastLineEmpty = false
	}

	m.flushParagraph(&minified)

	result := strings.Join(minified, "\n")
	return m.performFinalCleanup(result)
}

func (m *MarkdownMinifier) processLine(line string) string {
	leadingSpace := m.getLeadingSpace(line)
	line = strings.TrimSpace(line)

	if line == "" {
		return ""
	}

	if strings.HasPrefix(line, "#") {
		return m.minifyHeader(line)
	}

	if m.isListItem(line) {
		return leadingSpace + m.minifyListItem(line)
	}

	if m.isHorizontalRule(line) {
		return "---"
	}

	if strings.HasPrefix(line, ">") {
		return m.minifyBlockquote(line)
	}

	line = m.minifyEmphasis(line)

	return leadingSpace + line
}

func (m *MarkdownMinifier) getLeadingSpace(line string) string {
	spaces := ""
	for _, char := range line {
		if char == ' ' {
			spaces += " "
		} else {
			break
		}
	}
	return spaces
}

func (m *MarkdownMinifier) minifyHeader(line string) string {
	return m.headerSpaceRegex.ReplaceAllString(line, "$1 ")
}

func (m *MarkdownMinifier) minifyListItem(line string) string {
	line = m.listMarkerRegex.ReplaceAllString(line, "- ")
	line = m.orderedListRegex.ReplaceAllString(line, "- ")
	return line
}

func (m *MarkdownMinifier) minifyBlockquote(line string) string {
	return m.blockquoteRegex.ReplaceAllString(line, "> ")
}

func (m *MarkdownMinifier) minifyEmphasis(line string) string {
	for _, regex := range m.emphasisRegex {
		line = regex.ReplaceAllString(line, "*$1*")
	}
	return line
}

func (m *MarkdownMinifier) isListItem(line string) bool {
	return m.listMarkerRegex.MatchString(line) || m.orderedListRegex.MatchString(line)
}

func (m *MarkdownMinifier) isHorizontalRule(line string) bool {
	patterns := []string{
		`^-{3,}$`,
		`^\*{3,}$`,
		`^_{3,}$`,
		`^-\s+-\s+-[\s-]*$`,
		`^\*\s+\*\s+\*[\s\*]*$`,
		`^_\s+_\s+_[\s_]*$`,
	}

	trimmedLine := strings.TrimSpace(line)
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, trimmedLine)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func (m *MarkdownMinifier) performFinalCleanup(content string) string {
	content = m.multipleNewlineRegex.ReplaceAllString(content, "\n\n")
	content = m.trailingSpaceRegex.ReplaceAllString(content, "\n")
	content = strings.Trim(content, "\n")
	return content + "\n"
}
