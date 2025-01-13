// ./src/pkg/utils/path.go
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wizenheimer/byrd/src/internal/constants"
)

// GenerateURLHash normalizes the URL and generates a consistent hash
func GenerateURLHash(rawURL string) (string, error) {
	// pre-process the URL
	processedURL, err := PreProcessURL(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to normalize URL: %w", err)
	}

	hasher := sha256.New()
	hasher.Write([]byte(processedURL))
	return hex.EncodeToString(hasher.Sum(nil))[:32], nil
}

// GetScreenshotPath returns the path to the screenshot for a given url, year, week number and week day
func GetScreenshotPath(url string, year, weekNumber, weekDay int) (string, error) {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	var runID string
	if weekDay <= 3 { // Monday, Tuesday, Wednesday
		runID = constants.FirstRunID
	} else { // Thursday, Friday, Saturday, Sunday
		runID = constants.LastRunID
	}

	path, err := GeneratePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("images/%s", path), nil
}

// GetContentPath returns the path to the content for a given url, year, week number and week day
func GetContentPath(url string, year, weekNumber, weekDay int) (string, error) {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	var runID string
	if weekDay <= 3 { // Monday, Tuesday, Wednesday
		runID = constants.FirstRunID
	} else { // Thursday, Friday, Saturday, Sunday
		runID = constants.LastRunID
	}

	path, err := GeneratePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("text/%s", path), nil
}

// GeneratePath generates a path for a given hash, year, week number and run id
func GeneratePath(hash string, year, weekNumber int, runID string) (string, error) {
	if weekNumber < 1 || weekNumber > 52 {
		return "", fmt.Errorf("invalid week number: %d", weekNumber)
	}

	if year < 2000 || year > 2100 {
		return "", fmt.Errorf("invalid year: %d", year)
	}

	if runID != constants.FirstRunID && runID != constants.LastRunID {
		return "", fmt.Errorf("invalid run ID: %s", runID)
	}

	// Generates a path that sorts in reverse chronological order
	// This is useful for listing the most recent content first
	reverseYear := 9999 - year
	reverseWeek := 53 - weekNumber // 53 since ISO weeks go from 1-53
	reverseRun := constants.LastRunID
	if runID == constants.LastRunID {
		reverseRun = constants.FirstRunID // Reverse the run IDs too
	}

	return fmt.Sprintf("%s/%04d-%02d-%s", hash, reverseYear, reverseWeek, reverseRun), nil
}

// GetCurrentScreenshotPath returns the path to the current screenshot for a given url
func GetCurrentScreenshotPath(url string) (string, error) {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := GetCurrentTimeComponents(true)
	path, err := GeneratePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("images/%s", path), nil
}

// GetCurrentContentPath returns the path to the current content for a given url
func GetCurrentContentPath(url string) (string, error) {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := GetCurrentTimeComponents(true)
	path, err := GeneratePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("text/%s", path), nil
}

// GetPreviousScreenshotPath returns the path to the previous screenshot for a given url
func GetPreviousScreenshotPath(url string) (string, error) {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := GetPreviousTimeComponents(true)
	path, err := GeneratePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("images/%s", path), nil
}

// GetPreviousContentPath returns the path to the previous content for a given url
func GetPreviousContentPath(url string) (string, error) {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := GetPreviousTimeComponents(true)
	path, err := GeneratePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("text/%s", path), nil
}

// getCurrentTimeComponents returns the current year, week number, and run ID
func GetCurrentTimeComponents(enableBucketing bool) (year, weekNumber int, runID string) {
	now := time.Now()
	year, weekNumber = now.ISOWeek()

	currentWeekDay := now.Weekday()
	if currentWeekDay == 0 { // Convert Sunday (0) to 7
		currentWeekDay = 7
	}

	if enableBucketing {
		if int(currentWeekDay) <= 3 { // Monday, Tuesday, Wednesday
			runID = constants.FirstRunID
		} else { // Thursday, Friday, Saturday, Sunday
			runID = constants.LastRunID
		}
	} else {
		runID = strconv.Itoa(int(currentWeekDay))
	}

	return year, weekNumber, runID
}

// getPreviousTimeComponents returns the previous year, week number, and run ID
func GetPreviousTimeComponents(enableBucketing bool) (year, weekNumber int, runID string) {
	now := time.Now()
	currentYear, currentWeek := now.ISOWeek()

	currentWeekDay := now.Weekday()
	if currentWeekDay == 0 { // Convert Sunday (0) to 7
		currentWeekDay = 7
	}

	if enableBucketing {
		if int(currentWeekDay) <= 3 { // Monday, Tuesday, Wednesday
			runID = constants.LastRunID
			if currentWeek > 1 {
				year = currentYear
				weekNumber = currentWeek - 1
			} else {
				year = currentYear - 1
				weekNumber = 52
			}
		} else { // Thursday, Friday, Saturday, Sunday
			runID = constants.FirstRunID
			year = currentYear
			weekNumber = currentWeek
		}
	} else {
		runID = currentWeekDay.String()
	}

	return year, weekNumber, runID
}

func GetListingPrefixFromContentType(url, contentType string) (*string, error) {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return nil, fmt.Errorf("failed to generate URL hash: %w", err)
	}

	var prefix string
	switch contentType {
	case "image", "screenshot":
		contentType = "images"
	case "content", "text", "html":
		contentType = "text"
	default:
		return nil, fmt.Errorf("invalid content type: %s", contentType)
	}

	prefix = fmt.Sprintf("%s/%s", contentType, hash)

	return &prefix, nil
}

// ValidateHostname checks if a hostname is valid and not in the reject list
func ValidateHostname(hostname string) error {
	// Define rejected domains
	rejectedDomains := map[string]bool{
		"localhost":   true,
		"127.0.0.1":   true,
		"0.0.0.0":     true,
		"::1":         true,
		"example.com": true,
		"example.org": true,
		"example.net": true,
		"invalid":     true,
		"test":        true,
		"intranet":    true,
		"internal":    true,
		"local":       true,
		"private":     true,
		"corporative": true,
	}

	// Split host and port
	hostParts := strings.Split(strings.ToLower(hostname), ":")
	domainName := hostParts[0]

	// Check for rejected domains and subdomains
	domainParts := strings.Split(domainName, ".")
	for _, part := range domainParts {
		if rejectedDomains[part] {
			return fmt.Errorf("rejected domain: %s", hostname)
		}
	}

	// Reject IP addresses in private ranges
	if ip := net.ParseIP(domainName); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return fmt.Errorf("rejected IP address: %s", domainName)
		}
	}

	return nil
}

// NormalizeURL normalizes the URL by converting it to lowercase, removing www. prefix, cleaning the path, removing default ports, and sorting query parameters
func NormalizeURL(rawURL string) (string, error) {
	// Parse the URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Normalize the URL components
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// Remove www. prefix if present
	u.Host = strings.TrimPrefix(u.Host, "www.")

	// Clean the path (resolve relative paths, remove double slashes)
	u.Path = path.Clean(u.Path)

	// Remove default ports
	if (u.Scheme == "http" && strings.HasSuffix(u.Host, ":80")) ||
		(u.Scheme == "https" && strings.HasSuffix(u.Host, ":443")) {
		u.Host = strings.TrimSuffix(u.Host, ":80")
		u.Host = strings.TrimSuffix(u.Host, ":443")
	}

	// Sort query parameters
	if u.RawQuery != "" {
		query := u.Query()
		params := make([]string, 0, len(query))

		// Get sorted keys
		for k := range query {
			params = append(params, k)
		}
		sort.Strings(params)

		// Build sorted query string
		var sortedQuery strings.Builder
		for i, k := range params {
			if i > 0 {
				sortedQuery.WriteByte('&')
			}
			values := query[k]
			sort.Strings(values) // Sort multiple values for same key
			for j, v := range values {
				if j > 0 {
					sortedQuery.WriteByte('&')
				}
				sortedQuery.WriteString(url.QueryEscape(k))
				sortedQuery.WriteByte('=')
				sortedQuery.WriteString(url.QueryEscape(v))
			}
		}
		u.RawQuery = sortedQuery.String()
	}

	// Remove fragment as it doesn't affect the content
	u.Fragment = ""

	return u.String(), nil
}

// PreProcessURL follows URL redirects and validates the final destination
func PreProcessURL(rawURL string) (string, error) {
	// Validate the hostname
	if err := ValidateHostname(rawURL); err != nil {
		return "", fmt.Errorf("unsupported hostname for the URL: %w", err)
	}

	// URLShortenerPatterns contains known URL shortener domains
	var URLShortenerPatterns = []string{
		"bit.ly",
		"tinyurl.com",
		"t.co",
		"goo.gl",
		"ow.ly",
		"is.gd",
		"buff.ly",
		"tiny.cc",
		"adf.ly",
		"bit.do",
		"soo.gd",
		"s2r.co",
		"trib.al",
		"flip.it",
	}

	MaxRedirects := 3

	// Parse initial URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse initial URL: %w", err)
	}

	// Configure HTTP client with custom settings for URL shorteners
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= MaxRedirects {
				return fmt.Errorf("too many redirects (max %d)", MaxRedirects)
			}

			// Copy original headers for URL shortener services
			for key, values := range via[0].Header {
				req.Header[key] = values
			}

			// Add common headers that some URL shorteners check
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; URLValidator/1.0)")
			return nil
		},
		Timeout: 10 * time.Second,
	}

	// Check if it's a known URL shortener
	isShortener := false
	for _, pattern := range URLShortenerPatterns {
		if strings.Contains(parsedURL.Host, pattern) {
			isShortener = true
			break
		}
	}

	var resp *http.Response
	if isShortener {
		// For URL shorteners, try both HEAD and GET requests
		resp, err = client.Head(rawURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			// Some URL shorteners don't support HEAD requests
			// Fall back to GET request
			if resp != nil {
				resp.Body.Close()
			}
			resp, err = client.Get(rawURL)
		}
	} else {
		// For regular URLs, HEAD request is sufficient
		resp, err = client.Head(rawURL)
	}

	if err != nil {
		return "", fmt.Errorf("failed to follow redirects: %w", err)
	}
	defer resp.Body.Close()

	// Get the final URL after redirects
	finalURL := resp.Request.URL.String()

	// Handle cases where the final URL might be in a header
	// Some URL shorteners use custom headers for the final destination
	if location := resp.Header.Get("X-Final-Location"); location != "" {
		finalURL = location
	}

	// Additional validation for the final URL
	normalized, err := NormalizeURL(finalURL)
	if err != nil {
		return "", fmt.Errorf("invalid final URL: %w", err)
	}

	// Ensure the final URL is absolute
	finalParsed, err := url.Parse(normalized)
	if err != nil {
		return "", fmt.Errorf("failed to parse final URL: %w", err)
	}
	if !finalParsed.IsAbs() {
		return "", fmt.Errorf("final URL must be absolute")
	}

	return normalized, nil
}
