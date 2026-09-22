package fake

import (
	"context"
	"sync"

	smskit "github.com/Flying-Tea-Squad/smskit-go"
)

// Sender is a configurable SMS sender for application tests.
type Sender struct {
	mu       sync.Mutex
	Messages []smskit.SMSMessage
	Result   *smskit.SendResult
	Err      error
}

var _ smskit.SMSSender = (*Sender)(nil)

// SendSMS records the message and returns the configured result or error.
func (s *Sender) SendSMS(_ context.Context, message smskit.SMSMessage) (*smskit.SendResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, message)
	return s.Result, s.Err
}
