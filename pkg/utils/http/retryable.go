package http

import (
	"net/http"
	"time"
)

// RetryableClient wraps around http.Client and adds retry logic
type RetryableClient struct {
	client     *http.Client
	retryCodes []int
	timeout    time.Duration
	maxRetries int
}

// NewRetryableClient creates a new RetryableClient
func NewRetryableClient(client *http.Client, retryCodes []int, timeout time.Duration, maxRetries int) *RetryableClient {
	return &RetryableClient{
		client:     client,
		retryCodes: retryCodes,
		timeout:    timeout,
		maxRetries: maxRetries,
	}
}

// Do sends an HTTP request and retries on specified error codes
func (rc *RetryableClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rc.maxRetries; attempt++ {
		client := &http.Client{
			Timeout: rc.timeout,
		}

		resp, err = client.Do(req)
		if err == nil && !containsStatusCode(rc.retryCodes, resp.StatusCode) {
			return resp, nil
		}

		if attempt < rc.maxRetries {
			time.Sleep(time.Second * time.Duration(attempt+1)) // Exponential backoff
		}
	}

	return resp, err
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
