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
	SourceURL   string `json:"source_url"`
	RenderedURL string `json:"rendered_url"`
	Year        int    `json:"year"`
	WeekNumber  int    `json:"week_number"`
	WeekDay     int    `json:"week_day"`
}

func (s ScreenshotMetadata) ToMap() map[string]string {
	result := make(map[string]string)

	result["source_url"] = s.SourceURL
	result["rendered_url"] = s.RenderedURL
	result["year"] = strconv.Itoa(s.Year)
	result["week_day"] = strconv.Itoa(s.WeekDay)
	result["week_number"] = strconv.Itoa(s.WeekNumber)

	return result
}

// FromMap safely converts map[string]string to ScreenshotMetadata
func ScreenshotMetadataFromMap(m map[string]string) (ScreenshotMetadata, []error) {
	var result ScreenshotMetadata
	var errs []error

	// Required string fields
	if srcURL, exists := m["source_url"]; exists {
		result.SourceURL = srcURL
	} else {
		errs = append(errs, errors.New("missing required field: source_url"))
	}

	if rendURL, exists := m["rendered_url"]; exists {
		result.RenderedURL = rendURL
	} else {
		errs = append(errs, errors.New("missing required field: rendered_url"))
	}

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

	// Return errs if any occurred
	if len(errs) > 0 {
		return result, errs
	}

	return result, nil
}
