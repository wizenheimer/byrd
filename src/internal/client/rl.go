package client

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimitedClient wraps HTTPClient with rate limiting
type RateLimitedClient struct {
	client  HTTPClient
	limiter *rate.Limiter
}

// NewRateLimitedClient creates a new rate-limited HTTP client
func NewRateLimitedClient(client HTTPClient, qps float64) HTTPClient {
	return &RateLimitedClient{
		client:  client,
		limiter: rate.NewLimiter(rate.Limit(qps), 1), // burst size of 1
	}
}

// Do implements HTTPClient interface with rate limiting
func (c *RateLimitedClient) Do(req *http.Request) (*http.Response, error) {
	err := c.limiter.Wait(req.Context())
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

// Forward all other HTTPClient methods to the underlying client
func (c *RateLimitedClient) NewRequest() *RequestBuilder {
	return c.client.NewRequest()
}

func (c *RateLimitedClient) Get(path string) *RequestBuilder {
	return c.client.Get(path)
}

func (c *RateLimitedClient) Post(path string) *RequestBuilder {
	return c.client.Post(path)
}

func (c *RateLimitedClient) Put(path string) *RequestBuilder {
	return c.client.Put(path)
}

func (c *RateLimitedClient) Delete(path string) *RequestBuilder {
	return c.client.Delete(path)
}

func (c *RateLimitedClient) DoGet(ctx context.Context, path string) (*http.Response, error) {
	return c.client.DoGet(ctx, path)
}

func (c *RateLimitedClient) DoPost(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.client.DoPost(ctx, path, body)
}

func (c *RateLimitedClient) DoPut(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.client.DoPut(ctx, path, body)
}

func (c *RateLimitedClient) DoDelete(ctx context.Context, path string) (*http.Response, error) {
	return c.client.DoDelete(ctx, path)
}

func (c *RateLimitedClient) GetJSON(ctx context.Context, path string, v interface{}) error {
	return c.client.GetJSON(ctx, path, v)
}

func (c *RateLimitedClient) PostJSON(ctx context.Context, path string, body, v interface{}) error {
	return c.client.PostJSON(ctx, path, body, v)
}

func (c *RateLimitedClient) PutJSON(ctx context.Context, path string, body, v interface{}) error {
	return c.client.PutJSON(ctx, path, body, v)
}
