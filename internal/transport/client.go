package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second
const defaultMaxResponseSize = 1 << 20

// Client centralizes safe HTTP mechanics while leaving provider payloads to adapters.
type Client struct {
	HTTP            *http.Client
	UserAgent       string
	MaxResponseSize int64
}

// New returns a transport client with safe defaults.
func New(httpClient *http.Client, userAgent string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{HTTP: httpClient, UserAgent: userAgent, MaxResponseSize: defaultMaxResponseSize}
}

// Do executes a context-bound request and returns a bounded response body.
func (c *Client) Do(ctx context.Context, req *http.Request) (int, []byte, error) {
	req = req.WithContext(ctx)
	if c.UserAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	limit := c.MaxResponseSize
	if limit <= 0 {
		limit = defaultMaxResponseSize
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read response: %w", err)
	}
	if int64(len(body)) > limit {
		return resp.StatusCode, nil, fmt.Errorf("response exceeds %d bytes", limit)
	}
	return resp.StatusCode, body, nil
}
