package competitor

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type CompanyNameFinder struct {
	delimiters   []string
	minFrequency int
	minLength    int
	client       *http.Client
	logger       *logger.Logger
}

type PageInfo struct {
	URL      string
	Title    string
	SiteName string
	Error    error
}

// NewCompanyNameFinder creates a new CompanyNameFinder with default values
func NewCompanyNameFinder(logger *logger.Logger) *CompanyNameFinder {
	return &CompanyNameFinder{
		delimiters:   []string{"|", "-", "•", ":", "»", "—", "–"},
		minFrequency: 2,
		minLength:    3,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger.WithFields(map[string]interface{}{"module": "company_name_finder"}),
	}
}

func extractDomainName(url string) string {
	// Remove protocol if exists
	if idx := strings.Index(url, "://"); idx != -1 {
		url = url[idx+3:]
	}

	// Get domain part (before path)
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}

	// Split by dots and get parts
	parts := strings.Split(url, ".")

	// Get the main domain name (usually second-to-last part)
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}

	return ""
}

func (f *CompanyNameFinder) splitByDelimiters(title string) []string {
	parts := []string{title}

	for _, delimiter := range f.delimiters {
		var newParts []string
		for _, part := range parts {
			splitParts := strings.Split(part, delimiter)
			for _, p := range splitParts {
				trimmed := strings.TrimSpace(p)
				if len(trimmed) >= f.minLength {
					newParts = append(newParts, trimmed)
				}
			}
		}
		parts = newParts
	}

	return parts
}

func (f *CompanyNameFinder) findMostFrequent(titles []string) string {
	frequency := make(map[string]int)

	for _, title := range titles {
		if title == "" {
			continue
		}

		parts := f.splitByDelimiters(title)
		seen := make(map[string]bool)

		for _, part := range parts {
			if !seen[part] {
				frequency[part]++
				seen[part] = true
			}
		}
	}

	var mostFrequent string
	maxCount := f.minFrequency - 1

	for part, count := range frequency {
		if count > maxCount && len(part) >= f.minLength {
			if count > maxCount || (count == maxCount && len(part) > len(mostFrequent)) {
				mostFrequent = part
				maxCount = count
			}
		}
	}

	return mostFrequent
}

func (f *CompanyNameFinder) cleanCommonSuffixes(name string) string {
	suffixes := []string{
		"Inc", "Inc.", "LLC", "Ltd", "Ltd.", "Limited",
		"Corp", "Corp.", "Corporation", "Home", "Homepage",
		"Official Website", "Official Site", "Company", "Co",
	}

	cleaned := name
	for _, suffix := range suffixes {
		cleaned = strings.TrimSuffix(cleaned, " "+suffix)
	}

	return strings.TrimSpace(cleaned)
}

func (f *CompanyNameFinder) fetchPage(url string) PageInfo {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	resp, err := f.client.Get(url)
	if err != nil {
		return PageInfo{URL: url, Error: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PageInfo{URL: url, Error: err}
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return PageInfo{URL: url, Error: err}
	}

	var siteName, title string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "meta" {
				var property, content string
				for _, attr := range n.Attr {
					if attr.Key == "property" && attr.Val == "og:site_name" {
						property = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if property == "og:site_name" && content != "" {
					siteName = content
				}
			}

			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	return PageInfo{
		URL:      url,
		Title:    title,
		SiteName: siteName,
	}
}

func (f *CompanyNameFinder) ProcessURLs(urls []string) string {
	var wg sync.WaitGroup
	pageInfoChan := make(chan PageInfo, len(urls))

	// Fetch all pages concurrently
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			pageInfo := f.fetchPage(url)
			pageInfoChan <- pageInfo
		}(url)
	}

	go func() {
		wg.Wait()
		close(pageInfoChan)
	}()

	var siteNames []string
	var titles []string
	domainFreq := make(map[string]int)

	pageErrorCount := 0
	for pageInfo := range pageInfoChan {
		// Count domain frequencies
		domain := extractDomainName(pageInfo.URL)
		if domain != "" {
			domainFreq[domain]++
		}

		if pageInfo.Error != nil {
			f.logger.Warn("Failed to fetch page", zap.Error(pageInfo.Error), zap.String("url", pageInfo.URL))
			pageErrorCount++
			continue
		}

		if pageInfo.SiteName != "" {
			siteNames = append(siteNames, pageInfo.SiteName)
		}
		if pageInfo.Title != "" {
			titles = append(titles, pageInfo.Title)
		}
	}

	if pageErrorCount == len(urls) {
		return "New Page" // All pages failed to load
	}

	// 1. First try using og:site_name
	if len(siteNames) > 0 {
		allSame := true
		first := siteNames[0]
		for _, name := range siteNames[1:] {
			if name != first {
				allSame = false
				break
			}
		}
		if allSame {
			return f.cleanCommonSuffixes(first)
		}

		if companyName := f.findMostFrequent(siteNames); companyName != "" {
			return f.cleanCommonSuffixes(companyName)
		}
	}

	// 2. Try title analysis
	if companyName := f.cleanCommonSuffixes(f.findMostFrequent(titles)); companyName != "" {
		return companyName
	}

	// 3. Final fallback: use most common domain name or uuid
	if companyName := f.pickMostCommonDomain(domainFreq); companyName != "" {
		return companyName
	}

	// 4. Fallback to uuids
	return uuid.NewString()
}

func (f *CompanyNameFinder) pickMostCommonDomain(domainFreq map[string]int) string {
	var mostCommonDomain string
	maxFreq := 0
	for domain, freq := range domainFreq {
		if freq > maxFreq {
			maxFreq = freq
			mostCommonDomain = domain
		}
	}

	if mostCommonDomain != "" {
		// Clean up domain name
		cleaned := strings.ToTitle(strings.ToLower(mostCommonDomain))
		cleaned = strings.ReplaceAll(cleaned, "-", " ")
		cleaned = strings.ReplaceAll(cleaned, "_", " ")
		return f.cleanCommonSuffixes(cleaned)
	}

	return mostCommonDomain
}
