// ./src/internal/service/screenshot/path.go
package screenshot

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/wizenheimer/byrd/src/internal/constants"
)

// generateURLHash normalizes the URL and generates a consistent hash
func generateURLHash(rawURL string) (string, error) {
	// pre-process the URL
	processedURL, err := preProcessURL(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to normalize URL: %w", err)
	}

	hasher := sha256.New()
	hasher.Write([]byte(processedURL))
	return hex.EncodeToString(hasher.Sum(nil))[:32], nil
}

// GeneratePath generates a path for a given hash, year, week number and run id
func generatePath(hash string, year, weekNumber int, runID string) (string, error) {
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
func getCurrentScreenshotPath(url string) (string, error) {
	hash, err := generateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := getCurrentTimeComponents(true)
	path, err := generatePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("images/%s", path), nil
}

// GetCurrentContentPath returns the path to the current content for a given url
func getCurrentContentPath(url string) (string, error) {
	hash, err := generateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := getCurrentTimeComponents(true)
	path, err := generatePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("text/%s", path), nil
}

// GetPreviousScreenshotPath returns the path to the previous screenshot for a given url
func getPreviousScreenshotPath(url string) (string, error) {
	hash, err := generateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := getPreviousTimeComponents(true)
	path, err := generatePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("images/%s", path), nil
}

// GetPreviousContentPath returns the path to the previous content for a given url
func getPreviousContentPath(url string) (string, error) {
	hash, err := generateURLHash(url)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL hash: %w", err)
	}

	year, weekNumber, runID := getPreviousTimeComponents(true)
	path, err := generatePath(hash, year, weekNumber, runID)
	if err != nil {
		return "", fmt.Errorf("failed to generate path: %w", err)
	}

	return fmt.Sprintf("text/%s", path), nil
}

// getCurrentTimeComponents returns the current year, week number, and run ID
func getCurrentTimeComponents(enableBucketing bool) (year, weekNumber int, runID string) {
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
func getPreviousTimeComponents(enableBucketing bool) (year, weekNumber int, runID string) {
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
