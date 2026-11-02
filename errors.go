package smskit

import (
	"encoding/json"
	"fmt"
)

// ErrorKind categorizes an error without requiring callers to parse text.
type ErrorKind string

const (
	// ErrorInvalidRequest means the request was invalid or incomplete.
	ErrorInvalidRequest ErrorKind = "invalid_request"
	// ErrorAuthentication means the provider doesn't accept your login or API key.
	ErrorAuthentication ErrorKind = "authentication"
	// ErrorAuthorization means you're logged in but not allowed to do this.
	ErrorAuthorization ErrorKind = "authorization"
	// ErrorRateLimited means you sent too many requests too fast.
	ErrorRateLimited ErrorKind = "rate_limited"
	// ErrorRejected means the provider refused the message.
	ErrorRejected ErrorKind = "rejected"
	// ErrorTemporary means a short problem, so trying again may work.
	ErrorTemporary ErrorKind = "temporary"
	// ErrorTransport means a network problem, like a lost connection.
	ErrorTransport ErrorKind = "transport"
	// ErrorProvider means the provider itself had an error.
	ErrorProvider ErrorKind = "provider"
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
