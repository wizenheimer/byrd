// ./src/internal/service/screenshot/utils.go
package screenshot

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

// refreshScreenshot refreshes the screenshot and html content for the given URL
// it ensures that the screenshot and content are fetched and aren't null before returning
func (s *screenshotService) refreshScreenshot(ctx context.Context, opts models.ScreenshotRequestOptions) (*image.Image, *string, error) {
	defaultOpt := models.GetDefaultScreenshotRequestOptions(opts.URL)
	mergedOpt := models.MergeScreenshotRequestOptions(defaultOpt, opts)

	req, err := s.createScreenshotRequest(ctx, http.MethodGet, "take", mergedOpt)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
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
	if htmlResp == nil {
		return nil, errors.New("received nil response from content URL")
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

// CreateScreenshotRequest creates a new HTTP request with query parameters from a struct
func (s *screenshotService) createScreenshotRequest(ctx context.Context, requestMethod, requestPath string, opts interface{}) (*http.Request, error) {
	requestURL := fmt.Sprintf("%s/%s", strings.TrimRight(s.origin, "/"), strings.TrimLeft(requestPath, "/"))
	u, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}

	// Convert struct to query parameters
	q := u.Query()
	q.Set("access_key", s.key)
	addStructToQuery(q, opts)
	u.RawQuery = q.Encode()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, requestMethod, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	return req, nil
}

// addStructToQuery converts a struct to query parameters
func addStructToQuery(q url.Values, data interface{}) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		parts := strings.Split(jsonTag, ",")
		if len(parts) == 0 {
			continue
		}
		name := parts[0]

		if field.IsZero() {
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				addQueryParam(q, name, field.Elem())
			}
		case reflect.Slice:
			if field.Len() > 0 {
				switch name {
				case "wait_until":
					for j := 0; j < field.Len(); j++ {
						q.Add("wait_until", fmt.Sprint(field.Index(j).Interface()))
					}
				case "scripts_wait_until":
					for j := 0; j < field.Len(); j++ {
						q.Add("scripts_wait_until", fmt.Sprint(field.Index(j).Interface()))
					}
				case "hide_selectors":
					for j := 0; j < field.Len(); j++ {
						q.Add("hide_selector", fmt.Sprint(field.Index(j).Interface()))
					}
				case "block_requests":
					for j := 0; j < field.Len(); j++ {
						q.Add("block_request", fmt.Sprint(field.Index(j).Interface()))
					}
				case "block_resources":
					for j := 0; j < field.Len(); j++ {
						q.Add("block_resources", fmt.Sprint(field.Index(j).Interface()))
					}
				default:
					for j := 0; j < field.Len(); j++ {
						q.Add(name, fmt.Sprint(field.Index(j).Interface()))
					}
				}
			}
		case reflect.Map:
			if field.Len() > 0 {
				iter := field.MapRange()
				for iter.Next() {
					paramName := fmt.Sprintf("%s[%s]", name, iter.Key().String())
					q.Add(paramName, fmt.Sprint(iter.Value().Interface()))
				}
			}
		case reflect.String:
			if field.String() != "" {
				q.Add(name, field.String())
			}
		}
	}
}

// addQueryParam adds a single query parameter
func addQueryParam(q url.Values, name string, value reflect.Value) {
	switch value.Kind() {
	case reflect.String:
		q.Add(name, value.String())
	case reflect.Bool:
		q.Add(name, strconv.FormatBool(value.Bool()))
	case reflect.Int, reflect.Int64:
		q.Add(name, strconv.FormatInt(value.Int(), 10))
	case reflect.Float64:
		q.Add(name, strconv.FormatFloat(value.Float(), 'f', -1, 64))
	}
}
