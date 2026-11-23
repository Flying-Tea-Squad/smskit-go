package wasiliana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	smskit "github.com/Flying-Tea-Squad/smskit-go"
	"github.com/Flying-Tea-Squad/smskit-go/internal/transport"
)

const (
	providerName   = "wasiliana"
	defaultBaseURL = "https://api.wasiliana.com"
)

// Config configures a Wasiliana client.
type Config struct {
	APIKey     string
	BaseURL    string       // optional; tests point this at httptest
	HTTPClient *http.Client // optional; never modified
	UserAgent  string       // optional
}

// Client talks to the Wasiliana API. It is safe for concurrent use.
type Client struct {
	apiKey    string
	baseURL   string
	transport *transport.Client
}

// New validates cfg and returns a Client. It makes no network calls.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("wasiliana: api key is required: %w", smskit.ErrInvalidRequest)
	}
	base := cfg.BaseURL
	if base == "" {
		base = defaultBaseURL
	}
	return &Client{
		apiKey:    cfg.APIKey,
		baseURL:   strings.TrimRight(base, "/"),
		transport: transport.New(cfg.HTTPClient, cfg.UserAgent),
	}, nil
}

// postJSON sends payload as JSON with apiKey authentication.
func (c *Client) postJSON(ctx context.Context, path string, payload any) (int, []byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, fmt.Errorf("wasiliana: encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return 0, nil, fmt.Errorf("wasiliana: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("apiKey", c.apiKey)
	status, resp, err := c.transport.Do(ctx, req)
	if err != nil {
		// Never retried: the provider may already have the message.
		return status, nil, &smskit.ProviderError{
			Provider: providerName, Kind: smskit.ErrorTransport,
			Message: err.Error(), HTTPStatus: status,
		}
	}
	return status, resp, nil
}
