// ./src/internal/client/client.go
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	clf "github.com/wizenheimer/byrd/src/internal/interfaces/client"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// Client implements the HTTPClient interface
type Client struct {
	client     *http.Client
	authMethod clf.AuthMethod
	retryCodes []int
	maxRetries int
	logger     *logger.Logger
}

// ClientOption defines the functional options for Client
type ClientOption func(*Client)

// NewClient creates a new Client with the given options
func NewClient(options ...ClientOption) (clf.HTTPClient, error) {
	c := &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryCodes: []int{
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusRequestTimeout,
		},
		maxRetries: 3,
	}

	for _, opt := range options {
		opt(c)
	}

	if c.logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return c, nil
}

// Core implementation of Do method
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.authMethod != nil {
		c.logger.Debug("applying auth for request",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()))
		c.authMethod.Apply(req)
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		c.logger.Debug("executing request",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Int("attempt", attempt))

		resp, err = c.client.Do(req)
		if err == nil && !containsStatusCode(c.retryCodes, resp.StatusCode) {
			return resp, nil
		}

		c.logger.Debug("request failed",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Int("attempt", attempt),
			zap.Error(err),
			zap.Int("status_code", resp.StatusCode))

		if attempt < c.maxRetries {
			time.Sleep(time.Second * time.Duration(attempt+1))
		}
	}

	return resp, err
}

// NewRequest creates a new RequestBuilder
func (c *Client) NewRequest() clf.RequestBuilder {
	return &RequestBuilder{
		queryParams: url.Values{},
		headers:     make(map[string]string),
		ctx:         context.Background(),
	}
}

// Convenience methods
func (c *Client) Get(path string) clf.RequestBuilder {
	return c.NewRequest().Method(http.MethodGet).Path(path)
}

func (c *Client) Post(path string) clf.RequestBuilder {
	return c.NewRequest().Method(http.MethodPost).Path(path)
}

func (c *Client) Put(path string) clf.RequestBuilder {
	return c.NewRequest().Method(http.MethodPut).Path(path)
}

func (c *Client) Delete(path string) clf.RequestBuilder {
	return c.NewRequest().Method(http.MethodDelete).Path(path)
}

// Direct execution methods

// DoGet performs a GET request
func (c *Client) DoGet(ctx context.Context, path string) (*http.Response, error) {
	return c.Get(path).Context(ctx).Execute(c)
}

// DoPost performs a POST request
func (c *Client) DoPost(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Post(path).Context(ctx).JSON(body).Execute(c)
}

// DoPut performs a PUT request
func (c *Client) DoPut(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Put(path).Context(ctx).JSON(body).Execute(c)
}

// DoDelete performs a DELETE request
func (c *Client) DoDelete(ctx context.Context, path string) (*http.Response, error) {
	return c.Delete(path).Context(ctx).Execute(c)
}

// JSON convenience methods

// GetJSON performs a GET request and decodes the response into v
func (c *Client) GetJSON(ctx context.Context, path string, v interface{}) error {
	resp, err := c.DoGet(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// PostJSON performs a POST request with body and decodes the response into v
func (c *Client) PostJSON(ctx context.Context, path string, body, v interface{}) error {
	resp, err := c.DoPost(ctx, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// PutJSON performs a PUT request with body and decodes the response into v
func (c *Client) PutJSON(ctx context.Context, path string, body, v interface{}) error {
	resp, err := c.DoPut(ctx, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// contains checks if a slice contains a specific integer
func containsStatusCode(retryCodes []int, statusCode int) bool {
	for _, v := range retryCodes {
		if v == statusCode {
			return true
		}
	}
	return false
}
