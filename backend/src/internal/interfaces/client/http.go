// ./src/internal/interfaces/client/http.go
package interfaces

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
)

// HTTPClient defines the interface for HTTP operations
type HTTPClient interface {
	// Core methods
	Do(req *http.Request) (*http.Response, error)
	NewRequest() RequestBuilder

	// Convenience methods
	Get(path string) RequestBuilder
	Post(path string) RequestBuilder
	Put(path string) RequestBuilder
	Delete(path string) RequestBuilder

	// Direct execution methods
	DoGet(ctx context.Context, path string) (*http.Response, error)
	DoPost(ctx context.Context, path string, body interface{}) (*http.Response, error)
	DoPut(ctx context.Context, path string, body interface{}) (*http.Response, error)
	DoDelete(ctx context.Context, path string) (*http.Response, error)

	// JSON convenience methods
	GetJSON(ctx context.Context, path string, v interface{}) error
	PostJSON(ctx context.Context, path string, body, v interface{}) error
	PutJSON(ctx context.Context, path string, body, v interface{}) error
}

// RequestBuilder interface defines methods for building HTTP requests
type RequestBuilder interface {
	Method(method string) RequestBuilder
	BaseURL(baseURL string) RequestBuilder
	Path(path string) RequestBuilder
	QueryParam(key, value string) RequestBuilder
	Header(key, value string) RequestBuilder
	JSON(data interface{}) RequestBuilder
	Form(data url.Values) RequestBuilder
	MultipartForm(fn func(*multipart.Writer) error) RequestBuilder
	Context(ctx context.Context) RequestBuilder
	AddQueryParamsFromStruct(opts interface{}) RequestBuilder
	Build() (*http.Request, error)
	Execute(c HTTPClient) (*http.Response, error)
}
