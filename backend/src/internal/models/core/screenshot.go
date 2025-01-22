package models

import (
	"errors"
	"fmt"
	"image"
	"strconv"
)

// ScreenshotContent defines the response structure for screenshot content requests
type ScreenshotContent struct {
	StoragePath string              `json:"path"`
	Content     string              `json:"content"`
	Metadata    *ScreenshotMetadata `json:"metadata,omitempty"`
}

// ScreenshotImage defines the response structure for screenshot image requests
type ScreenshotImage struct {
	StoragePath string              `json:"path"`
	Image       image.Image         `json:"image"`
	Metadata    *ScreenshotMetadata `json:"metadata,omitempty"`
}

// ScreenshotMetadata defines complete metadata for a screenshot
type ScreenshotMetadata struct {
	Year       int `json:"year"`
	WeekNumber int `json:"week_number"`
	WeekDay    int `json:"week_day"`
}

// ToMap safely converts ScreenshotMetadata to map[string]string.
// Raises validation errors if any required field is missing or invalid.
func (s ScreenshotMetadata) ToMap() (map[string]string, error) {
	result := make(map[string]string)
	var errs []error

	// Validate Year (assuming valid years are between 2000 and 2100)
	if s.Year < 2000 || s.Year > 2100 {
		errs = append(errs, fmt.Errorf("year must be between 2000 and 2100, got: %d", s.Year))
	}
	result["year"] = strconv.Itoa(s.Year)

	// Validate WeekDay (0-6, where 0 is Sunday)
	if s.WeekDay < 0 || s.WeekDay > 6 {
		errs = append(errs, fmt.Errorf("week_day must be between 0 and 6, got: %d", s.WeekDay))
	}
	result["week_day"] = strconv.Itoa(s.WeekDay)

	// Validate WeekNumber (1-53)
	if s.WeekNumber < 1 || s.WeekNumber > 53 {
		errs = append(errs, fmt.Errorf("week_number must be between 1 and 53, got: %d", s.WeekNumber))
	}
	result["week_number"] = strconv.Itoa(s.WeekNumber)

	// Return errors if any occurred
	if len(errs) > 0 {
		return nil, fmt.Errorf("validation errors: %v", errs)
	}

	return result, nil
}

// FromMap safely converts map[string]string to ScreenshotMetadata.
// Raises an error if any required field is missing or invalid.
func ScreenshotMetadataFromMap(m map[string]string) (*ScreenshotMetadata, error) {
	var result ScreenshotMetadata
	var errs []error

	// Required integer fields
	if year, exists := m["year"]; exists {
		if y, err := strconv.Atoi(year); err == nil {
			result.Year = y
		} else {
			errs = append(errs, fmt.Errorf("invalid year: %s", err))
		}
	} else {
		errs = append(errs, errors.New("missing required field: year"))
	}

	if weekday, exists := m["week_day"]; exists {
		if wd, err := strconv.Atoi(weekday); err == nil {
			result.WeekDay = wd
		} else {
			errs = append(errs, fmt.Errorf("invalid week_day: %s", err))
		}
	} else {
		errs = append(errs, errors.New("missing required field: week_day"))
	}

	if weeknumber, exists := m["week_number"]; exists {
		if wn, err := strconv.Atoi(weeknumber); err == nil {
			result.WeekNumber = wn
		} else {
			errs = append(errs, fmt.Errorf("invalid week_number: %s", err))
		}
	} else {
		errs = append(errs, errors.New("missing required field: week_number"))
	}

	// Return errors if any occurred
	if len(errs) > 0 {
		return nil, fmt.Errorf("validation errors: %v", errs)
	}

	return &result, nil
}
