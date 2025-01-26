// ./src/internal/service/screenshot/utils.go
package screenshot

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"strconv"

	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// refreshScreenshot refreshes the screenshot and html content for the given URL
// it ensures that the screenshot and content are fetched and aren't null before returning
func (s *screenshotService) refreshScreenshot(_ context.Context, opts models.ScreenshotRequestOptions) (*image.Image, *string, error) {
	defaultOpt := models.GetDefaultScreenshotRequestOptions(opts.URL)
	mergedOpt := models.MergeScreenshotRequestOptions(defaultOpt, opts)

	resp, err := s.httpClient.NewRequest().
		BaseURL(s.origin).
		Method(http.MethodGet).
		Path("/take").
		QueryParam("access_key", s.key).
		AddQueryParamsFromStruct(mergedOpt).
		Execute(s.httpClient)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("failed to retrieve screenshot")
	}

	img, err := s.getScreenshot(resp)
	if err != nil {
		return nil, nil, err
	} else if img == nil {
		return nil, nil, errors.New("failed to retrieve screenshot")
	}

	content, err := s.getContent(resp)
	if err != nil {
		return nil, nil, err
	} else if content == nil {
		return nil, nil, errors.New("failed to retrieve content")
	}

	return img, content, nil
}

func getScreenshotMetadata(backDate bool) (*models.ScreenshotMetadata, error) {
	var currentDayString string
	var currentWeek int
	var currentYear int
	if backDate {
		currentYear, currentWeek, currentDayString = getPreviousTimeComponents(true)
	} else {
		currentYear, currentWeek, currentDayString = getCurrentTimeComponents(true)
	}
	currentDay, err := strconv.Atoi(currentDayString)
	if err != nil {
		return nil, err
	}

	return &models.ScreenshotMetadata{
		Year:       currentYear,
		WeekNumber: currentWeek,
		WeekDay:    currentDay,
	}, nil
}

func (s *screenshotService) getScreenshot(resp *http.Response) (*image.Image, error) {
	imageContentTypes := []string{
		"image/png",
	}

	contentType := resp.Header.Get("Content-Type")
	if !utils.Contains(imageContentTypes, contentType) {
		return nil, fmt.Errorf("received unexpected content type: %v", contentType)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func (s *screenshotService) getContent(resp *http.Response) (*string, error) {
	renderedURL := resp.Header.Get("X-ScreenshotOne-Content-URL")
	if renderedURL == "" {
		return nil, errors.New("no content URL found in headers, cannot proceed with rendering")
	}

	htmlResp, err := http.Get(renderedURL)
	if err != nil {
		return nil, err
	}
	defer htmlResp.Body.Close()
	if htmlResp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve content")
	}

	htmlBytes, err := io.ReadAll(htmlResp.Body)
	if err != nil {
		return nil, err
	}

	content := string(htmlBytes)
	return &content, nil
}
