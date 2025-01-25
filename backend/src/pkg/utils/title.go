package utils

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TitleStrategy func(*html.Node, string) string

func GetPageTitle(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; GoTitleBot/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	strategies := []TitleStrategy{
		getTitleTag,        // Standard <title> tag
		getOGTitle,         // Open Graph title
		getFallbackFromUrl, // URL fallback
	}

	for _, strategy := range strategies {
		if title := strategy(doc, url); title != "" {
			return cleanTitle(title), nil
		}
	}

	return "", fmt.Errorf("no title found using any method")
}

// Strategy 1: Standard <title> tag
func getTitleTag(n *html.Node, _ string) string {
	var title string
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" && node.FirstChild != nil {
			title = node.FirstChild.Data
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(n)
	return title
}

// Strategy 2: Open Graph title
func getOGTitle(n *html.Node, _ string) string {
	var traverse func(*html.Node) string
	traverse = func(node *html.Node) string {
		if node.Type == html.ElementNode && node.Data == "meta" {
			var property, content string
			for _, attr := range node.Attr {
				if attr.Key == "property" && attr.Val == "og:title" {
					property = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}
			if property == "og:title" && content != "" {
				return content
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if result := traverse(child); result != "" {
				return result
			}
		}
		return ""
	}
	return traverse(n)
}

// Strategy 3: Extract from URL as last resort
func getFallbackFromUrl(_ *html.Node, url string) string {
	// Get the last part of the path
	lastSegment := path.Base(url)

	// Remove extension if present
	lastSegment = strings.TrimSuffix(lastSegment, path.Ext(lastSegment))

	// Replace hyphens and underscores with spaces
	lastSegment = strings.ReplaceAll(lastSegment, "-", " ")
	lastSegment = strings.ReplaceAll(lastSegment, "_", " ")

	// Title case the result using the proper Unicode-aware method
	caser := cases.Title(language.English)
	return caser.String(lastSegment)
}

// Helper function to clean and normalize titles
func cleanTitle(title string) string {
	// Remove extra whitespace
	title = strings.TrimSpace(title)
	title = strings.Join(strings.Fields(title), " ")

	// Decode common HTML entities
	title = strings.ReplaceAll(title, "&amp;", "&")
	title = strings.ReplaceAll(title, "&quot;", "\"")
	title = strings.ReplaceAll(title, "&#39;", "'")

	return title
}
