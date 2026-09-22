package smskit

import "fmt"

// SMSMessage is the provider-independent representation of one SMS.
// Provider-specific fields belong in the provider package.
type SMSMessage struct {
	To          string
	From        string
	Body        string
	CallbackURL string
	Metadata    map[string]string
}

// Validate checks fields that are required by every SMS provider.
func (m SMSMessage) Validate() error {
	if m.To == "" {
		return fmt.Errorf("sms recipient is required: %w", ErrInvalidRequest)
	}
	if m.Body == "" {
		return fmt.Errorf("sms body is required: %w", ErrInvalidRequest)
	}
	return nil
}
