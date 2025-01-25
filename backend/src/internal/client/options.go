package client

import (
	"time"

	"github.com/wizenheimer/byrd/src/pkg/logger"
)

// WithMaxRetries sets the maximum number of retries for the client
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithRetryCodes sets the status codes to retry on
func WithRetryCodes(codes []int) ClientOption {
	return func(c *Client) {
		c.retryCodes = codes
	}
}

// WithTimeout sets the timeout for the client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithAuth adds authentication to the client
func WithAuth(auth AuthMethod) ClientOption {
	return func(c *Client) {
		c.authMethod = auth
	}
}

// WithLogger sets the logger for the client
func WithLogger(logger *logger.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}
