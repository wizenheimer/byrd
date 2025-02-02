package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Client is a minimal HTTP client wrapper with retry and rate limiting
type HTTPClient struct {
	logger     *logger.Logger
	httpClient *http.Client
	retryCodes []int
	maxRetries int
	limiter    *rate.Limiter
}

type ClientOption func(*HTTPClient)

func WithRetry(maxRetries int, retryCodes []int) ClientOption {
	return func(c *HTTPClient) {
		c.maxRetries = maxRetries
		c.retryCodes = retryCodes
	}
}

func WithRateLimit(rps float64, burst int) ClientOption {
	return func(c *HTTPClient) {
		c.limiter = rate.NewLimiter(rate.Limit(rps), burst)
	}
}

func NewClient(logger *logger.Logger, opts ...ClientOption) (*HTTPClient, error) {
	c := &HTTPClient{
		httpClient: http.DefaultClient,
		maxRetries: 3,                                   // default retries
		retryCodes: []int{408, 429, 500, 502, 503, 504}, // default retry codes
		logger: logger.WithFields(map[string]interface{}{
			"component": "http_client",
		}),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.logger == nil {
		return nil, errors.New("logger is required")
	}

	// Deduplicate the retry codes
	retryCodeMap := make(map[int]struct{})
	for _, code := range c.retryCodes {
		retryCodeMap[code] = struct{}{}
	}
	c.retryCodes = make([]int, 0, len(retryCodeMap))
	for code := range retryCodeMap {
		c.retryCodes = append(c.retryCodes, code)
	}

	return c, nil
}

// shouldRetry checks if the status code is in the retry codes list
func (c *HTTPClient) shouldRetry(statusCode int) bool {
	for _, code := range c.retryCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.limiter != nil {
		if err := c.limiter.Wait(req.Context()); err != nil {
			return nil, fmt.Errorf("rate limit wait: %w", err)
		}
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-req.Context().Done():
				if resp != nil {
					resp.Body.Close()
				}
				return nil, req.Context().Err()
			case <-time.After(backoff):
			}
		}

		// Execute the request
		resp, err = c.httpClient.Do(req)
		if err != nil {
			c.logger.Error("http request failed", zap.Error(err), zap.Any("attempt", attempt), zap.Any("maxRetries", c.maxRetries), zap.Any("url", req.URL))
			// Network-level error
			if attempt == c.maxRetries {
				return nil, fmt.Errorf("max retries reached: %w", err)
			}
			continue
		}

		// Check if we should retry based on status code
		if c.shouldRetry(resp.StatusCode) {
			resp.Body.Close()
			if attempt == c.maxRetries {
				return nil, fmt.Errorf("max retries reached, last status code: %d", resp.StatusCode)
			}
			continue
		}

		// Success! Return the response
		return resp, nil
	}

	// This should never happen, but just in case
	if resp != nil {
		resp.Body.Close()
	}
	return nil, fmt.Errorf("unexpected end of retry loop")
}
