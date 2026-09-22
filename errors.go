package messaging

import (
	"encoding/json"
	"fmt"
)

// ErrorKind categorizes an error without requiring callers to parse text.
type ErrorKind string

const (
	ErrorInvalidRequest ErrorKind = "invalid_request"
	ErrorAuthentication ErrorKind = "authentication"
	ErrorAuthorization  ErrorKind = "authorization"
	ErrorRateLimited    ErrorKind = "rate_limited"
	ErrorRejected       ErrorKind = "rejected"
	ErrorTemporary      ErrorKind = "temporary"
	ErrorTransport      ErrorKind = "transport"
	ErrorProvider       ErrorKind = "provider"
)

// ErrInvalidRequest marks input rejected before a provider request.
var ErrInvalidRequest = fmt.Errorf("invalid request")

// ProviderError is a safe, structured provider failure.
type ProviderError struct {
	Provider   string
	Kind       ErrorKind
	Code       string
	Message    string
	Retryable  bool
	HTTPStatus int
	Raw        json.RawMessage
}

// Error implements error.
func (e *ProviderError) Error() string {
	if e.Code == "" {
		return fmt.Sprintf("%s: %s", e.Provider, e.Message)
	}
	return fmt.Sprintf("%s: %s: %s", e.Provider, e.Code, e.Message)
}
