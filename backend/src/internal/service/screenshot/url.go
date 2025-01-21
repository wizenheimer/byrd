package screenshot

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

type ContentType string

const (
	ContentTypeImage   ContentType = "image"
	ContentTypeContent ContentType = "content"
)

// validateHostname checks if a hostname is valid and not in the reject list
func validateHostname(hostname string) error {
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

// normalizeURL normalizes the URL by converting it to lowercase, removing www. prefix, cleaning the path, removing default ports, and sorting query parameters
func normalizeURL(rawURL string) (string, error) {
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

// preProcessURL follows URL redirects and validates the final destination
func preProcessURL(rawURL string) (string, error) {
	// Validate the hostname
	if err := validateHostname(rawURL); err != nil {
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
	normalized, err := normalizeURL(finalURL)
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

func GetPath(url string, contentType ContentType, backDate bool) (string, error) {
	switch contentType {
	case ContentTypeImage:
		// Get the current screenshot path
		if backDate {
			return getPreviousScreenshotPath(url)
		} else {
			return getCurrentScreenshotPath(url)
		}
	case ContentTypeContent:
		// Get the current content path
		if backDate {
			return getPreviousContentPath(url)
		} else {
			return getCurrentContentPath(url)
		}
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}
