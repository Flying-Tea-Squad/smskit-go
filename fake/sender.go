package fake

import (
	"context"
	"sync"

	messaging "github.com/Flying-Tea-Squad/smskit-go"
)

// Sender is a configurable SMS sender for application tests.
type Sender struct {
	mu       sync.Mutex
	Messages []messaging.SMSMessage
	Result   *messaging.SendResult
	Err      error
}

var _ messaging.SMSSender = (*Sender)(nil)

// SendSMS records the message and returns the configured result or error.
func (s *Sender) SendSMS(_ context.Context, message messaging.SMSMessage) (*messaging.SendResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, message)
	return s.Result, s.Err
}
