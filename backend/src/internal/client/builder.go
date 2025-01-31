// ./src/internal/client/builder.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// RequestBuilder helps construct HTTP requests
type requestBuilder struct {
	method      string
	baseURL     string
	path        string
	queryParams url.Values
	headers     map[string]string
	body        io.Reader
	ctx         context.Context
	err         error
}

func (rb *requestBuilder) Method(method string) RequestBuilder {
	rb.method = method
	return rb
}

func (rb *requestBuilder) BaseURL(baseURL string) RequestBuilder {
	rb.baseURL = baseURL
	return rb
}

func (rb *requestBuilder) Path(path string) RequestBuilder {
	rb.path = path
	return rb
}

func (rb *requestBuilder) QueryParam(key, value string) RequestBuilder {
	rb.queryParams.Add(key, value)
	return rb
}

func (rb *requestBuilder) Header(key, value string) RequestBuilder {
	rb.headers[key] = value
	return rb
}

func (rb *requestBuilder) JSON(data interface{}) RequestBuilder {
	if rb.err != nil {
		return rb
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		rb.err = fmt.Errorf("marshaling JSON: %w", err)
		return rb
	}

	rb.body = bytes.NewBuffer(jsonData)
	rb.headers["Content-Type"] = "application/json"
	return rb
}

func (rb *requestBuilder) Form(data url.Values) RequestBuilder {
	rb.body = strings.NewReader(data.Encode())
	rb.headers["Content-Type"] = "application/x-www-form-urlencoded"
	return rb
}

func (rb *requestBuilder) MultipartForm(fn func(*multipart.Writer) error) RequestBuilder {
	if rb.err != nil {
		return rb
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := fn(writer); err != nil {
		rb.err = fmt.Errorf("writing multipart form: %w", err)
		return rb
	}

	if err := writer.Close(); err != nil {
		rb.err = fmt.Errorf("closing multipart writer: %w", err)
		return rb
	}

	rb.body = &body
	rb.headers["Content-Type"] = writer.FormDataContentType()
	return rb
}

func (rb *requestBuilder) Context(ctx context.Context) RequestBuilder {
	rb.ctx = ctx
	return rb
}

func (rb *requestBuilder) Build() (*http.Request, error) {
	if rb.err != nil {
		return nil, rb.err
	}

	url := rb.baseURL + rb.path
	if len(rb.queryParams) > 0 {
		url += "?" + rb.queryParams.Encode()
	}

	req, err := http.NewRequestWithContext(rb.ctx, rb.method, url, rb.body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	for k, v := range rb.headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (rb *requestBuilder) Execute(c HTTPClient) (*http.Response, error) {
	req, err := rb.Build()
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// AddQueryParamsFromStruct adds query parameters from a struct using reflection
func (rb *requestBuilder) AddQueryParamsFromStruct(opts interface{}) RequestBuilder {
	val := reflect.ValueOf(opts)
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

		// Split the json tag to get the name and options
		parts := strings.Split(jsonTag, ",")
		if len(parts) == 0 {
			// TODO: figure out migration
			continue
		}
		name := parts[0]

		// Skip empty fields
		if field.IsZero() {
			continue
		}

		// Handle different field types
		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				rb.addQueryParam(name, field.Elem())
			}
		case reflect.Slice:
			if field.Len() > 0 {
				// Special handling for specific array fields
				switch name {
				case "wait_until":
					for j := 0; j < field.Len(); j++ {
						rb.QueryParam("wait_until", fmt.Sprint(field.Index(j).Interface()))
					}
				case "scripts_wait_until":
					for j := 0; j < field.Len(); j++ {
						rb.QueryParam("scripts_wait_until", fmt.Sprint(field.Index(j).Interface()))
					}
				case "hide_selectors":
					for j := 0; j < field.Len(); j++ {
						rb.QueryParam("hide_selector", fmt.Sprint(field.Index(j).Interface()))
					}
				case "block_requests":
					for j := 0; j < field.Len(); j++ {
						rb.QueryParam("block_request", fmt.Sprint(field.Index(j).Interface()))
					}
				case "block_resources":
					for j := 0; j < field.Len(); j++ {
						rb.QueryParam("block_resources", fmt.Sprint(field.Index(j).Interface()))
					}
				default:
					for j := 0; j < field.Len(); j++ {
						rb.QueryParam(name, fmt.Sprint(field.Index(j).Interface()))
					}
				}
			}
		case reflect.Map:
			if field.Len() > 0 {
				iter := field.MapRange()
				for iter.Next() {
					paramName := fmt.Sprintf("%s[%s]", name, iter.Key().String())
					rb.QueryParam(paramName, fmt.Sprint(iter.Value().Interface()))
				}
			}
		case reflect.String:
			if field.String() != "" {
				rb.QueryParam(name, field.String())
			}
		}
	}
	return rb
}

// addQueryParam adds a single query parameter based on the field value
func (rb *requestBuilder) addQueryParam(name string, value reflect.Value) {
	switch value.Kind() {
	case reflect.String:
		rb.QueryParam(name, value.String())
	case reflect.Bool:
		rb.QueryParam(name, strconv.FormatBool(value.Bool()))
	case reflect.Int, reflect.Int64:
		rb.QueryParam(name, strconv.FormatInt(value.Int(), 10))
	case reflect.Float64:
		rb.QueryParam(name, strconv.FormatFloat(value.Float(), 'f', -1, 64))
	}
}
