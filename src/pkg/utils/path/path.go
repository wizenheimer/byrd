package path

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/wizenheimer/iris/src/internal/constants"
)

// GenerateURLHash normalizes the URL and generates a consistent hash
func GenerateURLHash(rawURL string) (string, error) {
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

	// Generate hash from normalized URL
	normalizedURL := u.String()
	hasher := sha256.New()
	hasher.Write([]byte(normalizedURL))
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

	year, weekNumber, runID := getCurrentTimeComponents()
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

	year, weekNumber, runID := getCurrentTimeComponents()
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

	year, weekNumber, runID := getPreviousTimeComponents()
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

	year, weekNumber, runID := getPreviousTimeComponents()
	path, err := GeneratePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("text/%s", path), nil
}

// getCurrentTimeComponents returns the current year, week number, and run ID
func getCurrentTimeComponents() (year, weekNumber int, runID string) {
	now := time.Now()
	year, weekNumber = now.ISOWeek()

	currentWeekDay := now.Weekday()
	if currentWeekDay == 0 { // Convert Sunday (0) to 7
		currentWeekDay = 7
	}

	if int(currentWeekDay) <= 3 { // Monday, Tuesday, Wednesday
		runID = constants.FirstRunID
	} else { // Thursday, Friday, Saturday, Sunday
		runID = constants.LastRunID
	}

	return year, weekNumber, runID
}

// getPreviousTimeComponents returns the previous year, week number, and run ID
func getPreviousTimeComponents() (year, weekNumber int, runID string) {
	now := time.Now()
	currentYear, currentWeek := now.ISOWeek()

	currentWeekDay := now.Weekday()
	if currentWeekDay == 0 { // Convert Sunday (0) to 7
		currentWeekDay = 7
	}

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

	return year, weekNumber, runID
}
