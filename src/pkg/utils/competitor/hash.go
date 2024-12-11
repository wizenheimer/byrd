package competitor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/wizenheimer/iris/src/internal/constants"
)

// GeneratePathHash generates a hash for a given url
func GeneratePathHash(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))[:32]
}

// GetScreenshotPath returns the path to the screenshot for a given url, week number and week day
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

// GetCurrentScreenshotPath returns the path to the screenshot for a given url
func GetCurrentScreenshotPath(url string) string {
	path := GenerateCurrentPath(url)
	return fmt.Sprintf("screenshot/%s", path)
}

// GetPreviousScreenshotPath returns the path to the screenshot for a given url
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
	hash := GeneratePathHash(url)
	return fmt.Sprintf("%s/%s/%s", hash, weekNumber, runID)
}

func GenerateCurrentPath(url string) string {
	// Generate a hash for the url
	hash := GeneratePathHash(url)

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

	// Generate the path using the hash, week number and runID
	path := GeneratePath(hash, weekNumber, runID)
	return path
}

func GeneratePreviousPath(url string) string {
	// Generate a hash for the url
	hash := GeneratePathHash(url)

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
	path := GeneratePath(hash, prevWeek, prevRun)
	return path
}
