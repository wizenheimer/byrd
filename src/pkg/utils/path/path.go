package path

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wizenheimer/iris/src/internal/constants"
)

// GenerateURLHash normalizes the URL and generates a consistent hash
func GenerateURLHash(rawURL string) (string, error) {
	// Parse the URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
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

// GetScreenshotPath returns the path to
// the screenshot for a given url, week number and week day
func GetScreenshotPath(url, weekNumber, weekDay string) string {
	// Bucket the week day into a runID
	weekDayInt, _ := strconv.Atoi(weekDay)
	var runID string
	if weekDayInt <= 3 { // Monday, Tuesday, Wednesday
		runID = constants.FirstRunID
	} else { // Thursday, Friday, Saturday, Sunday
		runID = constants.LastRunID
	}

	path := GeneratePath(url, weekNumber, runID)
	return fmt.Sprintf("screenshot/%s", path)
}

// GetCurrentScreenshotPath returns the path to
// the screenshot for a given url
func GetCurrentScreenshotPath(url string) string {
	path := GenerateCurrentPath(url)
	return fmt.Sprintf("screenshot/%s", path)
}

// GetPreviousScreenshotPath returns the path to
// the screenshot for a given url
func GetPreviousScreenshotPath(url string) string {
	path := GeneratePreviousPath(url)
	return fmt.Sprintf("screenshot/%s", path)
}

// GetContentPath returns the path to the content for a given url, week number and run id
func GetContentPath(url, weekNumber, weekDay string) string {
	// Bucket the week day into a runID
	weekDayInt, _ := strconv.Atoi(weekDay)
	var runID string
	if weekDayInt <= 3 { // Monday, Tuesday, Wednesday
		runID = constants.FirstRunID
	} else { // Thursday, Friday, Saturday, Sunday
		runID = constants.LastRunID
	}

	path := GeneratePath(url, weekNumber, runID)
	return fmt.Sprintf("content/%s", path)
}

// GetCurrentContentPath returns the path to the content for a given url
func GetCurrentContentPath(url string) string {
	path := GenerateCurrentPath(url)
	return fmt.Sprintf("content/%s", path)
}

// GetPreviousContentPath returns the path to the content for a given url
func GetPreviousContentPath(url string) string {
	path := GeneratePreviousPath(url)
	return fmt.Sprintf("content/%s", path)
}

// GeneratePath generates a path for a given url, week number and run id
func GeneratePath(url, weekNumber, runID string) string {
	hash, err := GenerateURLHash(url)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", hash, weekNumber, runID)
}

func GenerateCurrentPath(url string) string {

	// Get the current weekday
	currentWeekDay := time.Now().Weekday()
	if currentWeekDay == 0 { // Convert Sunday (0) to 7
		currentWeekDay = 7
	}

	// Bucket the current weekday into a runID
	var runID string
	if currentWeekDay <= 3 { // Monday, Tuesday, Wednesday
		runID = constants.FirstRunID
	} else { // Thursday, Friday, Saturday, Sunday
		runID = constants.LastRunID
	}

	// Generate a week number
	_, currentWeekNumber := time.Now().ISOWeek()
	weekNumber := fmt.Sprintf("%d", currentWeekNumber)

	// Generate the path using the url, week number and runID
	return GeneratePath(url, weekNumber, runID)
}

func GeneratePreviousPath(url string) string {
	// Get the current weekday
	currentWeekday := time.Now().Weekday()
	if currentWeekday == 0 { // Convert Sunday (0) to 7
		currentWeekday = 7
	}

	// Bucket the current weekday into a runID
	var prevRunID string
	if currentWeekday <= 3 { // Monday, Tuesday, Wednesday
		prevRunID = constants.LastRunID
	} else { // Thursday, Friday, Saturday, Sunday
		prevRunID = constants.FirstRunID
	}

	// Generate a week number
	_, week := time.Now().ISOWeek()

	// Calculate the previous runID and weekNumber
	var prevWeek string
	var prevRun string

	// If the current runID is the last run of the current week
	// then the previous runID is the first run of the current week
	if prevRunID == constants.LastRunID {
		prevWeek = fmt.Sprintf("%02d", week)
	} else {
		// If the current runID is the first run of the current week
		// then the previous runID is the last run of the previous week

		// Perform a wrap around to week 52 if the current week is 1
		if week > 1 {
			prevWeek = fmt.Sprintf("%02d", week-1)
		} else {
			prevWeek = "52" // Wrap around to week 52
		}
	}

	// Generate the path
	return GeneratePath(url, prevWeek, prevRun)
}
